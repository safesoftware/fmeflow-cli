/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type PublishedParameter struct {
	Name string `json:"name"`
}

type SimpleParameter struct {
	Value string `json:"value"`
	PublishedParameter
}

type ListParameter struct {
	Value []string `json:"value"`
	PublishedParameter
}

type Directive struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type JobId struct {
	Id int `json:"id"`
}

type Job struct {
	PublishedParameters    []interface{}     `json:"-"`
	RawPublishedParameters []json.RawMessage `json:"publishedParameters,omitempty"`
	TMDirectives           struct {
		Rtc         bool   `json:"rtc"`
		Ttc         int    `json:"ttc,omitempty"`
		Description string `json:"description,omitempty"`
		Tag         string `json:"tag,omitempty"`
		TTL         int    `json:"ttl,omitempty"`
	} `json:"TMDirectives,omitempty"`
	NMDirectives struct {
		Directives    []Directive `json:"directives,omitempty"`
		SuccessTopics []string    `json:"successTopics,omitempty"`
		FailureTopics []string    `json:"failureTopics,omitempty"`
	} `json:"NMDirectives,omitempty"`
}

type JobResult struct {
	TimeRequested       time.Time `json:"timeRequested"`
	RequesterResultPort int       `json:"requesterResultPort"`
	NumFeaturesOutput   int       `json:"numFeaturesOutput"`
	RequesterHost       string    `json:"requesterHost"`
	TimeStarted         time.Time `json:"timeStarted"`
	ID                  int       `json:"id"`
	TimeFinished        time.Time `json:"timeFinished"`
	Priority            int       `json:"priority"`
	StatusMessage       string    `json:"statusMessage"`
	Status              string    `json:"status"`
}

