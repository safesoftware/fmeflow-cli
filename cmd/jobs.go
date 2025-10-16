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
	"github.com/spf13/viper"
)

type JobStatusV4 struct {
	ID                       int       `json:"id"`
	Description              string    `json:"description"`
	EngineHost               string    `json:"engineHost"`
	EngineName               string    `json:"engineName"`
	Repository               string    `json:"repository"`
	Queue                    string    `json:"queue"`
	QueueType                string    `json:"queueType"`
	ResultDatasetDownloadURL string    `json:"resultDatasetDownloadUrl"`
	Status                   string    `json:"status"`
	TimeFinished             time.Time `json:"timeFinished"`
	TimeQueued               time.Time `json:"timeQueued"`
	TimeStarted              time.Time `json:"timeStarted"`
	RuntimeUsername          string    `json:"runtimeUsername"`
	RuntimeUserID            string    `json:"runtimeUserID"`
	Workspace                string    `json:"workspace"`
	ElapsedTime              int       `json:"elapsedTime"`
	CPUTime                  int       `json:"cpuTime"`
	CPUPercent               float64   `json:"cpuPercent"`
	PeakMemoryUsage          int       `json:"peakMemoryUsage"`
	LineCount                int       `json:"lineCount"`
	WarningCount             int       `json:"warningCount"`
	ErrorCount               int       `json:"errorCount"`
}

type JobStatusV3 struct {
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

type JobsV4 struct {
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
	TotalCount int           `json:"totalCount"`
	Items      []JobStatusV4 `json:"items"`
}

type JobsV3 struct {
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
	TotalCount int           `json:"totalCount"`
	Items      []JobStatusV3 `json:"items"`
}

type jobsFlags struct {
	outputType     string
	noHeaders      bool
	jobsRunning    bool
	jobsCompleted  bool
	jobsActive     bool
	jobsAll        bool
	jobsQueued     bool
	jobsFailed     bool
	jobsSucceeded  bool
	jobsCancelled  bool
	jobsRepository string
	jobsUserName   string
	jobsWorkspace  string
	jobsSourceID   string
	jobsSourceType string
	jobId          int
	jobStatus      []string
	engineName     string
	queue          string
	sort           string
	apiVersion     apiVersionFlag
}

type account struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	FullName       string  `json:"fullName"`
	Email          string  `json:"email"`
	IsSuperUser    bool    `json:"isSuperUser"`
	Enabled        bool    `json:"enabled"`
	SharingEnabled bool    `json:"sharingEnabled"`
	Type           string  `json:"type"`
	Password       *string `json:"password"`
}

