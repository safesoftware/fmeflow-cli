package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type EngineV4 struct {
	Name                   string   `json:"name"`
	Hostname               string   `json:"hostname"`
	EngineManagerHostname  string   `json:"engineManagerHostname"`
	Platform               string   `json:"platform"`
	CurrentJobID           int      `json:"currentJobID"`
	BuildNumber            int      `json:"buildNumber"`
	Type                   string   `json:"type"`
	State                  string   `json:"state"`
	AssignedQueues         []string `json:"assignedQueues"`
	RegistrationProperties []string `json:"registrationProperties"`
	HostProperties         struct {
		PhysicalMemory int `json:"physicalMemory"`
		ProcessorCount int `json:"processorCount"`
	} `json:"hostProperties"`
}

type EngineV3 struct {
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

type EnginesV4 struct {
	Offset     int        `json:"offset"`
	Limit      int        `json:"limit"`
	TotalCount int        `json:"totalCount"`
	Items      []EngineV4 `json:"items"`
}

type EnginesV3 struct {
	Offset     int        `json:"offset"`
	Limit      int        `json:"limit"`
	TotalCount int        `json:"totalCount"`
	Items      []EngineV3 `json:"items"`
}

type engineFlags struct {
	count      bool
	outputType string
	noHeaders  bool
	apiVersion apiVersionFlag
}

// enginesCmd represents the engines command
func newEnginesCmd() *cobra.Command {
	f := engineFlags{}
	cmd := &cobra.Command{
		Use:   "engines",
		Short: "Get information about the FME Engines",
		Long:  "Gets information and status about FME Engines currently connected to FME Server",
		Example: `
  # List all engines
  fmeflow engines
	
  # Output number of engines
  fmeflow engines --count
	
  # Output engines in json form
  fmeflow engines --json
	
  # Output just the names of the engines with no column headers
  fmeflow engines --output=custom-columns=NAME:.instanceName --no-headers`,
		Args: NoArgs,
		RunE: enginesRun(&f),
	}
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.Flags().BoolVar(&f.count, "count", false, "Prints the total count of engines.")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.MarkFlagsMutuallyExclusive("output", "count")
	cmd.MarkFlagsMutuallyExclusive("no-headers", "count")
	//enginesCmd.MarkFlagsMutuallyExclusive("json", "count")
	return cmd

}

func enginesRun(f *engineFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		if f.apiVersion == "" {
			if viper.GetInt("build") < projectsV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		// set up http
		client := &http.Client{}

		if f.apiVersion == "v4" {
			request, err := buildFmeFlowRequest("/fmeapiv4/engines", "GET", nil)
			if err != nil {
				return err
			}
			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != http.StatusOK {
				responseData, err := io.ReadAll(response.Body)
				if err == nil {
					var responseMessage Message
					if err := json.Unmarshal(responseData, &responseMessage); err == nil {

						// if json output is requested, output the JSON to stdout before erroring
						if jsonOutput {
							prettyJSON, err := prettyPrintJSON(responseData)
							if err == nil {
								fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
							} else {
								return errors.New(response.Status)
							}
						}
						return errors.New(responseMessage.Message)
					} else {
						return errors.New(response.Status)
					}
				} else {
					return errors.New(response.Status)
				}
			}

			// unmarshal into struct
			var result EnginesV4
			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			} else {
				if f.count {
					// simply return the count of engines
					fmt.Fprintln(cmd.OutOrStdout(), result.TotalCount)
				} else if f.outputType == "table" { // output a table with some default fields selected
					t := table.NewWriter()
					t.SetStyle(defaultStyle)

					t.AppendHeader(table.Row{"Name", "Host", "Build", "Platform", "Type", "Current Job ID", "Registration Properties", "Queues"})

					for _, element := range result.Items {
						t.AppendRow(table.Row{element.Name, element.Hostname, element.BuildNumber, element.Platform, element.Type, element.CurrentJobID, element.RegistrationProperties, element.AssignedQueues})
					}
					if f.noHeaders {
						t.ResetHeaders()
					}
					fmt.Fprintln(cmd.OutOrStdout(), t.Render())
					// output the raw json but formatted
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
					if f.noHeaders {
						t.ResetHeaders()
					}
					fmt.Fprintln(cmd.OutOrStdout(), t.Render())

				} else {
					return errors.New("invalid output format specified")
				}

			}

			return nil
		} else if f.apiVersion == "v3" {
			// call the status endpoint to see if it is finished
			request, err := buildFmeFlowRequest("/fmerest/v3/transformations/engines", "GET", nil)
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

			var result EnginesV3
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			} else {
				if f.count {
					// simply return the count of engines
					fmt.Fprintln(cmd.OutOrStdout(), result.TotalCount)
				} else if f.outputType == "table" { // output a table with some default fields selected
					t := table.NewWriter()
					t.SetStyle(defaultStyle)

					t.AppendHeader(table.Row{"Name", "Host", "Build", "Platform", "Type", "Current Job ID", "Registration Properties", "Queues"})

					for _, element := range result.Items {
						t.AppendRow(table.Row{element.InstanceName, element.HostName, element.BuildNumber, element.Platform, element.Type, element.CurrentJobID, element.RegistrationProperties, element.AssignedQueues})
					}
					if f.noHeaders {
						t.ResetHeaders()
					}
					fmt.Fprintln(cmd.OutOrStdout(), t.Render())
					// output the raw json but formatted
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
		return nil
	}
}
