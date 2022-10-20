package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// requestStatusCmd represents the requestStatus command
var requestStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check status of license request",
	Long: `Check the status of a license request.
    
Example:
fmeserver license request status`,
	RunE: func(cmd *cobra.Command, args []string) error {
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
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result RequestStatus
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				// output all values returned by the JSON in a table
				v := reflect.ValueOf(result)
				typeOfS := v.Type()
				header := table.Row{}
				row := table.Row{}
				for i := 0; i < v.NumField(); i++ {
					header = append(header, convertCamelCaseToTitleCase(typeOfS.Field(i).Name))
					row = append(row, v.Field(i).Interface())
				}

				t := table.NewWriter()
				t.SetStyle(defaultStyle)

				t.AppendHeader(header)
				t.AppendRow(row)

				if noHeaders {
					t.ResetHeaders()
				}
				fmt.Println(t.Render())
			} else {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Println(prettyJSON)
			}

		}
		return nil

	},
}

func init() {
	requestCmd.AddCommand(requestStatusCmd)
}