type accountsResponse struct {
	Items      []account `json:"items"`
	TotalCount int       `json:"totalCount"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
}

var jobsV4BuildThreshold = 25208
var activeStatuses = []string{"queued", "running"}
var completedStatuses = []string{"success", "failure", "cancelled"}

func newJobsCmd() *cobra.Command {
	f := jobsFlags{}
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Lists jobs on FME Server",
		Long:  "Lists running, queued, and/or queued jobs on FME Server. Pass in a job id to get information on a specific job.",

		Example: `
  # V3

  # List all jobs (currently limited to the most recent 1000)
  fmeflow jobs --all
	
  # List all running jobs
  fmeflow jobs --running
	
  # List all jobs from a given repository
  fmeflow jobs --repository Samples
	
  # List all jobs that ran a given workspace
  fmeflow jobs --repository Samples --workspace austinApartments.fmw
	
  # List all jobs in JSON format
  fmeflow jobs --json
	
  # List the workspace, CPU time and peak memory usage for a given repository
  fmeflow jobs --repository Samples --output="custom-columns=WORKSPACE:.workspace,CPU Time:.cpuTime,Peak Memory:.peakMemUsage"

  # V4
  # List all jobs
  fmeflow jobs

  # List all jobs with status of success or failure
  fmeflow jobs --status success --status failure

  # List all jobs from a given repository
  fmeflow jobs --repository Samples
	
  # List all jobs that ran a given workspace
  fmeflow jobs --repository Samples --workspace austinApartments.fmw

  # List all jobs in JSON format
  fmeflow jobs --json

  # List the workspace, CPU time and peak memory usage for a given repository
  fmeflow jobs --repository Samples --output="custom-columns=WORKSPACE:.workspace,CPU Time:.cpuTime,Peak Memory:.peakMemoryUsage"
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
	cmd.Flags().BoolVar(&f.jobsCompleted, "completed", false, "Retrieve completed jobs. For v3 API only")
	cmd.Flags().BoolVar(&f.jobsQueued, "queued", false, "Retrieve queued jobs")
	cmd.Flags().BoolVar(&f.jobsAll, "all", false, "Retrieve all jobs")
	cmd.Flags().BoolVar(&f.jobsActive, "active", false, "Retrieve active jobs. For v3 API only")
	cmd.Flags().BoolVar(&f.jobsFailed, "failed", false, "Retrieve failed jobs")
	cmd.Flags().BoolVar(&f.jobsSucceeded, "succeeded", false, "Retrieve succeeded jobs")
	cmd.Flags().BoolVar(&f.jobsCancelled, "cancelled", false, "Retrieve cancelled jobs")
	cmd.Flags().StringVar(&f.jobsRepository, "repository", "", "If specified, only jobs from the specified repository will be returned.")
	cmd.Flags().StringVar(&f.jobsWorkspace, "workspace", "", "If specified along with repository, only jobs from the specified repository and workspace will be returned.")
	cmd.Flags().StringVar(&f.jobsUserName, "user-name", "", "If specified, only jobs run by the specified user will be returned.")
	cmd.Flags().StringVar(&f.jobsSourceID, "source-id", "", "If specified along with source type, only jobs from the specified type with the specified id will be returned. For Automations, the source id is the automation id. For WorkspaceSubscriber, the source id is the id of the subscription. For Scheduler, the source id is the category and name of the schedule separated by '/'. For example, 'Category/Name'.")
	cmd.Flags().StringVar(&f.jobsSourceType, "source-type", "", "If specified, only jobs run by this source type will be returned.")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().IntVar(&f.jobId, "id", -1, "Specify the job id to display")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.Flags().StringVar(&f.engineName, "engine-name", "", "If specified, only jobs run by the specified engine will be returned. Queued jobs cannot be filtered by engine. For v4 API only")
	cmd.Flags().StringVar(&f.queue, "queue", "", "If specified, only jobs routed through the specified queue will be returned. For v4 API only")
	cmd.Flags().StringVar(&f.sort, "sort", "", "Sort jobs by one of: workspace, timeFinished, timeStarted, status. Append _asc or _desc to specify ascending or descending order. For example: workspace_asc. For v4 API only")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.MarkFlagsMutuallyExclusive("queued", "active")
	cmd.MarkFlagsMutuallyExclusive("running", "active")
	cmd.MarkFlagsMutuallyExclusive("id", "running")
	cmd.MarkFlagsMutuallyExclusive("id", "completed")
	cmd.MarkFlagsMutuallyExclusive("id", "queued")
	cmd.MarkFlagsMutuallyExclusive("id", "all")
	cmd.MarkFlagsMutuallyExclusive("id", "active")
	cmd.MarkFlagsMutuallyExclusive("id", "failed")
	cmd.MarkFlagsMutuallyExclusive("id", "succeeded")
	cmd.MarkFlagsMutuallyExclusive("id", "cancelled")
	cmd.MarkFlagsMutuallyExclusive("id", "repository")
	cmd.MarkFlagsMutuallyExclusive("id", "workspace")
	cmd.MarkFlagsMutuallyExclusive("id", "user-name")
	cmd.MarkFlagsMutuallyExclusive("id", "source-id")
	cmd.MarkFlagsMutuallyExclusive("id", "source-type")
	cmd.MarkFlagsMutuallyExclusive("id", "engine-name")
	cmd.MarkFlagsMutuallyExclusive("id", "queue")
	cmd.MarkFlagsMutuallyExclusive("id", "sort")
	cmd.MarkFlagsMutuallyExclusive("id", "user-name")
	cmd.MarkFlagsMutuallyExclusive("active", "engine-name")
	cmd.MarkFlagsMutuallyExclusive("queued", "engine-name")
	return cmd

}

func jobsRun(f *jobsFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		if f.apiVersion == "" {
			if viper.GetInt("build") < jobsV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		if f.apiVersion == apiVersionFlagV4 {
			var allJobs JobsV4
			if f.jobId == -1 {
				if !f.jobsQueued && !f.jobsRunning && !f.jobsFailed && !f.jobsSucceeded && !f.jobsCancelled && !f.jobsAll {
					f.jobsAll = true
				}

				if f.jobsQueued {
					f.jobStatus = append(f.jobStatus, "queued")
				}

				if f.jobsRunning {
					f.jobStatus = append(f.jobStatus, "running")
				}

				if f.jobsFailed {
					f.jobStatus = append(f.jobStatus, "failure")
				}

				if f.jobsSucceeded {
					f.jobStatus = append(f.jobStatus, "success")
				}

				if f.jobsCancelled {
					f.jobStatus = append(f.jobStatus, "cancelled")
				}

				err := getJobsV4("/fmeapiv4/jobs", &allJobs, f)
				if err != nil {
					return err
				}
			} else {
				// get specific job
				err := getJobsV4("/fmeapiv4/jobs/"+strconv.Itoa(f.jobId), &allJobs, f)
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
		} else if f.apiVersion == apiVersionFlagV3 {
			if !f.jobsActive && !f.jobsCompleted && !f.jobsQueued && !f.jobsRunning && !f.jobsAll && f.jobId == -1 {
				// if no filter is passed in, show all jobs
				f.jobsAll = true
			}
			var allJobs JobsV3
			if f.jobsActive || f.jobsAll {
				err := getJobsV3("/fmerest/v3/transformations/jobs/active", &allJobs, f)
				if err != nil {
					return err
				}
			}

			if f.jobsCompleted || f.jobsAll {
				err := getJobsV3("/fmerest/v3/transformations/jobs/completed", &allJobs, f)
				if err != nil {
					return err
				}
			}

			if f.jobsRunning {
				err := getJobsV3("/fmerest/v3/transformations/jobs/running", &allJobs, f)
				if err != nil {
					return err
				}
			}

			if f.jobsQueued {
				err := getJobsV3("/fmerest/v3/transformations/jobs/queued", &allJobs, f)
				if err != nil {
					return err
				}
			}

			if f.jobId != -1 {
				err := getJobsV3("/fmerest/v3/transformations/jobs/id/"+strconv.Itoa(f.jobId), &allJobs, f)
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
		return nil
	}
}

func getJobsV3(endpoint string, allJobs *JobsV3, f *jobsFlags) error {
	client := &http.Client{}
	request, err := buildFmeFlowRequest(endpoint, "GET", nil)
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
		var result JobsV3
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			// merge with existing jobs
			allJobs.TotalCount += result.TotalCount
			allJobs.Items = append(allJobs.Items, result.Items...)
		}
	} else {
		var result JobStatusV3
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			allJobs.TotalCount += 1
			allJobs.Items = append(allJobs.Items, result)
		}
	}
	return nil

}

