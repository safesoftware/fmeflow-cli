package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	DisableProjectItems bool      `json:"disableProjectItems"`
	Result              string    `json:"result"`
	ImportMode          string    `json:"importMode"`
	ProjectsImportMode  string    `json:"projectsImportMode"`
	PauseNotifications  bool      `json:"pauseNotifications"`
	ID                  int       `json:"id"`
	Type                string    `json:"type"`
	UserName            string    `json:"userName"`
	ContentType         string    `json:"contentType"`
	StartDate           time.Time `json:"startDate"`
	FinishedDate        time.Time `json:"finishedDate"`
	Status              string    `json:"status"`
}

var migrationTaskId int
var migrationTaskLog bool
var migrationTaskFile string

// migrationTasksCmd represents the migrationTasks command
var migrationTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Retrieves the records for all migration tasks.",
	Long:  `Retrieves the records for all migration tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		// --json overrides --output
		if jsonOutput {
			outputType = "json"
		}
		client := &http.Client{}

		var outputTasks []migrationTask

		if !migrationTaskLog { // output one or more tasks
			var responseData []byte
			if migrationTaskId == -1 {
				request, err := buildFmeServerRequest("/fmerest/v3/migration/tasks", "GET", nil)
				if err != nil {
					return err
				}
				response, err := client.Do(&request)
				if err != nil {
					return err
				}

				responseData, err = ioutil.ReadAll(response.Body)
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
				endpoint := "/fmerest/v3/migration/tasks/id/" + strconv.Itoa(migrationTaskId)
				request, err := buildFmeServerRequest(endpoint, "GET", nil)
				if err != nil {
					return err
				}
				response, err := client.Do(&request)
				if err != nil {
					return err
				} else if response.StatusCode != 200 {
					return errors.New(response.Status)
				}

				responseData, err := ioutil.ReadAll(response.Body)
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

			if outputType == "table" {
				t := table.NewWriter()
				t.SetStyle(defaultStyle)

				t.AppendHeader(table.Row{"ID", "Type", "Username", "Start Time", "End Time", "Status"})

				for _, element := range outputTasks {
					t.AppendRow(table.Row{element.ID, element.Type, element.UserName, element.StartDate, element.FinishedDate, element.Status})
				}
				if noHeaders {
					t.ResetHeaders()
				}
				fmt.Println(t.Render())
				// output the raw json but formatted
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
				if noHeaders {
					t.ResetHeaders()
				}
				fmt.Println(t.Render())
			} else {
				return errors.New("invalid output format specified")
			}

		} else if migrationTaskId != -1 && migrationTaskLog {
			endpoint := "/fmerest/v3/migration/tasks/id/" + strconv.Itoa(migrationTaskId) + "/log"
			request, err := buildFmeServerRequest(endpoint, "GET", nil)
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

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}

			if migrationTaskFile == "" {
				fmt.Println(string(responseData))
			} else {
				// Create the output file
				out, err := os.Create(migrationTaskFile)
				if err != nil {
					return err
				}
				defer out.Close()

				// use Copy so that it doesn't store the entire file in memory
				_, err = io.Copy(out, strings.NewReader(string(responseData)))
				if err != nil {
					return err
				}

				fmt.Println("Log file downloaded to " + migrationTaskFile)
			}

		}

		return nil
	},
}

func init() {
	migrationCmd.AddCommand(migrationTasksCmd)

	migrationTasksCmd.Flags().IntVar(&migrationTaskId, "id", -1, "Retrieves the record for a migration task according to the given ID.")
	migrationTasksCmd.Flags().BoolVar(&migrationTaskLog, "log", false, "Downloads the log file of a migration task.")
	migrationTasksCmd.Flags().StringVar(&migrationTaskFile, "file", "", "File to save the log to.")
	migrationTasksCmd.Flags().StringVarP(&outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	migrationTasksCmd.Flags().BoolVar(&noHeaders, "no-headers", false, "Don't print column headers")
	migrationTasksCmd.MarkFlagsRequiredTogether("id", "log")
}
