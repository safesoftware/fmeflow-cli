/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

type Jobs struct {
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	TotalCount int `json:"totalCount"`
	Items      []struct {
		Request       Job       `json:"request"`
		TimeDelivered time.Time `json:"timeDelivered"`
		Workspace     string    `json:"workspace"`
		NumErrors     int       `json:"numErrors"`
		NumLines      int       `json:"numLines"`
		EngineHost    string    `json:"engineHost"`
		TimeQueued    time.Time `json:"timeQueued"`
		CPUPct        float64   `json:"cpuPct"`
		Description   string    `json:"description"`
		TimeStarted   time.Time `json:"timeStarted"`
		Repository    string    `json:"repository"`
		UserName      string    `json:"userName"`
		Result        JobResult `json:"result"`
		CPUTime       int       `json:"cpuTime"`
		ID            int       `json:"id"`
		TimeFinished  time.Time `json:"timeFinished"`
		EngineName    string    `json:"engineName"`
		NumWarnings   int       `json:"numWarnings"`
		TimeSubmitted time.Time `json:"timeSubmitted"`
		ElapsedTime   int       `json:"elapsedTime"`
		PeakMemUsage  int       `json:"peakMemUsage"`
		Status        string    `json:"status"`
	} `json:"items"`
}

var jobsRunning bool
var jobsCompleted bool
var jobsActive bool
var jobsAll bool
var jobsQueued bool
var jobsRepository string
var jobsUserName string
var jobsWorkspace string
var jobsSourceID string
var jobsSourceType string

// jobsCmd represents the jobs command
var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Lists jobs on FME Server",
	Long:  `Lists jobs on FME Server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var allJobs Jobs
		if jobsActive {
			if jobsQueued {
				return errors.New("cannot specify both --queued and --active")
			}
			if jobsRunning {
				return errors.New("cannot specify both --running and --active")
			}
			err := getJobs("/fmerest/v3/transformations/jobs/active", &allJobs)
			if err != nil {
				return err
			}
		}

		if jobsCompleted {
			err := getJobs("/fmerest/v3/transformations/jobs/completed", &allJobs)
			if err != nil {
				return err
			}
		}

		if jobsRunning {
			err := getJobs("/fmerest/v3/transformations/jobs/running", &allJobs)
			if err != nil {
				return err
			}
		}

		if jobsQueued {
			err := getJobs("/fmerest/v3/transformations/jobs/queued", &allJobs)
			if err != nil {
				return err
			}
		}

		if !jsonOutput {
			// output all values returned by the JSON in a table
			fmt.Printf("Total Jobs: %v\n", allJobs.TotalCount)
			for _, job := range allJobs.Items {
				fmt.Println("-------------------")
				fmt.Printf("Job ID: %v\n", job.ID)
				fmt.Printf("Engine Name: %v\n", job.EngineName)
				fmt.Printf("Workspace: %v\n", job.Workspace)
				fmt.Printf("Status: %v\n", job.Status)

			}
		} else {
			outputjson, err := json.Marshal(allJobs)
			fmt.Printf("%s\n", outputjson)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(jobsCmd)
	jobsCmd.Flags().BoolVar(&jobsRunning, "running", false, "Retrieve running jobs")
	jobsCmd.Flags().BoolVar(&jobsCompleted, "completed", false, "Retrieve completed jobs")
	jobsCmd.Flags().BoolVar(&jobsQueued, "queued", false, "Retrieve queued jobs")
	jobsCmd.Flags().BoolVar(&jobsAll, "all", false, "Retrieve all jobs")
	jobsCmd.Flags().BoolVar(&jobsActive, "active", false, "Retrieve active jobs")
	jobsCmd.Flags().StringVar(&jobsRepository, "repository", "", "If specified, only jobs from the specified repository will be returned.")
	jobsCmd.Flags().StringVar(&jobsWorkspace, "workspace", "", "If specified along with repository, only jobs from the specified repository and workspace will be returned.")
	jobsCmd.Flags().StringVar(&jobsUserName, "user-name", "", "If specified, only jobs run by the specified user will be returned.")
	jobsCmd.Flags().StringVar(&jobsSourceID, "source-id", "", "If specified along with source type, only jobs from the specified type with the specified id will be returned. For Automations, the source id is the automation id. For WorkspaceSubscriber, the source id is the id of the subscription. For Scheduler, the source id is the category and name of the schedule separated by '/'. For example, 'Category/Name'.")
	jobsCmd.Flags().StringVar(&jobsSourceType, "source-type", "", "If specified, only jobs run by this source type will be returned.")

}

func getJobs(endpoint string, allJobs *Jobs) error {
	client := &http.Client{}
	request, err := buildFmeServerRequest(endpoint, "GET", nil)
	if err != nil {
		return err
	}

	q := request.URL.Query()

	if jobsRepository != "" {
		q.Add("repository", jobsRepository)
	}

	if jobsWorkspace != "" {
		q.Add("workspace", jobsWorkspace)
	}

	if jobsUserName != "" {
		q.Add("userName", jobsUserName)
	}

	if jobsSourceID != "" {
		q.Add("sourceID", jobsSourceID)
	}

	if jobsSourceType != "" {
		q.Add("sourceType", jobsSourceType)
	}

	request.URL.RawQuery = q.Encode()

	response, err := client.Do(&request)
	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var result Jobs
	if err := json.Unmarshal(responseData, &result); err != nil {
		return err
	} else {
		// merge with existing jobs
		allJobs.TotalCount += result.TotalCount
		allJobs.Items = append(allJobs.Items, result.Items...)
	}
	return nil
}
