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

type ConnectionDeleteFlags struct {
	name     string
	noprompt bool
}

func newConnectionDeleteCmd() *cobra.Command {
	f := ConnectionDeleteFlags{}
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a connection",
		Long:  `Delete a connection.`,
		Example: `
	Examples:
	# Delete a connection with the name "myConnection"
	fmeflow connections delete --name myConnection
`,

		Args: NoArgs,
		RunE: connectionDeleteRun(&f),
	}

	cmd.Flags().StringVar(&f.name, "name", "", "Name of the connection to delete.")
	cmd.Flags().BoolVarP(&f.noprompt, "no-prompt", "y", false, "Description of the new repository.")

	cmd.MarkFlagRequired("name")
	return cmd
}

func connectionDeleteRun(f *ConnectionDeleteFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		// check if deployment parameter exists first and error if it does not
		request, err := buildFmeFlowRequest("/fmeapiv4/connections/"+f.name, "GET", nil)
		if err != nil {
			return err
		}
		// send the request
		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusOK {
			// if we didn't get a 200 OK, then the deployment parameter does not exist
			// get the JSON response and throw a new error using the message
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
			}
		}

		// the parameter exists. Confirm deletion.
		if !f.noprompt {
			// prompt to confirm deletion
			confirm := false
			promptUser := &survey.Confirm{
				Message: "Are you sure you want to delete the deployment parameter " + f.name + "?",
			}
			survey.AskOne(promptUser, &confirm)
			if !confirm {
				return nil
			}
		}

		// get the current values of the connection we are going to update
		url := "/fmeapiv4/connections/" + f.name
		request, err = buildFmeFlowRequest(url, "DELETE", nil)
		if err != nil {
			return err
		}

		response, err = client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusNoContent {
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
			fmt.Fprintln(cmd.OutOrStdout(), "Connection successfully deleted.")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "{}")
		}

		return nil
	}
}
