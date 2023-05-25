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

// generated using https://mholt.github.io/json-to-go/
type migrationTasks struct {
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	TotalCount int `json:"totalCount"`
	Items      []migrationTask
}

type migrationTask struct {
	DisableProjectItems  bool      `json:"disableProjectItems"`
	Result               string    `json:"result"`
	ImportMode           string    `json:"importMode"`
	ProjectsImportMode   string    `json:"projectsImportMode"`
	PauseNotifications   bool      `json:"pauseNotifications"`
	ID                   int       `json:"id"`
	Type                 string    `json:"type"`
	UserName             string    `json:"userName"`
	ContentType          string    `json:"contentType"`
	StartDate            time.Time `json:"startDate"`
	FinishedDate         time.Time `json:"finishedDate"`
	Status               string    `json:"status"`
	ExcludeSensitiveInfo bool      `json:"excludeSensitiveInfo"`
	FailureTopic         string    `json:"failureTopic"`
	SuccessTopic         string    `json:"successTopic"`
	PackageName          string    `json:"packageName"`
	PackagePath          string    `json:"packagePath"`
	ProjectNames         []string  `json:"projectNames"`
	ResourceName         string    `json:"resourceName"`
}

type migrationTasksFlags struct {
	migrationTaskId   int
	migrationTaskLog  bool
	migrationTaskFile string
	outputType        string
	noHeaders         bool
}

func newMigrationTasksCmd() *cobra.Command {
	f := migrationTasksFlags{}
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Retrieves the records for migration tasks.",
		Long:  "Retrieves the records for migration tasks. Get all migration tasks or for a specific task by passing in the id.",
		Example: `
  # Get all migration tasks
  fmeflow migration tasks
	
  # Get all migration tasks in json
  fmeflow migration tasks --json
	
  # Get the migration task for a given id
  fmeflow migration tasks --id 1
	
  # Output the migration log for a given id to the console
  fmeflow migration tasks --id 1 --log
	
  # Output the migration log for a given id to a local file
  fmeflow migration tasks --id 1 --log --file my-backup-log.txt
	
  # Output just the start and end time of the a given id
  fmeflow migration tasks --id 1 --output="custom-columns=Start Time:.startDate,End Time:.finishedDate"`,
		Args: NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			if f.migrationTaskLog {
				cmd.MarkFlagsRequiredTogether("id", "log")
			}
		},
		RunE: migrationTasksRun(&f),
	}

	cmd.Flags().IntVar(&f.migrationTaskId, "id", -1, "Retrieves the record for a migration task according to the given ID.")
	cmd.Flags().BoolVar(&f.migrationTaskLog, "log", false, "Downloads the log file of a migration task.")
	cmd.Flags().StringVar(&f.migrationTaskFile, "file", "", "File to save the log to.")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")

	return cmd
}

func migrationTasksRun(f *migrationTasksFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		// set up http
		client := &http.Client{}

		var outputTasks []migrationTask

		if !f.migrationTaskLog { // output one or more tasks
			var responseData []byte
			if f.migrationTaskId == -1 {
				request, err := buildFmeFlowRequest("/fmerest/v3/migration/tasks", "GET", nil)
				if err != nil {
					return err
				}
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

				var result migrationTasks
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				} else {
					outputTasks = result.Items
				}
			} else {
				endpoint := "/fmerest/v3/migration/tasks/id/" + strconv.Itoa(f.migrationTaskId)
				request, err := buildFmeFlowRequest(endpoint, "GET", nil)
				if err != nil {
					return err
				}
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

				var result migrationTask
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				} else {
					outputTasks = []migrationTask{result}
				}
			}

			if f.outputType == "table" {
				t := table.NewWriter()
				t.SetStyle(defaultStyle)

				t.AppendHeader(table.Row{"ID", "Type", "Username", "Start Time", "End Time", "Status"})

				for _, element := range outputTasks {
					t.AppendRow(table.Row{element.ID, element.Type, element.UserName, element.StartDate, element.FinishedDate, element.Status})
				}
				if f.noHeaders {
					t.ResetHeaders()
				}
				fmt.Fprintln(cmd.OutOrStdout(), t.Render())
				// output the raw json but formatted
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
				for _, element := range outputTasks {
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

		} else if f.migrationTaskId != -1 && f.migrationTaskLog {
			endpoint := "/fmerest/v3/migration/tasks/id/" + strconv.Itoa(f.migrationTaskId) + "/log"
			request, err := buildFmeFlowRequest(endpoint, "GET", nil)
			if err != nil {
				return err
			}

			request.Header.Add("Accept", "application/octet-stream")
			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != 200 {
				return errors.New(response.Status)
			}

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			if f.migrationTaskFile == "" {
				fmt.Fprintln(cmd.OutOrStdout(), string(responseData))
			} else {
				// Create the output file
				out, err := os.Create(f.migrationTaskFile)
				if err != nil {
					return err
				}
				defer out.Close()

				// use Copy so that it doesn't store the entire file in memory
				_, err = io.Copy(out, strings.NewReader(string(responseData)))
				if err != nil {
					return err
				}

				fmt.Fprintln(cmd.OutOrStdout(), "Log file downloaded to "+f.migrationTaskFile)
			}

		}

		return nil
	}
}
