/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var outputLicenseFile string

// requestfileCmd represents the requestfile command
var requestfileCmd = &cobra.Command{
	Use:   "requestfile",
	Short: "Generates a JSON file for requesting a FME Server license file.",
	Long:  `Generates a JSON file for requesting a FME Server license file.`,
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

		request, err := buildFmeServerRequest("/fmerest/v3/licensing/requestfile", "POST", strings.NewReader(data.Encode()))
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Accept", "application/octet-stream")

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 200 {
			return errors.New(response.Status)
		}

		// read the body which should be the contents of the file
		d, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if outputLicenseFile != "" {
			tmpfile, err := os.Create(outputLicenseFile)
			if err != nil {
				return err
			}
			defer tmpfile.Close()
			tmpfile.Write(d)
		} else {
			fmt.Println(string(d))
		}
		return nil
	},
}

func init() {
	licenseCmd.AddCommand(requestfileCmd)

	requestfileCmd.Flags().StringVar(&firstName, "first-name", "", "First name to use for license request.")
	requestfileCmd.Flags().StringVar(&lastName, "last-name", "", "Last name to use for license request.")
	requestfileCmd.Flags().StringVar(&email, "email", "", "Email address for license request.")
	requestfileCmd.Flags().StringVar(&serialNumber, "serial-number", "", "Serial Number for the license request.")
	requestfileCmd.Flags().StringVar(&company, "company", "", "Company for the licensing request")
	requestfileCmd.Flags().StringVar(&outputLicenseFile, "file", "", "Path to file to output to.")
	requestfileCmd.MarkFlagRequired("first-name")
	requestfileCmd.MarkFlagRequired("last-name")
	requestfileCmd.MarkFlagRequired("email")
}
