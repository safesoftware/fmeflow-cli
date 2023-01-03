package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type JobStatus struct {
	Request       JobRequest `json:"request"`
	TimeDelivered time.Time  `json:"timeDelivered"`
	Workspace     string     `json:"workspace"`
	NumErrors     int        `json:"numErrors"`
	NumLines      int        `json:"numLines"`
	EngineHost    string     `json:"engineHost"`
	TimeQueued    time.Time  `json:"timeQueued"`
	CPUPct        float64    `json:"cpuPct"`
	Description   string     `json:"description"`
	TimeStarted   time.Time  `json:"timeStarted"`
	Repository    string     `json:"repository"`
	UserName      string     `json:"userName"`
	Result        JobResult  `json:"result"`
	CPUTime       int        `json:"cpuTime"`
	ID            int        `json:"id"`
	TimeFinished  time.Time  `json:"timeFinished"`
	EngineName    string     `json:"engineName"`
	NumWarnings   int        `json:"numWarnings"`
	TimeSubmitted time.Time  `json:"timeSubmitted"`
	ElapsedTime   int        `json:"elapsedTime"`
	PeakMemUsage  int        `json:"peakMemUsage"`
	Status        string     `json:"status"`
}

type Jobs struct {
	Offset     int         `json:"offset"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"totalCount"`
	Items      []JobStatus `json:"items"`
}

type jobsFlags struct {
	outputType     string
	noHeaders      bool
	jobsRunning    bool
	jobsCompleted  bool
	jobsActive     bool
	jobsAll        bool
	jobsQueued     bool
	jobsRepository string
	jobsUserName   string
	jobsWorkspace  string
	jobsSourceID   string
	jobsSourceType string
	jobId          int
}

func newJobsCmd() *cobra.Command {
	f := jobsFlags{}
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Lists jobs on FME Server",
		Long:  "Lists jobs on FME Server",

		Example: `
  # List all jobs (currently limited to the most recent 1000)
  fmeserver jobs --all
	
  # List all running jobs
  fmeserver jobs --running
	
  # List all jobs from a given repository
  fmeserver jobs --repository Samples
	
  # List all jobs that ran a given workspace
  fmeserver jobs --repository Samples --workspace austinApartments.fmw
	
  # List all jobs in JSON format
  fmeserver jobs --json
	
  # List the workspace, CPU time and peak memory usage for a given repository
  fmeserver jobs --repository Samples --output="custom-columns=WORKSPACE:.workspace,CPU Time:.cpuTime,Peak Memory:.peakMemUsage"
	`,
		Args: NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			if f.jobsWorkspace != "" {
				cmd.MarkFlagRequired("repository")
			}
		},
		RunE: jobsRun(&f),
	}
	cmd.Flags().BoolVar(&f.jobsRunning, "running", false, "Retrieve running jobs")
	cmd.Flags().BoolVar(&f.jobsCompleted, "completed", false, "Retrieve completed jobs")
	cmd.Flags().BoolVar(&f.jobsQueued, "queued", false, "Retrieve queued jobs")
	cmd.Flags().BoolVar(&f.jobsAll, "all", false, "Retrieve all jobs")
	cmd.Flags().BoolVar(&f.jobsActive, "active", false, "Retrieve active jobs")
	cmd.Flags().StringVar(&f.jobsRepository, "repository", "", "If specified, only jobs from the specified repository will be returned.")
	cmd.Flags().StringVar(&f.jobsWorkspace, "workspace", "", "If specified along with repository, only jobs from the specified repository and workspace will be returned.")
	cmd.Flags().StringVar(&f.jobsUserName, "user-name", "", "If specified, only jobs run by the specified user will be returned.")
	cmd.Flags().StringVar(&f.jobsSourceID, "source-id", "", "If specified along with source type, only jobs from the specified type with the specified id will be returned. For Automations, the source id is the automation id. For WorkspaceSubscriber, the source id is the id of the subscription. For Scheduler, the source id is the category and name of the schedule separated by '/'. For example, 'Category/Name'.")
	cmd.Flags().StringVar(&f.jobsSourceType, "source-type", "", "If specified, only jobs run by this source type will be returned.")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().IntVar(&f.jobId, "id", -1, "Specify the job id to display")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.MarkFlagsMutuallyExclusive("queued", "active")
	cmd.MarkFlagsMutuallyExclusive("running", "active")
	cmd.MarkFlagsMutuallyExclusive("id", "running")
	cmd.MarkFlagsMutuallyExclusive("id", "completed")
	cmd.MarkFlagsMutuallyExclusive("id", "queued")
	cmd.MarkFlagsMutuallyExclusive("id", "all")
	cmd.MarkFlagsMutuallyExclusive("id", "active")
	cmd.MarkFlagsMutuallyExclusive("id", "repository")
	cmd.MarkFlagsMutuallyExclusive("id", "workspace")
	cmd.MarkFlagsMutuallyExclusive("id", "user-name")
	cmd.MarkFlagsMutuallyExclusive("id", "source-id")
	cmd.MarkFlagsMutuallyExclusive("id", "source-type")
	return cmd

}

