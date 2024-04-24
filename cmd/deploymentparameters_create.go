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

type NewDeploymentParameter struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Value          string `json:"value"`
	ChoiceSettings struct {
		ChoiceSet        string   `json:"choiceSet"`
		Services         []string `json:"services,omitempty"`
		ExcludedServices []string `json:"excludedServices,omitempty"`
		Family           string   `json:"family,omitempty"`
	} `json:"choiceSettings,omitempty"`
}

type deploymentParameterCreateFlags struct {
	dpType           deploymentParameterTypeFlag
	value            string
	name             string
	dbType           string
	includedServices []string
	excludedServices []string
}

func newDeploymentParameterCreateCmd() *cobra.Command {
	f := deploymentParameterCreateFlags{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a deployment parameter",
		Long:  `Create a deployment parameter.`,
		Example: `
  # Create a Web deployment parameter including the slack service and specifying a slack connection
  fmeflow deploymentparameters create --type web --name slack_connection --included-service Slack --value slack_connection_value

  # Create a Database deployment parameter for PostgreSQL specifying a pgsql connection
  fmeflow deploymentparameters create --type database --name pgsql_param --database-type PostgreSQL --value pgsql_connection_value

  # Create a Text deployment parameter
  fmeflow deploymentparameters create --name text_connection --value text_connection_value --type text
`,

		Args: NoArgs,
		RunE: deploymentParametersCreateRun(&f),
	}

	cmd.Flags().Var(&f.dpType, "type", "Type of parameter to create. Must be one of text, database, or web. Default is text.")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the deployment parameter to create.")
	cmd.Flags().StringVar(&f.value, "value", "", "The value to set the deployment parameter to. (Optional)")
	cmd.Flags().StringArrayVar(&f.includedServices, "included-service", []string{}, "Service to include in the deployment parameter. Can be passed in multiple times if there are multiple Web services to include.")
	cmd.Flags().StringArrayVar(&f.excludedServices, "excluded-service", []string{}, "Service to exclude in the deployment parameter. Can be passed in multiple times if there are multiple Web services to exclude.")
	cmd.Flags().StringVar(&f.dbType, "database-type", "", "The type of the database to use for the database deployment parameter. (Optional)")
	cmd.RegisterFlagCompletionFunc("type", deploymentParameterTypeFlagCompletion)
	cmd.MarkFlagRequired("name")
	return cmd
}

func deploymentParametersCreateRun(f *deploymentParameterCreateFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		if f.dpType == "" {
			f.dpType = deploymentParameterTypeFlagText
		}

		// if type is not web and includeServices or excludedServices is set, then error
		if f.dpType != deploymentParameterTypeFlagWeb && (len(f.includedServices) > 0 || len(f.excludedServices) > 0) {
			return errors.New("cannot include or exclude services for a non-web connection deployment parameter")
		}

		// if the type is not database and dbType is set, then error
		if f.dpType != deploymentParameterTypeFlagDatabase && f.dbType != "" {
			return errors.New("cannot set a database family for a non-database deployment parameter")
		}

		// set up http
		client := &http.Client{}

		var newDepParam NewDeploymentParameter
		newDepParam.Name = f.name
		newDepParam.Value = f.value

		// set type specific settings
		if f.dpType == deploymentParameterTypeFlagDatabase {
			newDepParam.Type = "dropdown"
			newDepParam.ChoiceSettings.ChoiceSet = "dbConnections"
			if f.dbType != "" {
				newDepParam.ChoiceSettings.Family = f.dbType
			}
		} else if f.dpType == deploymentParameterTypeFlagWeb {
			newDepParam.Type = "dropdown"
			newDepParam.ChoiceSettings.ChoiceSet = "webConnections"
			newDepParam.ChoiceSettings.Services = f.includedServices
			newDepParam.ChoiceSettings.ExcludedServices = f.excludedServices
		} else {
			// text type just sets type and nothing else
			newDepParam.Type = f.dpType.String()
		}

		jsonData, err := json.Marshal(newDepParam)
		if err != nil {
			return err
		}

		request, err := buildFmeFlowRequest("/fmeapiv4/deploymentparameters", "POST", bytes.NewBuffer(jsonData))
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
