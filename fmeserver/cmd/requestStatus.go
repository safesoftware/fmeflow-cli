/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

// requestStatusCmd represents the requestStatus command
var requestStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check status of license request",
	Long: `Check the status of a license request.
	
Example:
fmeserver license request status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/licensing/request/status", "GET", nil)
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

		var result RequestStatus
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Println(result.Status)
				fmt.Println(result.Message)
			} else {
				fmt.Println(string(responseData))
			}

		}
		return nil

	},
}

func init() {
	requestCmd.AddCommand(requestStatusCmd)
}
