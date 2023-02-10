package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

type licenseRequestStatusFlags struct {
	outputType string
	noHeaders  bool
}

func newLicenseRequestStatusCmd() *cobra.Command {
	f := licenseRequestStatusFlags{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check the status of a license request.",
		Long:  "Check the status of a license request.",
		Example: `
	# Output the license request status as a table
	fmeserver license request status
	
	# Output the license Request status in json
	fmeserver license request status --json
	
	# Output just the status message
	fmeserver license request status --output custom-columns=STATUS:.status --no-headers`,
		Args: NoArgs,
		RunE: licenseRequestStatusRun(&f),
	}
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	return cmd
}

func licenseRequestStatusRun(f *licenseRequestStatusFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

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
		} else if response.StatusCode != 200 {
			return errors.New(response.Status)
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result RequestStatus
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
