package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

type NewDeplymentParameter struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type deploymentParameterCreateFlags struct {
	dpType string
	value  string
	name   string
}

func newDeploymentParameterCreateCmd() *cobra.Command {
	f := deploymentParameterCreateFlags{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a deployment parameter",
		Long:  `Create a deployment parameter.`,
		Example: `
	Examples:
	# Create a deployment parameter with the name "myParam" and the value "myValue"
	fmeserver deploymentparameters create --name myParam --value myValue
`,

		Args: NoArgs,
		RunE: deploymentParametersCreateRun(&f),
	}

	cmd.Flags().StringVar(&f.dpType, "type", "text", "Type of parameter")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the deployment parameter to create.")
	cmd.Flags().StringVar(&f.value, "value", "", "The value to set the deployment parameter to.")
	// currently type can only be set to "text". But in the future maybe there will be more options?
	cmd.Flags().MarkHidden("type")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("value")
	return cmd
}

func deploymentParametersCreateRun(f *deploymentParameterCreateFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		var newDepParam NewDeplymentParameter
		newDepParam.Name = f.name
		newDepParam.Value = f.value
		newDepParam.Type = f.dpType
		jsonData, err := json.Marshal(newDepParam)
		if err != nil {
			return err
		}

		request, err := buildFmeServerRequest("/fmeapiv4/deploymentparameters", "POST", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", "application/json")

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusCreated {
			// attempt to parse the body into JSON as there could be a valuable message in there
			// if fail, just output the status code
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

		if !jsonOutput {
			fmt.Fprintln(cmd.OutOrStdout(), "Deployment Parameter successfully created.")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "{}")
		}

		return nil
	}
}
