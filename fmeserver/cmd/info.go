/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/spf13/cobra"
)

type FMEServerInfo struct {
	CurrentTime       string `json:"currentTime"`
	LicenseManagement bool   `json:"licenseManagement"`
	Build             string `json:"build"`
	TimeZone          string `json:"timeZone"`
	Version           int    `json:"version"`
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Retrieves build, version and time information about FME Server",
	Long:  `Retrieves build, version and time information about FME Server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/licensing/license/status", "GET", nil)
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

		var result LicenseStatus
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
			} else {
				fmt.Println(string(responseData))
			}

		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
