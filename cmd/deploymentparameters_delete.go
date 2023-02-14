package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type deploymentParameterDeleteFlags struct {
	name     string
	noprompt bool
}

func newDeploymentParameterDeleteCmd() *cobra.Command {
	f := deploymentParameterDeleteFlags{}
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a deployment parameter",
		Long:  `Delete a deployment parameter.`,
		Example: `
	Examples:
	# Delete adeployment parameter with the name "myParam"
	fmeserver deploymentparameters delete --name myParam
	
	# Delete a repository with the name "myRepository" and no confirmation
	fmeserver deploymentparameters delete --name myParam --no-prompt
`,
		Args: NoArgs,
		RunE: deploymentParameterDeleteRun(&f),
	}

	cmd.Flags().BoolVarP(&f.noprompt, "no-prompt", "y", false, "Description of the new repository.")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the repository to create.")
	cmd.MarkFlagRequired("name")
	return cmd
}

func deploymentParameterDeleteRun(f *deploymentParameterDeleteFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		if !f.noprompt {
			// prompt for a user and password
			confirm := false
			promptUser := &survey.Confirm{
				Message: "Are you sure you want to delete the deployment parameter " + f.name + "?",
			}
			survey.AskOne(promptUser, &confirm)
			if !confirm {
				return nil
			}
		}

		request, err := buildFmeServerRequest("/fmeapiv4/deploymentparameters/"+f.name, "DELETE", nil)
		if err != nil {
			return err
		}

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
			fmt.Fprintln(cmd.OutOrStdout(), "Deployment Parameter successfully deleted.")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "{}")
		}
		return nil
	}
}