var runWorkspace string
var runRepository string
var runWait bool
var runRtc bool
var runTtc int
var runTtl int
var runTag string
var runDescription string
var runSourceData string
var runSuccessTopics []string
var runFailureTopics []string
var runPublishedParameter []string
var runNodeManagerDirective []string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a workspace on FME Server",
	Long:  `Run a workspace on FME Server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			outputType = "json"
		}
		// set up http
		client := &http.Client{
			// set a long timeout for jobs that are long running.
			// maybe this should be a parameter?
			Timeout: 604800 * time.Second,
		}

		var result JobResult
		var responseData []byte

		if runSourceData == "" {
			job := &Job{}

			// get published parameters
			for _, parameter := range runPublishedParameter {
				this_parameter := strings.Split(parameter, "=")
				if strings.Contains(this_parameter[1], ",") {
					var a ListParameter
					a.Name = this_parameter[0]
					this_list := strings.Split(this_parameter[1], ",")
					a.Value = this_list
					job.PublishedParameters = append(job.PublishedParameters, a)
				} else {
					var a SimpleParameter
					a.Name = this_parameter[0]
					a.Value = this_parameter[1]
					job.PublishedParameters = append(job.PublishedParameters, a)
				}
			}

			// get node manager directives
			for _, directive := range runNodeManagerDirective {
				this_directive := strings.Split(directive, "=")
				var a Directive
				a.Name = this_directive[0]
				a.Value = this_directive[1]
				job.NMDirectives.Directives = append(job.NMDirectives.Directives, a)
			}

			if runTtc != -1 {
				job.TMDirectives.Ttc = runTtc
			}
			if runTtl != -1 {
				job.TMDirectives.TTL = runTtl
			}

			job.TMDirectives.Rtc = runRtc

			// append slice to slice
			job.NMDirectives.SuccessTopics = append(job.NMDirectives.SuccessTopics, runSuccessTopics...)
			job.NMDirectives.FailureTopics = append(job.NMDirectives.FailureTopics, runFailureTopics...)

			if runDescription != "" {
				job.TMDirectives.Description = runDescription
			}

			jobJson, err := json.Marshal(job)
			if err != nil {
				return err
			}

			submitEndpoint := "submit"
			if runWait {
				submitEndpoint = "transact"
			}

			endpoint := "/fmerest/v3/transformations/" + submitEndpoint + "/" + runRepository + "/" + runWorkspace

			request, err := buildFmeServerRequest(endpoint, "POST", strings.NewReader(string(jobJson)))
			if err != nil {
				return err
			}

			request.Header.Add("Content-Type", "application/json")

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != 200 && response.StatusCode != 202 {
				return errors.New(response.Status)
			}

			responseData, err = io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			if !runWait {
				var result JobId
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				} else {
					if !jsonOutput {
						fmt.Println("Job submitted with id: " + strconv.Itoa(result.Id))
					} else {
						prettyJSON, err := prettyPrintJSON(responseData)
						if err != nil {
							return err
						}
						fmt.Println(prettyJSON)
					}
				}
			} else {
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				}
			}
		} else {
			// we are uploading a source file, so we want to send the file in the body as octet stream, and parameters as url parameters
			file, err := os.Open(runSourceData)
			if err != nil {
				return err
			}
			defer file.Close()

			endpoint := "/fmerest/v3/transformations/transactdata/" + runRepository + "/" + runWorkspace
			request, err := buildFmeServerRequest(endpoint, "POST", file)
			if err != nil {
				return err
			}

			q := request.URL.Query()

			if runDescription != "" {
				q.Add("opt_description", runDescription)
			}

			for _, topic := range runSuccessTopics {
				q.Add("opt_successtopics", topic)
			}

			for _, topic := range runFailureTopics {
				q.Add("opt_failuretopics", topic)
			}

			if runDescription != "" {
				endpoint += "opt_description=" + runDescription
			}

			if runTag != "" {
				q.Add("opt_tag", runTag)
			}

			if runTtl != -1 {
				q.Add("opt_ttl", strconv.Itoa(runTtl))
			}

			if runTtc != -1 {
				q.Add("opt_ttc", strconv.Itoa(runTtc))
			}

			// TODO: I'm not sure this is the correct way to pass published parameters in the query string
			for _, parameter := range runPublishedParameter {
				this_parameter := strings.Split(parameter, "=")
				q.Add(this_parameter[0], this_parameter[1])
			}

			request.URL.RawQuery = q.Encode()

			request.Header.Set("Content-Type", "application/octet-stream")

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != 200 {
				return errors.New(response.Status)
			}

			responseData, err = io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			}
		}

		if runWait {
			if outputType == "table" {
				t := table.NewWriter()
				t.SetStyle(defaultStyle)

				t.AppendHeader(table.Row{"ID", "Status", "Status Message", "Features Output"})

				t.AppendRow(table.Row{result.ID, result.Status, result.StatusMessage, result.NumFeaturesOutput})

				if noHeaders {
					t.ResetHeaders()
				}
				fmt.Println(t.Render())

			} else if outputType == "json" {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Println(prettyJSON)
			} else if strings.HasPrefix(outputType, "custom-columns=") {
				// parse the columns and json queries
				columnsString := outputType[len("custom-columns="):]
				if len(columnsString) == 0 {
					return errors.New("custom-columns format specified but no custom columns given")
				}

				// we have to marshal the Items array, then create an array of marshalled items
				// to pass to the creation of the table.
				marshalledItems := [][]byte{}

				mJson, err := json.Marshal(result)
				if err != nil {
					return err
				}

				marshalledItems = append(marshalledItems, mJson)

				columnsInput := strings.Split(columnsString, ",")
				t, err := createTableFromCustomColumns(marshalledItems, columnsInput)
				if err != nil {
					return err
				}
				if noHeaders {
					t.ResetHeaders()
				}
				fmt.Println(t.Render())

			} else {
				return errors.New("invalid output format specified")
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&runRepository, "repository", "", "The name of the repository containing the workspace to run.")
	runCmd.Flags().StringVar(&runWorkspace, "workspace", "", "The name of the workspace to run.")
	runCmd.Flags().BoolVar(&runWait, "wait", false, "Submit job and wait for it to finish.")
	runCmd.Flags().BoolVar(&runRtc, "run-until-canceled", false, "Runs a job until it is explicitly canceled. The job will run again regardless of whether the job completed successfully, failed, or the server crashed or was shut down.")
	runCmd.Flags().IntVar(&runTtc, "time-until-canceled", -1, "Time (in seconds) elapsed for a running job before it's cancelled. The minimum value is 1 second, values less than 1 second are ignored.")
	runCmd.Flags().IntVar(&runTtl, "time-to-live", -1, "Time to live in the job queue (in seconds)")
	runCmd.Flags().StringVar(&runTag, "tag", "", "The job routing tag for the request")
	runCmd.Flags().StringVar(&runDescription, "description", "", "Description of the request.")
	runCmd.Flags().StringVar(&runSourceData, "file", "", "Upload a local file Source dataset to use to run the workspace.")
	runCmd.Flags().StringVarP(&outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	runCmd.Flags().BoolVar(&noHeaders, "no-headers", false, "Don't print column headers")

	runCmd.Flags().StringArrayVar(&runSuccessTopics, "success-topic", []string{}, "Topics to notify when the job succeeds. Can be specified more than once.")
	runCmd.Flags().StringArrayVar(&runFailureTopics, "failure-topic", []string{}, "Topics to notify when the job fails. Can be specified more than once.")
	runCmd.Flags().StringArrayVar(&runPublishedParameter, "published-parameter", []string{}, "Workspace published parameters defined for this job. Specify as Key=Value. Can be passed in multiple times. For list parameters, specify as Key=Value1,Value2. This means parameter values can't contain = or , at the moment. That should probably be fixed.")
	runCmd.Flags().StringArrayVar(&runNodeManagerDirective, "node-manager-directive", []string{}, "Additional NM Directives, which can include client-configured keys, to pass to the notification service for custom use by subscriptions. Specify as Key=Value Can be passed in multiple times.")
}
