package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var ready bool

type Healthcheck struct {
	Status string `json:"status"`
}

// healthcheckCmd represents the healthcheck command
var healthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Retrieves the health status of FME Server",
	Long:  `Retrieves the health status of FME Server. The health status is normal if the FME Server REST API is responsive. Note that this endpoint does not require authentication. Load balancer or other systems can monitor FME Server using this endpoint without supplying token or password credentials.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			outputType = "json"
		}

		// set up http
		client := &http.Client{}

		endpoint := "/fmerest/v3/healthcheck"
		if ready {
			endpoint += "?ready=true"
		}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest(endpoint, "GET", nil)
		if err != nil {
			return err
		}
		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 200 {
			return errors.New(response.Status)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result Healthcheck
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
					header = append(header, typeOfS.Field(i).Name)
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
	rootCmd.AddCommand(healthcheckCmd)
	healthcheckCmd.Flags().BoolVar(&ready, "ready", false, "The health check will report the status of FME Server if it is ready to process jobs.")
	healthcheckCmd.Flags().StringVarP(&outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	healthcheckCmd.Flags().BoolVar(&noHeaders, "no-headers", false, "Don't print column headers")
}
