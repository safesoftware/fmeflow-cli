package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type DeploymentParameters struct {
	Items      []DeploymentParameter `json:"items"`
	Limit      int                   `json:"limit"`
	Offset     int                   `json:"offset"`
	TotalCount int                   `json:"totalCount"`
}

type DeploymentParameter struct {
	Name    string    `json:"name"`
	Owner   string    `json:"owner"`
	Type    string    `json:"type"`
	Updated time.Time `json:"updated"`
	Value   string    `json:"value"`
}

type deploymentparametersFlags struct {
	name       string
	outputType string
	noHeaders  bool
}

var deploymentParametersBuildThreshold = 23170

func newDeploymentParametersCmd() *cobra.Command {
	f := deploymentparametersFlags{}
	cmd := &cobra.Command{
		Use:   "deploymentparameters",
		Short: "List Deployment Parameters",
		Long:  `Lists Deployment Parameters on the given FME Server.`,
		Example: `
	Examples:
	# List all deployment parameters
	fmeflow deploymentparameters
	
	# List a single deployment parameter
	fmeflow deploymentparameters --name testParameter
	
	# Output all deploymentparameters in json format
	fmeflow deploymentparameters --json`,
		Args: NoArgs,
		RunE: deploymentParametersRun(&f),
	}

	cmd.Flags().StringVar(&f.name, "name", "", "If specified, only the repository with that name will be returned")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	cmd.AddCommand(newDeploymentParameterCreateCmd())
	cmd.AddCommand(newDeploymentParameterDeleteCmd())
	cmd.AddCommand(newDeploymentParameterUpdateCmd())

	return cmd
}

func deploymentParametersRun(f *deploymentparametersFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		// set up http
		client := &http.Client{}

		// set up the URL to query
		url := "/fmeapiv4/deploymentparameters"
		if f.name != "" {
			// add the repository name to the request if specified
			url = url + "/" + f.name
		}
		request, err := buildFmeFlowRequest(url, "GET", nil)
		if err != nil {
			return err
		}

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusOK {
			// attempt to parse the body into JSON as there could be a valuable message in there
			// if fail, just output the status code
			responseData, err := io.ReadAll(response.Body)
			if err == nil {

				var responseMessage Message
				if err := json.Unmarshal(responseData, &responseMessage); err == nil {

					// if json output is requested, output the JSON to stdout before erroring
					if f.outputType == "json" {
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

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result DeploymentParameters

		if f.name == "" {
			// if no name specified, request will return the full struct
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			}
		} else {
			// else, we aree getting a single repository. We will just append this
			// to the Item list in the full struct for easier parsing
			var singleResult DeploymentParameter
			if err := json.Unmarshal(responseData, &singleResult); err != nil {
				return err
			}
			result.TotalCount = 1
			result.Items = append(result.Items, singleResult)
		}

		if f.outputType == "table" {

			t := table.NewWriter()
			t.SetStyle(defaultStyle)

			t.AppendHeader(table.Row{"Name", "Owner", "Type", "Value", "Last Updated"})

			for _, element := range result.Items {
				t.AppendRow(table.Row{element.Name, element.Owner, element.Type, element.Value, element.Updated})
			}
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

		return nil
	}
}