func getJobsV4(endpoint string, allJobs *JobsV4, f *jobsFlags) error {
	client := &http.Client{}
	request, err := buildFmeFlowRequest(endpoint, "GET", nil)
	if err != nil {
		return err
	}

	q := request.URL.Query()

	if len(f.jobStatus) > 0 {
		for _, status := range f.jobStatus {
			q.Add("status", status)
		}
	} else if f.jobId == -1 {
		if f.engineName != "" {
			activeStatuses = []string{"running"}
		}

		f.jobStatus = activeStatuses
		err := getJobsV4(endpoint, allJobs, f)
		if err != nil {
			return err
		}

		f.jobStatus = completedStatuses
		err = getJobsV4(endpoint, allJobs, f)
		if err != nil {
			return err
		}
		return nil
	}

	if f.jobsRepository != "" {
		q.Add("repository", f.jobsRepository)
	}

	if f.jobsWorkspace != "" {
		q.Add("workspace", f.jobsWorkspace)
	}

	if f.jobsUserName != "" {
		userID, err := GetAccountIDByName(f.jobsUserName)
		if err != nil {
			return fmt.Errorf("failed to find user '%s': %v", f.jobsUserName, err)
		}
		q.Add("runtimeUserID", userID)
	}

	if f.engineName != "" {
		q.Add("engineName", f.engineName)
	}

	if f.jobsSourceType != "" {
		q.Add("sourceType", f.jobsSourceType)
	}

	if f.jobsSourceID != "" {
		q.Add("sourceID", f.jobsSourceID)
	}

	if f.queue != "" {
		q.Add("queue", f.queue)
	}

	if f.sort != "" {
		errorMsg := "invalid sort format, append _asc or _desc to one of: workspace, timeFinished, timeStarted, status"
		elements := strings.SplitN(f.sort, "_", 2)
		property := elements[0]
		order := ""

		if len(elements) > 1 {
			order = elements[1]
		} else {
			return errors.New(errorMsg)
		}

		validProperties := []string{"workspace", "timeFinished", "timeStarted", "status"}
		isValidProperty := false
		for _, validProperty := range validProperties {
			if property == validProperty {
				isValidProperty = true
				break
			}
		}

		if !isValidProperty || (order != "asc" && order != "desc") {
			return errors.New(errorMsg)
		}

		q.Add("sort", f.sort)
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
		return fmt.Errorf("%s: %s", response.Status, string(responseData))
	}

	if f.jobId == -1 {
		var result JobsV4
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			// merge with existing jobs
			allJobs.TotalCount += result.TotalCount
			allJobs.Items = append(allJobs.Items, result.Items...)
		}
	} else {
		var result JobStatusV4
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			allJobs.TotalCount += 1
			allJobs.Items = append(allJobs.Items, result)
		}
	}
	return nil
}

func GetAccountIDByName(accountName string) (string, error) {
	client := &http.Client{}
	limit := 100
	offset := 0

	for {
		request, err := buildFmeFlowRequest("/fmeapiv4/accounts", "GET", nil)
		if err != nil {
			return "", err
		}

		q := request.URL.Query()
		q.Add("limit", strconv.Itoa(limit))
		q.Add("offset", strconv.Itoa(offset))
		q.Add("summary", "true")
		request.URL.RawQuery = q.Encode()

		response, err := client.Do(&request)
		if err != nil {
			return "", err
		}

		if response.StatusCode != 200 {
			response.Body.Close()
			return "", errors.New(response.Status)
		}

		responseData, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return "", err
		}

		var accResponse accountsResponse
		if err := json.Unmarshal(responseData, &accResponse); err != nil {
			return "", err
		}

		for _, acc := range accResponse.Items {
			if acc.Name == accountName {
				return acc.ID, nil
			}
		}

		if len(accResponse.Items) < limit || offset+limit >= accResponse.TotalCount {
			break
		}

		offset += limit
	}

	return "", fmt.Errorf("account name '%s' not found", accountName)
}
