package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type MachineKey struct {
	MachineKey string `json:"machineKey"`
}

// machinekeyCmd represents the machinekey command
var machinekeyCmd = &cobra.Command{
	Use:   "machinekey",
	Short: "Retrieves machine key of the machine running FME Server.",
	Long:  `Retrieves machine key of the machine running FME Server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			outputType = "json"
		}
		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/licensing/machinekey", "GET", nil)
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

		var result MachineKey
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if outputType == "table" {
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

			} else if outputType == "json" {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Println(prettyJSON)
			} else {
				return errors.New("invalid output format specified")
			}
		}
		return nil
	},
}

func init() {
	licenseCmd.AddCommand(machinekeyCmd)
	machinekeyCmd.Flags().StringVarP(&outputType, "output", "o", "table", "Specify the output type. Should be one of table or json")
	machinekeyCmd.Flags().BoolVar(&noHeaders, "no-headers", false, "Don't print column headers")
}
