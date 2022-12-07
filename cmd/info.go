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

type FMEServerInfo struct {
	CurrentTime       string `json:"currentTime"`
	LicenseManagement bool   `json:"licenseManagement"`
	Build             string `json:"build"`
	TimeZone          string `json:"timeZone"`
	Version           string `json:"version"`
}

type infoFlags struct {
	outputType string
	noHeaders  bool
}

func newInfoCmd() *cobra.Command {
	f := infoFlags{}
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Retrieves build, version and time information about FME Server",
		Long:  "Retrieves build, version and time information about FME Server",
		Example: `
  # Output FME Server information in a table
  fmeserver info

  # Output FME Server information in json
  fmeserver info --json

  # Output just the build string with no column headers
  fmeserver info --output=custom-columns="BUILD:.build" --no-headers
	`,
		Args: NoArgs,
		RunE: infoRun(&f),
	}
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	return cmd
}

func infoRun(f *infoFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/info", "GET", nil)
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

		var result FMEServerInfo
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if f.outputType == "table" {

				// output all values returned by the JSON in a table
				t := createTableWithDefaultColumns(result)

				if noHeaders {
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
				if noHeaders {
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
