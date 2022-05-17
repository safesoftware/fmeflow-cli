/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

type RefreshStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var refreshWait bool

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		request, err := buildFmeServerRequest("/fmerest/v3/licensing/refresh", "POST", nil)
		if err != nil {
			return err
		}

		//request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 202 {
			return errors.New(response.Status)
		}

		fmt.Println("License Refresh Successfully sent.")

		if refreshWait {
			// check the license refresh status until it is finished
			complete := false
			for {
				fmt.Print(".")
				time.Sleep(1 * time.Second)
				// call the status endpoint to see if it is finished
				request, err := buildFmeServerRequest("/fmerest/v3/licensing/refresh/status", "GET", nil)
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

				var result RefreshStatus
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				} else if result.Status != "REQUESTING" {
					complete = true
					fmt.Println(result.Message)
				}

				if complete {
					break
				}
			}
		}

		return nil
	},
}

func init() {
	licenseCmd.AddCommand(refreshCmd)
	refreshCmd.Flags().BoolVar(&refreshWait, "wait", false, "Wait for licensing refresh to finish")
}
