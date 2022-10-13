package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/spf13/cobra"
)

type LicenseStatus struct {
	ExpiryDate        string `json:"expiryDate"`
	MaximumEngines    int    `json:"maximumEngines"`
	SerialNumber      string `json:"serialNumber"`
	IsLicensedExpired bool   `json:"isLicenseExpired"`
	IsLicensed        bool   `json:"isLicensed"`
	MaximumAuthors    int    `json:"maximumAuthors"`
}

// licenseStatusCmd represents the licenseStatus command
var licenseStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Retrieves status of the installed FME Server license.",
	Long:  `Retrieves status of the installed FME Server license.`,
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
	licenseCmd.AddCommand(licenseStatusCmd)
}
