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

type licenseRefreshStatusFlags struct {
	outputType string
	noHeaders  bool
	apiVersion apiVersionFlag
}

func newRefreshStatusCmd() *cobra.Command {
	f := licenseRefreshStatusFlags{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check the status of a license refresh request.",
		Long:  "Check the status of a license refresh request.",
		Example: `
	# Output the license refresh status as a table
	fmeflow license refresh status
	
	# Output the license refresh status in json
	fmeflow license refresh status --json
	
	# Output just the status message
	fmeflow license refresh status --output custom-columns=STATUS:.status --no-headers`,
		Args: NoArgs,
		RunE: refreshStatusRun(&f),
	}
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	return cmd

}

func refreshStatusRun(f *licenseRefreshStatusFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		// get build to decide if we should use v3 or v4
		// FME Server 2023.0+ and later can use v4. Otherwise fall back to v3
		if f.apiVersion == "" {
			fmeflowBuild := viper.GetInt("build")
			if fmeflowBuild < refreshV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		var statusEndpoint string
		if f.apiVersion == "v4" {
			statusEndpoint = "/fmeapiv4/license/refresh/status"
		} else {
			statusEndpoint = "/fmerest/v3/licensing/refresh/status"
		}

		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeFlowRequest(statusEndpoint, "GET", nil)
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

		var result RefreshStatus
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if f.outputType == "table" {
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

				// we have to marshal the Items array, then create an array of marshalled items
				// to pass to the creation of the table.
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
