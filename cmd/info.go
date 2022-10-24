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

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Retrieves build, version and time information about FME Server",
	Long: `Retrieves build, version and time information about FME Server

Examples:

# Output FME Server information in a table
fmeserver info

# Output FME Server information in json
fmeserver info --json

# Output just the build string with no column headers
fmeserver info --output=custom-columns="BUILD:.build" --no-headers
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			outputType = "json"
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
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result FMEServerInfo
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if outputType == "table" {

				// output all values returned by the JSON in a table
				t := createTableWithDefaultColumns(result)

				if noHeaders {
					t.ResetHeaders()
				}
				fmt.Println(t.Render())

			} else if outputType == "json" {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Println(prettyJSON)
			} else if strings.HasPrefix(outputType, "custom-columns=") {
				// parse the columns and json queries
				columnsString := outputType[len("custom-columns="):]
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
				fmt.Println(t.Render())

			} else {
				return errors.New("invalid output format specified")
			}

		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	infoCmd.Flags().BoolVar(&noHeaders, "no-headers", false, "Don't print column headers")
}
