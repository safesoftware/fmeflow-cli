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

type JobRequest struct {
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

type runFlags struct {
	runWorkspace            string
	runRepository           string
	runWait                 bool
	runRtc                  bool
	runTtc                  int
	runTtl                  int
	runTag                  string
	runDescription          string
	runSourceData           string
	runSuccessTopics        []string
	runFailureTopics        []string
	runPublishedParameter   []string
	runNodeManagerDirective []string
	outputType              string
	noHeaders               bool
}

func newRunCmd() *cobra.Command {
	f := runFlags{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a workspace on FME Server",
		Long: `Run a workspace on FME Server
		
	Examples:
	# Submit a job asynchronously
	fmeserver run --repository Samples --workspace austinApartments.fmw
	
	# Submit a job and wait for it to complete
	fmeserver run --repository Samples --workspace austinApartments.fmw --wait
	
	# Submit a job to a specific queue and set a time to live in the queue
	fmeserver run --repository Samples --workspace austinApartments.fmw --tag Queue1 --time-to-live 120
	
	# Submit a job and pass in a few published parameters
	fmeserver run --repository Samples --workspace austinDownload.fmw --published-parameter THEMES=railroad,airports --published-parameter COORDSYS=TX83-CF
	
	# Submit a job, wait for it to complete, and customize the output
	fmeserver run --repository Samples --workspace austinApartments.fmw --wait --output="custom-columns=Time Requested:.timeRequested,Time Started:.timeStarted,Time Finished:.timeFinished"
	
	# Upload a local file to use as the source data for the translation
	fmeserver run --repository Samples --workspace austinApartments.fmw --file Landmarks-edited.sqlite --wait`,
		Args: NoArgs,
		RunE: runRun(&f),
	}

	cmd.Flags().StringVar(&f.runRepository, "repository", "", "The name of the repository containing the workspace to run.")
	cmd.Flags().StringVar(&f.runWorkspace, "workspace", "", "The name of the workspace to run.")
	cmd.Flags().BoolVar(&f.runWait, "wait", false, "Submit job and wait for it to finish.")
	cmd.Flags().BoolVar(&f.runRtc, "run-until-canceled", false, "Runs a job until it is explicitly canceled. The job will run again regardless of whether the job completed successfully, failed, or the server crashed or was shut down.")
	cmd.Flags().IntVar(&f.runTtc, "time-until-canceled", -1, "Time (in seconds) elapsed for a running job before it's cancelled. The minimum value is 1 second, values less than 1 second are ignored.")
	cmd.Flags().IntVar(&f.runTtl, "time-to-live", -1, "Time to live in the job queue (in seconds)")
	cmd.Flags().StringVar(&f.runTag, "tag", "", "The job routing tag for the request")
	cmd.Flags().StringVar(&f.runDescription, "description", "", "Description of the request.")
	cmd.Flags().StringVar(&f.runSourceData, "file", "", "Upload a local file Source dataset to use to run the workspace.")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")

	cmd.Flags().StringArrayVar(&f.runSuccessTopics, "success-topic", []string{}, "Topics to notify when the job succeeds. Can be specified more than once.")
	cmd.Flags().StringArrayVar(&f.runFailureTopics, "failure-topic", []string{}, "Topics to notify when the job fails. Can be specified more than once.")
	cmd.Flags().StringArrayVar(&f.runPublishedParameter, "published-parameter", []string{}, "Workspace published parameters defined for this job. Specify as Key=Value. Can be passed in multiple times. For list parameters, specify as Key=Value1,Value2. This means parameter values can't contain = or , at the moment. That should probably be fixed.")
	cmd.Flags().StringArrayVar(&f.runNodeManagerDirective, "node-manager-directive", []string{}, "Additional NM Directives, which can include client-configured keys, to pass to the notification service for custom use by subscriptions. Specify as Key=Value Can be passed in multiple times.")

	cmd.MarkFlagRequired("repository")
	cmd.MarkFlagRequired("workspace")

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

		var result JobResult
		var responseData []byte

		if f.runSourceData == "" {
			job := &JobRequest{}

			// get published parameters
			for _, parameter := range f.runPublishedParameter {
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
			for _, directive := range f.runNodeManagerDirective {
				this_directive := strings.Split(directive, "=")
				var a Directive
				a.Name = this_directive[0]
				a.Value = this_directive[1]
				job.NMDirectives.Directives = append(job.NMDirectives.Directives, a)
			}

			if f.runTtc != -1 {
				job.TMDirectives.Ttc = f.runTtc
			}
			if f.runTtl != -1 {
				job.TMDirectives.TTL = f.runTtl
			}

			job.TMDirectives.Rtc = f.runRtc

			// append slice to slice
			job.NMDirectives.SuccessTopics = append(job.NMDirectives.SuccessTopics, f.runSuccessTopics...)
			job.NMDirectives.FailureTopics = append(job.NMDirectives.FailureTopics, f.runFailureTopics...)

			if f.runDescription != "" {
				job.TMDirectives.Description = f.runDescription
			}

			jobJson, err := json.Marshal(job)
			if err != nil {
				return err
			}

			submitEndpoint := "submit"
			if f.runWait {
				submitEndpoint = "transact"
			}

			endpoint := "/fmerest/v3/transformations/" + submitEndpoint + "/" + f.runRepository + "/" + f.runWorkspace

			request, err := buildFmeServerRequest(endpoint, "POST", strings.NewReader(string(jobJson)))
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

			if !f.runWait {
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
			file, err := os.Open(f.runSourceData)
			if err != nil {
				return err
			}
			defer file.Close()

			endpoint := "/fmerest/v3/transformations/transactdata/" + f.runRepository + "/" + f.runWorkspace
			request, err := buildFmeServerRequest(endpoint, "POST", file)
			if err != nil {
				return err
			}

			q := request.URL.Query()

			if f.runDescription != "" {
				q.Add("opt_description", f.runDescription)
			}

			for _, topic := range f.runSuccessTopics {
				q.Add("opt_successtopics", topic)
			}

			for _, topic := range f.runFailureTopics {
				q.Add("opt_failuretopics", topic)
			}

			if f.runDescription != "" {
				endpoint += "opt_description=" + f.runDescription
			}

			if f.runTag != "" {
				q.Add("opt_tag", f.runTag)
			}

			if f.runTtl != -1 {
				q.Add("opt_ttl", strconv.Itoa(f.runTtl))
			}

			if f.runTtc != -1 {
				q.Add("opt_ttc", strconv.Itoa(f.runTtc))
			}

			// TODO: I'm not sure this is the correct way to pass published parameters in the query string
			for _, parameter := range f.runPublishedParameter {
				this_parameter := strings.Split(parameter, "=")
				q.Add(this_parameter[0], this_parameter[1])
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

		if f.runWait {
			if f.outputType == "table" {
				t := table.NewWriter()
				t.SetStyle(defaultStyle)

				t.AppendHeader(table.Row{"ID", "Status", "Status Message", "Features Output"})

				t.AppendRow(table.Row{result.ID, result.Status, result.StatusMessage, result.NumFeaturesOutput})

				if f.noHeaders {
					t.ResetHeaders()
				}
				fmt.Println(t.Render())

			} else if f.outputType == "json" {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Println(prettyJSON)
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
				fmt.Println(t.Render())

			} else {
				return errors.New("invalid output format specified")
			}
		}
		return nil
	}
}
