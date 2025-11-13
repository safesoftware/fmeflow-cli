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
	"github.com/spf13/viper"
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

type JobRequestV4 struct {
	Directives          map[string]string      `json:"directives,omitempty"`
	FailureTopics       []string               `json:"failureTopics,omitempty"`
	SuccessTopics       []string               `json:"successTopics,omitempty"`
	MaxJobRuntime       int                    `json:"maxJobRuntime,omitempty"`
	MaxTimeInQueue      int                    `json:"maxTimeInQueue,omitempty"`
	Queue               string                 `json:"queue,omitempty"`
	Repository          string                 `json:"repository,omitempty"`
	Workspace           string                 `json:"workspace,omitempty"`
	PublishedParameters map[string]interface{} `json:"publishedParameters,omitempty"`
	MaxTotalLifeTime    int                    `json:"maxTotalLifeTime,omitempty"`
}

type JobRequestV3 struct {
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

type JobResultV4 struct {
	ID                  int       `json:"id"`
	FeatureOutputCount  int       `json:"featureOutputCount"`
	RequesterHost       string    `json:"requesterHost"`
	RequesterResultPort int       `json:"requesterResultPort"`
	Status              string    `json:"status"`
	StatusMessage       string    `json:"statusMessage"`
	TimeFinished        time.Time `json:"timeFinished"`
	TimeQueued          time.Time `json:"timeQueued"`
	TimeStarted         time.Time `json:"timeStarted"`
}

type JobResultV3 struct {
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

type runFlags struct {
	workspace              string
	repository             string
	wait                   bool
	rtc                    bool
	description            string
	sourceData             string
	successTopics          []string
	failureTopics          []string
	publishedParameter     []string
	listPublishedParameter []string
	nodeManagerDirective   []string
	directive              []string
	queue                  string
	maxJobRuntime          int
	maxTimeInQueue         int
	maxTotalLifeTime       int
	outputType             string
	noHeaders              bool
}

func newRunCmd() *cobra.Command {
	f := runFlags{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a workspace on FME Server.",
		Long:  `Run a workspace on FME Server.`,
		Example: `
  # Submit a job asynchronously
  fmeflow run --repository Samples --workspace austinApartments.fmw
	
  # Submit a job and wait for it to complete
  fmeflow run --repository Samples --workspace austinApartments.fmw --wait
	
  # Submit a job to a specific queue and set a time to live in the queue
  fmeflow run --repository Samples --workspace austinApartments.fmw --queue Queue1 --max-time-in-queue 120
	
  # Submit a job and pass in a few published parameters
  fmeflow run --repository Samples --workspace austinDownload.fmw --published-parameter-list THEMES=railroad,airports --published-parameter COORDSYS=TX83-CF
	
  # Submit a job, wait for it to complete, and customize the output
  fmeflow run --repository Samples --workspace austinApartments.fmw --wait --output="custom-columns=Time Requested:.timeRequested,Time Started:.timeStarted,Time Finished:.timeFinished"
	
  # Upload a local file to use as the source data for the translation
  fmeflow run --repository Samples --workspace austinApartments.fmw --file Landmarks-edited.sqlite --wait`,
		Args: NoArgs,
		RunE: runRun(&f),
	}

	cmd.Flags().StringVar(&f.repository, "repository", "", "The name of the repository containing the workspace to run.")
	cmd.Flags().StringVar(&f.workspace, "workspace", "", "The name of the workspace to run.")
	cmd.Flags().BoolVar(&f.wait, "wait", false, "Submit job and wait for it to finish.")
	cmd.Flags().StringArrayVar(&f.publishedParameter, "published-parameter", []string{}, "Published parameters defined for this workspace. Specify as Key=Value. Can be passed in multiple times. For list parameters, use the --list-published-parameter flag.")
	cmd.Flags().StringArrayVar(&f.listPublishedParameter, "published-parameter-list", []string{}, "A List-type published parameters defined for this workspace. Specify as Key=Value1,Value2. Can be passed in multiple times.")
	cmd.Flags().StringVar(&f.sourceData, "file", "", "Upload a local file Source dataset to use to run the workspace. Note this causes the translation to run in synchonous mode whether the --wait flag is passed in or not. For v3 API only.")
	cmd.Flags().BoolVar(&f.rtc, "run-until-canceled", false, "Runs a job until it is explicitly canceled. The job will run again regardless of whether the job completed successfully, failed, or the server crashed or was shut down. For v3 API only.")
	cmd.Flags().StringVar(&f.description, "description", "", "Description of the request. For v3 API only.")
	cmd.Flags().StringArrayVar(&f.successTopics, "success-topic", []string{}, "Topics to notify when the job succeeds. Can be specified more than once.")
	cmd.Flags().StringArrayVar(&f.failureTopics, "failure-topic", []string{}, "Topics to notify when the job fails. Can be specified more than once.")
	cmd.Flags().StringArrayVar(&f.nodeManagerDirective, "node-manager-directive", []string{}, "Additional NM Directives, which can include client-configured keys, to pass to the notification service for custom use by subscriptions. Specify as Key=Value Can be passed in multiple times. For v3 API only.")
	cmd.Flags().StringArrayVar(&f.directive, "directive", []string{}, "Additional directives to pass to the job submission. Specify as Key=Value. Can be passed in multiple times. For v4 API only.")
	cmd.Flags().StringVar(&f.queue, "queue", "", "Queue of the job to submit.")
	cmd.Flags().StringVar(&f.queue, "tag", "", "The queue (job routing tag) for the request.")
	cmd.Flags().IntVar(&f.maxJobRuntime, "max-job-runtime", -1, "Time (in seconds) elapsed for a running job before it's cancelled. The minimum value is 1 second, values less than 1 second are ignored.")
	cmd.Flags().IntVar(&f.maxJobRuntime, "time-until-canceled", -1, "Time (in seconds) elapsed for a running job before it's cancelled. The minimum value is 1 second, values less than 1 second are ignored.")
	cmd.Flags().IntVar(&f.maxTimeInQueue, "max-time-in-queue", -1, "Time to live in the job queue (in seconds).")
	cmd.Flags().IntVar(&f.maxTimeInQueue, "time-to-live", -1, "Time to live in the job queue (in seconds).")
	cmd.Flags().IntVar(&f.maxTotalLifeTime, "max-total-life-time", -1, "Time to live including both time in the queue and run time (in seconds). The maximum value is 86400 and the minimum value is 1. For v4 API only.")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")

	// since there are a lot of flags in this command, using the sorting above with more important flags first seems helpful
	cmd.Flags().SortFlags = false

	cmd.MarkFlagRequired("repository")
	cmd.MarkFlagRequired("workspace")

	//node manager directives and run until canceled do not work with transactdata
	cmd.MarkFlagsMutuallyExclusive("file", "node-manager-directive")
	cmd.MarkFlagsMutuallyExclusive("file", "run-until-canceled")

	// deprecated flags can't be used with the equavalent new flags
	cmd.MarkFlagsMutuallyExclusive("tag", "queue")
	cmd.MarkFlagsMutuallyExclusive("time-until-canceled", "max-job-runtime")
	cmd.MarkFlagsMutuallyExclusive("time-to-live", "max-time-in-queue")

	// mark v3 deprecated flags
	cmd.Flags().MarkDeprecated("tag", "please use --queue instead")
	cmd.Flags().MarkDeprecated("time-until-canceled", "please use --max-job-runtime instead")
	cmd.Flags().MarkDeprecated("time-to-live", "please use --max-time-in-queue instead")

	return cmd
}

func runRun(f *runFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}
		// set up http
		client := &http.Client{
			// set a long timeout for jobs that are long running.
			// maybe this should be a parameter?
			Timeout: 604800 * time.Second,
		}

		if viper.GetInt("build") >= 26018 {
			var result JobResultV4
			var responseData []byte

			job := &JobRequestV4{}
			job.PublishedParameters = make(map[string]interface{})

			// get published parameters
			for _, parameter := range f.publishedParameter {
				this_parameter := strings.SplitN(parameter, "=", 2)
				job.PublishedParameters[this_parameter[0]] = this_parameter[1]
			}

			// get list published parameters
			for _, parameter := range f.listPublishedParameter {
				this_parameter := strings.SplitN(parameter, "=", 2)
				job.PublishedParameters[this_parameter[0]] = splitEscapedString(this_parameter[1], ',')
			}

			job.Directives = make(map[string]string)
			for _, directive := range f.directive {
				this_directive := strings.SplitN(directive, "=", 2)
				job.Directives[this_directive[0]] = this_directive[1]
			}

			job.SuccessTopics = append(job.SuccessTopics, f.successTopics...)
			job.FailureTopics = append(job.FailureTopics, f.failureTopics...)
			job.Queue = f.queue
			job.Repository = f.repository
			job.Workspace = f.workspace

			if f.maxJobRuntime > 0 {
				job.MaxJobRuntime = f.maxJobRuntime
			}

			if f.maxTimeInQueue > 0 {
				job.MaxTimeInQueue = f.maxTimeInQueue
			}

			if f.wait && f.maxTotalLifeTime > 0 && f.maxTotalLifeTime < 86401 {
				job.MaxTotalLifeTime = f.maxTotalLifeTime
			}

			jobJson, err := json.Marshal(job)
			if err != nil {
				return err
			}

			syncEndpoint := ""
			if f.wait {
				syncEndpoint = "/sync"
			}

			endpoint := "/fmeapiv4/jobs" + syncEndpoint

			request, err := buildFmeFlowRequest(endpoint, "POST", strings.NewReader(string(jobJson)))
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

			if !f.wait {
				var result JobId
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				} else {
					if !jsonOutput {
						fmt.Fprintln(cmd.OutOrStdout(), "Job submitted with id: "+strconv.Itoa(result.Id))
					} else {
						prettyJSON, err := prettyPrintJSON(responseData)
						if err != nil {
							return err
						}
						fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
					}
				}
			} else {
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				}
			}

			if f.wait {
				if f.outputType == "table" {
					t := table.NewWriter()
					t.SetStyle(defaultStyle)

					t.AppendHeader(table.Row{"ID", "Status", "Status Message", "Features Output"})

					t.AppendRow(table.Row{result.ID, result.Status, result.StatusMessage, result.FeatureOutputCount})

					if f.noHeaders {
						t.ResetHeaders()
					}
					fmt.Fprintln(cmd.OutOrStdout(), t.Render())

				} else if f.outputType == "json" {
					prettyJSON, err := prettyPrintJSON(responseData)
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
					if f.noHeaders {
						t.ResetHeaders()
					}
					fmt.Fprintln(cmd.OutOrStdout(), t.Render())

				} else {
					return errors.New("invalid output format specified")
				}

			}
			return nil

		} else {

			var result JobResultV3
			var responseData []byte

			if f.sourceData == "" {
				job := &JobRequestV3{}

				// get published parameters
				for _, parameter := range f.publishedParameter {
					this_parameter := strings.SplitN(parameter, "=", 2)
					var a SimpleParameter
					a.Name = this_parameter[0]
					a.Value = this_parameter[1]
					job.PublishedParameters = append(job.PublishedParameters, a)
				}

				// get list published parameters
				for _, parameter := range f.listPublishedParameter {
					this_parameter := strings.SplitN(parameter, "=", 2)
					var a ListParameter
					a.Name = this_parameter[0]
					// split on commas, unless they are escaped
					a.Value = splitEscapedString(this_parameter[1], ',')
					job.PublishedParameters = append(job.PublishedParameters, a)

				}

				// get node manager directives
				for _, directive := range f.nodeManagerDirective {
					this_directive := strings.Split(directive, "=")
					var a Directive
					a.Name = this_directive[0]
					a.Value = this_directive[1]
					job.NMDirectives.Directives = append(job.NMDirectives.Directives, a)
				}

				if f.maxJobRuntime != -1 {
					job.TMDirectives.Ttc = f.maxJobRuntime
				}
				if f.maxTimeInQueue != -1 {
					job.TMDirectives.TTL = f.maxTimeInQueue
				}

				if f.queue != "" {
					job.TMDirectives.Tag = f.queue
				}

				job.TMDirectives.Rtc = f.rtc

				// append slice to slice
				job.NMDirectives.SuccessTopics = append(job.NMDirectives.SuccessTopics, f.successTopics...)
				job.NMDirectives.FailureTopics = append(job.NMDirectives.FailureTopics, f.failureTopics...)

				if f.description != "" {
					job.TMDirectives.Description = f.description
				}

				jobJson, err := json.Marshal(job)
				if err != nil {
					return err
				}

				submitEndpoint := "submit"
				if f.wait {
					submitEndpoint = "transact"
				}

				endpoint := "/fmerest/v3/transformations/" + submitEndpoint + "/" + f.repository + "/" + f.workspace

				request, err := buildFmeFlowRequest(endpoint, "POST", strings.NewReader(string(jobJson)))
				if err != nil {
					return err
				}

				request.Header.Add("Content-Type", "application/json")

				response, err := client.Do(&request)

				if err != nil {
					return err
				} else if response.StatusCode != 200 && response.StatusCode != 202 {
					if response.StatusCode == 404 {
						return fmt.Errorf("%w: check that the specified workspace and repository exist", errors.New(response.Status))
					} else if response.StatusCode == 422 {
						return fmt.Errorf("%w: either job failed or published parameters are invalid", errors.New(response.Status))
					} else {
						return errors.New(response.Status)
					}
				}

				responseData, err = io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				if !f.wait {
					var result JobId
					if err := json.Unmarshal(responseData, &result); err != nil {
						return err
					} else {
						if !jsonOutput {
							fmt.Fprintln(cmd.OutOrStdout(), "Job submitted with id: "+strconv.Itoa(result.Id))
						} else {
							prettyJSON, err := prettyPrintJSON(responseData)
							if err != nil {
								return err
							}
							fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
						}
					}
				} else {
					if err := json.Unmarshal(responseData, &result); err != nil {
						return err
					}
				}
			} else {
				// we are uploading a source file, so we want to send the file in the body as octet stream, and parameters as url parameters
				file, err := os.Open(f.sourceData)
				if err != nil {
					return err
				}
				defer file.Close()

				endpoint := "/fmerest/v3/transformations/transactdata/" + f.repository + "/" + f.workspace
				request, err := buildFmeFlowRequest(endpoint, "POST", file)
				if err != nil {
					return err
				}

				q := request.URL.Query()

				if f.description != "" {
					q.Add("opt_description", f.description)
				}

				for _, topic := range f.successTopics {
					q.Add("opt_successtopics", topic)
				}

				for _, topic := range f.failureTopics {
					q.Add("opt_failuretopics", topic)
				}

				if f.description != "" {
					endpoint += "opt_description=" + f.description
				}

				if f.queue != "" {
					q.Add("opt_tag", f.queue)
				}

				if f.maxTimeInQueue != -1 {
					q.Add("opt_ttl", strconv.Itoa(f.maxTimeInQueue))
				}

				if f.maxJobRuntime != -1 {
					q.Add("opt_ttc", strconv.Itoa(f.maxJobRuntime))
				}

				for _, parameter := range f.publishedParameter {
					this_parameter := strings.SplitN(parameter, "=", 2)
					q.Add(this_parameter[0], this_parameter[1])
				}
				for _, parameter := range f.listPublishedParameter {
					this_parameter := strings.SplitN(parameter, "=", 2)
					this_list := splitEscapedString(this_parameter[1], ',')
					for _, item := range this_list {
						q.Add(this_parameter[0], item)
					}

				}

				request.URL.RawQuery = q.Encode()

				request.Header.Set("Content-Type", "application/octet-stream")

				response, err := client.Do(&request)
				if err != nil {
					return err
				} else if response.StatusCode != 200 {
					if response.StatusCode == 404 {
						return fmt.Errorf("%w: check that the specified workspace and repository exist", errors.New(response.Status))
					} else {
						return errors.New(response.Status)
					}

				}

				responseData, err = io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				}
			}

			// the transactdata endpoint only runs synchonously
			if f.wait || f.sourceData != "" {
				if f.outputType == "table" {
					t := table.NewWriter()
					t.SetStyle(defaultStyle)

					t.AppendHeader(table.Row{"ID", "Status", "Status Message", "Features Output"})

					t.AppendRow(table.Row{result.ID, result.Status, result.StatusMessage, result.NumFeaturesOutput})

					if f.noHeaders {
						t.ResetHeaders()
					}
					fmt.Fprintln(cmd.OutOrStdout(), t.Render())

				} else if f.outputType == "json" {
					prettyJSON, err := prettyPrintJSON(responseData)
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
					if f.noHeaders {
						t.ResetHeaders()
					}
					fmt.Fprintln(cmd.OutOrStdout(), t.Render())

				} else {
					return errors.New("invalid output format specified")
				}
			}
			return nil
		}

		return nil

	}
}

// split a string on delimiter, unless it is escaped
func splitEscapedString(s string, delimiter rune) []string {
	var result []string
	var builder strings.Builder
	var escaped bool
	for _, r := range s {
		if escaped {
			if r != '\\' && r != delimiter {
				builder.WriteRune('\\')
			}
			builder.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == delimiter && !escaped {
			result = append(result, builder.String())
			builder.Reset()
			continue
		}
		builder.WriteRune(r)
	}
	result = append(result, builder.String())
	return result
}
