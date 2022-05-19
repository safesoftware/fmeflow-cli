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
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var firstName string
var lastName string
var email string
var serialNumber string
var company string
var wait bool

type RequestStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request a license from the FME Server licensing server",
	Long: `Request a license file from the FME Server licensing server. First name, Last name and email are required for requesting a license file.
If no serial number is passed in, a trial license will be requested.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// add mandatory values
		data := url.Values{
			"firstName": {firstName},
			"lastName":  {lastName},
			"email":     {email},
		}

		// add optional values
		if serialNumber != "" {
			data.Add("serialNumber", serialNumber)
		}
		if company != "" {
			data.Add("company", company)
		}

		request, err := buildFmeServerRequest("/fmerest/v3/licensing/request", "POST", strings.NewReader(data.Encode()))
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 202 {
			return errors.New(response.Status)
		}

		if !jsonOutput {
			fmt.Println("License Request Successfully sent.")
		} else {
			if !wait {
				fmt.Println("{}")
			}
		}

		if wait {
			// check the license status until it is finished
			complete := false
			for {
				if !jsonOutput {
					fmt.Print(".")
				}

				time.Sleep(1 * time.Second)
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
				} else if result.Status != "REQUESTING" {
					complete = true
					if !jsonOutput {
						fmt.Println(result.Message)
					} else {
						fmt.Println(string(responseData))
					}
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
	licenseCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringVar(&firstName, "first-name", "", "First name to use for license request.")
	requestCmd.Flags().StringVar(&lastName, "last-name", "", "Last name to use for license request.")
	requestCmd.Flags().StringVar(&email, "email", "", "Email address for license request.")
	requestCmd.Flags().StringVar(&serialNumber, "serial-number", "", "Serial Number for the license request.")
	requestCmd.Flags().StringVar(&company, "company", "", "Company for the licensing request")
	requestCmd.Flags().BoolVar(&wait, "wait", false, "Wait for licensing request to finish")
	requestCmd.MarkFlagRequired("first-name")
	requestCmd.MarkFlagRequired("last-name")
	requestCmd.MarkFlagRequired("email")
}
