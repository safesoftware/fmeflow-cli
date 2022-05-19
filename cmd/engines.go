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
	"strconv"

	"github.com/spf13/cobra"
)

type Engines struct {
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	TotalCount int `json:"totalCount"`
	Items      []struct {
		HostName                    string        `json:"hostName"`
		AssignedQueues              []string      `json:"assignedQueues"`
		ResultFailureCount          int           `json:"resultFailureCount"`
		InstanceName                string        `json:"instanceName"`
		RegistrationProperties      []string      `json:"registrationProperties"`
		EngineManagerNodeName       string        `json:"engineManagerNodeName"`
		MaxTransactionResultFailure int           `json:"maxTransactionResultFailure"`
		Type                        string        `json:"type"`
		BuildNumber                 int           `json:"buildNumber"`
		Platform                    string        `json:"platform"`
		ResultSuccessCount          int           `json:"resultSuccessCount"`
		MaxTransactionResultSuccess int           `json:"maxTransactionResultSuccess"`
		AssignedStreams             []interface{} `json:"assignedStreams"`
		TransactionPort             int           `json:"transactionPort"`
		CurrentJobID                int           `json:"currentJobID"`
	} `json:"items"`
}

// enginesCmd represents the engines command
var enginesCmd = &cobra.Command{
	Use:   "engines",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/transformations/engines", "GET", nil)
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

		var result Engines
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Println("Total Engines: " + strconv.Itoa(result.TotalCount))

				for _, element := range result.Items {
					fmt.Println("------")
					fmt.Printf("Instance Name: %v\n", element.InstanceName)
					fmt.Printf("Engine Host: %v\n", element.HostName)
					fmt.Printf("FME Build: %v\n", element.BuildNumber)
					fmt.Printf("Platform: %v\n", element.Platform)
					fmt.Printf("Type: %v\n", element.Type)
					fmt.Printf("Current Job ID: %v\n", element.CurrentJobID)
					fmt.Printf("RegistrationProperties:\n")
					for _, property := range element.RegistrationProperties {
						fmt.Printf("\t%v\n", property)
					}
					fmt.Printf("Queues:\n")
					for _, queue := range element.AssignedQueues {
						fmt.Printf("\t%v\n", queue)
					}
				}

			} else {
				fmt.Println(string(responseData))
			}

		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(enginesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// enginesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// enginesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}