func jobsRun(f *jobsFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}
		if !f.jobsActive && !f.jobsCompleted && !f.jobsQueued && !f.jobsRunning && !f.jobsAll && f.jobId == -1 {
			// if no filter is passed in, show all jobs
			f.jobsAll = true
		}
		var allJobs Jobs
		if f.jobsActive || f.jobsAll {
			err := getJobs("/fmerest/v3/transformations/jobs/active", &allJobs, f)
			if err != nil {
				return err
			}
		}

		if f.jobsCompleted || f.jobsAll {
			err := getJobs("/fmerest/v3/transformations/jobs/completed", &allJobs, f)
			if err != nil {
				return err
			}
		}

		if f.jobsRunning {
			err := getJobs("/fmerest/v3/transformations/jobs/running", &allJobs, f)
			if err != nil {
				return err
			}
		}

		if f.jobsQueued {
			err := getJobs("/fmerest/v3/transformations/jobs/queued", &allJobs, f)
			if err != nil {
				return err
			}
		}

		if f.jobId != -1 {
			err := getJobs("/fmerest/v3/transformations/jobs/id/"+strconv.Itoa(f.jobId), &allJobs, f)
			if err != nil {
				return err
			}
		}

		if f.outputType == "table" {
			// output all values returned by the JSON in a table
			t := table.NewWriter()
			t.SetStyle(defaultStyle)

			t.AppendHeader(table.Row{"Job ID", "Engine Name", "Workspace", "Status"})

			for _, job := range allJobs.Items {
				t.AppendRow(table.Row{job.ID, job.EngineName, job.Workspace, job.Status})
			}
			if f.noHeaders {
				t.ResetHeaders()
			}
			fmt.Fprintln(cmd.OutOrStdout(), t.Render())

		} else if f.outputType == "json" {
			outputjson, err := json.Marshal(allJobs)
			if err != nil {
				return err
			}
			prettyJSON, err := prettyPrintJSON(outputjson)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)

		} else if strings.HasPrefix(f.outputType, "custom-columns") {
			// parse the columns and json queries
			columnsString := ""
			if strings.HasPrefix(f.outputType, "custom-columns=") {
				columnsString = f.outputType[len("custom-columns="):]
			}
			if len(columnsString) == 0 {
				return errors.New("custom-columns format specified but no custom columns given")
			}

			// we have to marshal the Items array, then create an array of marshalled items
			// to pass to the creation of the table.
			marshalledItems := [][]byte{}
			for _, element := range allJobs.Items {
				mJson, err := json.Marshal(element)
				if err != nil {
					return err
				}

				marshalledItems = append(marshalledItems, mJson)
			}

			columnsInput := strings.Split(columnsString, ",")
			t, err := createTableFromCustomColumns(marshalledItems, columnsInput)
			if err != nil {
				return err
			}
			if f.noHeaders {
				t.ResetHeaders()
			}
			fmt.Fprintln(cmd.OutOrStdout(), t.Render())

		} else {
			return errors.New("invalid output format specified")
		}
		return nil
	}
}

func getJobs(endpoint string, allJobs *Jobs, f *jobsFlags) error {
	client := &http.Client{}
	request, err := buildFmeServerRequest(endpoint, "GET", nil)
	if err != nil {
		return err
	}

	q := request.URL.Query()

	if f.jobsRepository != "" {
		q.Add("repository", f.jobsRepository)
	}

	if f.jobsWorkspace != "" {
		q.Add("workspace", f.jobsWorkspace)
	}

	if f.jobsUserName != "" {
		q.Add("userName", f.jobsUserName)
	}

	if f.jobsSourceID != "" {
		q.Add("sourceID", f.jobsSourceID)
	}

	if f.jobsSourceType != "" {
		q.Add("sourceType", f.jobsSourceType)
	}

	request.URL.RawQuery = q.Encode()

	response, err := client.Do(&request)
	if err != nil {
		return err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return errors.New(response.Status)
	}

	if f.jobId == -1 {
		var result Jobs
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			// merge with existing jobs
			allJobs.TotalCount += result.TotalCount
			allJobs.Items = append(allJobs.Items, result.Items...)
		}
	} else {
		var result JobStatus
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			allJobs.TotalCount += 1
			allJobs.Items = append(allJobs.Items, result)
		}
	}

	return nil
}
