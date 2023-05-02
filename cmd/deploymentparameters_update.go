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

type UpdateDeplymentParameter struct {
	Value string `json:"value"`
}

type deploymentParameterUpdateFlags struct {
	value string
	name  string
}

func newDeploymentParameterUpdateCmd() *cobra.Command {
	f := deploymentParameterUpdateFlags{}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a deployment parameter",
		Long:  `Update a deployment parameter.`,
		Example: `
	Examples:
	# Update a deployment parameter with the name "myParam" and the value "myValue"
	fmeflow deploymentparameters update --name myParam --value myValue
`,

		Args: NoArgs,
		RunE: deploymentParametersUpdateRun(&f),
	}

	cmd.Flags().StringVar(&f.name, "name", "", "Name of the deployment parameter to update.")
	cmd.Flags().StringVar(&f.value, "value", "", "The value to set the deployment parameter to.")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("value")
	return cmd
}

func deploymentParametersUpdateRun(f *deploymentParameterUpdateFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		var newDepParam NewDeplymentParameter
		newDepParam.Value = f.value
		jsonData, err := json.Marshal(newDepParam)
		if err != nil {
			return err
		}

		request, err := buildFmeFlowRequest("/fmeapiv4/deploymentparameters/"+f.name, "PUT", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", "application/json")

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusNoContent {
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
			fmt.Fprintln(cmd.OutOrStdout(), "Deployment Parameter successfully updated.")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "{}")
		}

		return nil
	}
}
