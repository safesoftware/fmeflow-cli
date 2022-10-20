package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/spf13/cobra"
)

type Engine struct {
	HostName                    string        `json:"hostName"`
	AssignedQueues              []string      `json:"assignedQueues"`
	ResultFailureCount          int           `json:"resultFailureCount"`
	InstanceName                string        `json:"instanceName"`
	RegistrationProperties      []string      `json:"registrationProperties"`
	EngineManagerNodeName       string        `json:"engineManagerNodeName"`
	MaxTransactionResultFailure int           `json:"maxTransactionResultFailure"`
	Type                        string        `json:"type"`
	BuildNumber                 int           `json:"buildNumber"`
	Platform                    string        `json:"platform"`
	ResultSuccessCount          int           `json:"resultSuccessCount"`
	MaxTransactionResultSuccess int           `json:"maxTransactionResultSuccess"`
	AssignedStreams             []interface{} `json:"assignedStreams"`
	TransactionPort             int           `json:"transactionPort"`
	CurrentJobID                int           `json:"currentJobID"`
}
type Engines struct {
	Offset     int      `json:"offset"`
	Limit      int      `json:"limit"`
	TotalCount int      `json:"totalCount"`
	Items      []Engine `json:"items"`
}

var count bool

// enginesCmd represents the engines command
var enginesCmd = &cobra.Command{
	Use:   "engines",
	Short: "Get information about the FME Engines",
	Long: `Gets information and status about FME Engines currently connected to FME Server
	
Examples:

# List all engines
fmeserver engines

# Output number of engines
fmeserver engines --count

# Output engines in json form
fmeserver engines --json

# Output just the names of the engines with no column headers
fmeserver engines --output=custom-columns=NAME:.instanceName --no-headers`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			outputType = "json"
		}
		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/transformations/engines", "GET", nil)
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

		var result Engines
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if count {
				// simply return the count of engines
				fmt.Println(result.TotalCount)
			} else if outputType == "table" { // output a table with some default fields selected
				t := table.NewWriter()
				t.SetStyle(defaultStyle)

				t.AppendHeader(table.Row{"Name", "Host", "Build", "Platform", "Type", "Current Job ID", "Registration Properties", "Queues"})

				for _, element := range result.Items {
					t.AppendRow(table.Row{element.InstanceName, element.HostName, element.BuildNumber, element.Platform, element.Type, element.CurrentJobID, element.RegistrationProperties, element.AssignedQueues})
				}
				if noHeaders {
					t.ResetHeaders()
				}
				fmt.Println(t.Render())
				// output the raw json but formatted
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
				for _, element := range result.Items {
					mJson, err := json.Marshal(element)
					if err != nil {
						return err
					}

					marshalledItems = append(marshalledItems, mJson)
				}

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
	rootCmd.AddCommand(enginesCmd)
	enginesCmd.Flags().StringVarP(&outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	enginesCmd.Flags().BoolVar(&noHeaders, "no-headers", false, "Don't print column headers")
	enginesCmd.Flags().BoolVar(&count, "count", false, "Prints the total count of engines.")
	enginesCmd.MarkFlagsMutuallyExclusive("output", "count")
	enginesCmd.MarkFlagsMutuallyExclusive("no-headers", "count")
	//enginesCmd.MarkFlagsMutuallyExclusive("json", "count")
}
