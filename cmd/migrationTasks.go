/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// generated using https://mholt.github.io/json-to-go/
type migrationTasks struct {
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	TotalCount int `json:"totalCount"`
	Items      []struct {
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
	} `json:"items"`
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
		client := &http.Client{}

		if migrationTaskId == -1 && !migrationTaskLog {

			request, err := buildFmeServerRequest("/fmerest/v3/migration/tasks", "GET", nil)
			if err != nil {
				return err
			}
			response, err := client.Do(&request)
			if err != nil {
				return err
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result migrationTasks
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			} else {
				if !jsonOutput {
					// TODO: Figure out a nice way of outputting all the tasks. For now always output json

					// output all values returned by the JSON in a table
					/*v := reflect.ValueOf(result)
					typeOfS := v.Type()

					for i := 0; i < v.NumField(); i++ {
						fmt.Printf("%s:\t%v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
					}*/
					//fmt.Printf("%+v\n", result)
					fmt.Println(string(responseData))
				} else {
					fmt.Println(string(responseData))
				}
			}
		} else if migrationTaskId != -1 && !migrationTaskLog {
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
				if !jsonOutput {
					// output all values returned by the JSON in a table
					v := reflect.ValueOf(result)
					typeOfS := v.Type()

					for i := 0; i < v.NumField(); i++ {
						fmt.Printf("%s:\t%v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
					}
					//fmt.Printf("%+v\n", result)
					//fmt.Println(string(responseData))
				} else {
					fmt.Println(string(responseData))
				}
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
}
