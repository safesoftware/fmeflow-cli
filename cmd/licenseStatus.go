package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type LicenseStatusV3 struct {
	ExpiryDate       string `json:"expiryDate"`
	MaximumEngines   int    `json:"maximumEngines"`
	SerialNumber     string `json:"serialNumber"`
	IsLicenseExpired bool   `json:"isLicenseExpired"`
	IsLicensed       bool   `json:"isLicensed"`
	MaximumAuthors   int    `json:"maximumAuthors"`
}

type LicenseStatusV4 struct {
	Licensed       bool   `json:"licensed"`
	Expiration     string `json:"expiration"`
	MaximumEngines int    `json:"maximumEngines"`
	Expired        bool   `json:"expired"`
	SerialNumber   string `json:"serialNumber"`
	MaximumAuthors int    `json:"maximumAuthors"`
}

type licenseStatusFlags struct {
	outputType string
	noHeaders  bool
	apiVersion apiVersionFlag
}

var licenseStatusV4BuildThreshold = 23319

func newLicenseStatusCmd() *cobra.Command {
	f := licenseStatusFlags{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Retrieves status of the installed FME Server license.",
		Long:  `Retrieves status of the installed FME Server license.`,
		Args:  NoArgs,
		RunE:  licenseStatusRun(&f),
	}
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	return cmd
}

func licenseStatusRun(f *licenseStatusFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		// get build to decide if we should use v3 or v4
		// FME Server 2023.0+ and later can use v4. Otherwise fall back to v3
		if f.apiVersion == "" {
			fmeflowBuild := viper.GetInt("build")
			if fmeflowBuild < licenseStatusV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		var endpoint string
		if f.apiVersion == "v4" {
			endpoint = "/fmeapiv4/license/status"
		} else {
			endpoint = "/fmerest/v3/licensing/license/status"
		}

		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
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

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if f.apiVersion == "v4" {
			var result LicenseStatusV4
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			}

			if f.outputType == "table" {
				// output all values returned by the JSON in a table
				t := createTableWithDefaultColumns(result)

				if f.noHeaders {
					t.ResetHeaders()
				}
				fmt.Fprintln(cmd.OutOrStdout(), t.Render())
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

				marshalledItems := [][]byte{}
				mJson, err := json.Marshal(result)
				if err != nil {
					return err
				}
				marshalledItems = append(marshalledItems, mJson)

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
		} else {
			var result LicenseStatusV3
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			}

			if f.outputType == "table" {
				// output all values returned by the JSON in a table
				t := createTableWithDefaultColumns(result)

				if f.noHeaders {
					t.ResetHeaders()
				}
				fmt.Fprintln(cmd.OutOrStdout(), t.Render())
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

				marshalledItems := [][]byte{}
				mJson, err := json.Marshal(result)
				if err != nil {
					return err
				}
				marshalledItems = append(marshalledItems, mJson)

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
		}
		return nil

	}
}
