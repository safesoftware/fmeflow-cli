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

type UpdateDeploymentParameter struct {
	Type           string          `json:"type,omitempty"`
	Value          string          `json:"value"`
	ChoiceSettings *ChoiceSettings `json:"choiceSettings,omitempty"`
}

type ChoiceSettings struct {
	ChoiceSet        string   `json:"choiceSet,omitempty"`
	Services         []string `json:"services,omitempty"`
	ExcludedServices []string `json:"excludedServices,omitempty"`
	Family           string   `json:"family,omitempty"`
}

type deploymentParameterUpdateFlags struct {
	dpType           deploymentParameterTypeFlag
	value            string
	name             string
	dbType           string
	includedServices []string
	excludedServices []string
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

	cmd.Flags().Var(&f.dpType, "type", "Update the type of the parameter. Must be one of text, database, or web. Default is text.")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the deployment parameter to update.")
	cmd.Flags().StringVar(&f.value, "value", "", "The value to set the deployment parameter to.")
	cmd.Flags().StringArrayVar(&f.includedServices, "included-service", []string{}, "Service to include in the deployment parameter. Can be passed in multiple times if there are multiple Web services to include.")
	cmd.Flags().StringArrayVar(&f.excludedServices, "excluded-service", []string{}, "Service to exclude in the deployment parameter. Can be passed in multiple times if there are multiple Web services to exclude.")
	cmd.Flags().StringVar(&f.dbType, "database-type", "", "The type of the database to use for the database deployment parameter. (Optional)")
	cmd.RegisterFlagCompletionFunc("type", deploymentParameterTypeFlagCompletion)
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("value")
	return cmd
}

func deploymentParametersUpdateRun(f *deploymentParameterUpdateFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// if type isnot web and includeServices or excludedServices is set, then error
		if f.dpType != deploymentParameterTypeFlagWeb && (len(f.includedServices) > 0 || len(f.excludedServices) > 0) {
			return errors.New("cannot include or exclude services for a non-web connection deployment parameter")
		}

		// if the type is not database and dbType is set, then error
		if f.dpType != deploymentParameterTypeFlagDatabase && f.dbType != "" {
			return errors.New("cannot set a database family for a non-database deployment parameter")
		}

		// set up http
		client := &http.Client{}

		// check if deployment parameter exists first and error if it does not
		var currParam DeploymentParameter
		request, err := buildFmeFlowRequest("/fmeapiv4/deploymentparameters/"+f.name, "GET", nil)
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
		} else {
			// get the current parameter
			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(responseData, &currParam); err != nil {
				return err
			}
		}

		var newDepParam UpdateDeploymentParameter
		newDepParam.ChoiceSettings = new(ChoiceSettings)
		newDepParam.Value = f.value

		// if type is passed in, set the type and check for other related settings
		// default to the current value if not passed in
		if cmd.Flags().Changed("type") {
			// set type specific settings
			if f.dpType == deploymentParameterTypeFlagDatabase {
				newDepParam.Type = "dropdown"
				newDepParam.ChoiceSettings.ChoiceSet = "dbConnections"
				if cmd.Flags().Changed("database-type") {
					newDepParam.ChoiceSettings.Family = f.dbType
				} else {
					newDepParam.ChoiceSettings.Family = currParam.ChoiceSettings.Family
				}
			} else if f.dpType == deploymentParameterTypeFlagWeb {
				newDepParam.Type = "dropdown"
				newDepParam.ChoiceSettings.ChoiceSet = "webConnections"
				if cmd.Flags().Changed(("included-service")) {
					newDepParam.ChoiceSettings.Services = f.includedServices
				} else {
					newDepParam.ChoiceSettings.Services = currParam.ChoiceSettings.Services
				}
				if cmd.Flags().Changed(("excluded-service")) {
					newDepParam.ChoiceSettings.ExcludedServices = f.excludedServices
				} else {
					newDepParam.ChoiceSettings.ExcludedServices = currParam.ChoiceSettings.ExcludedServices
				}
			} else {
				// text type just sets type and nothing else
				newDepParam.Type = f.dpType.String()
				newDepParam.ChoiceSettings = nil
			}
		} else {
			// if type isn't passed in, then use all the same settings as the current parameter
			newDepParam.Type = currParam.Type
			newDepParam.ChoiceSettings = (*ChoiceSettings)(&currParam.ChoiceSettings)
		}

		jsonData, err := json.Marshal(newDepParam)
		if err != nil {
			return err
		}

		request, err = buildFmeFlowRequest("/fmeapiv4/deploymentparameters/"+f.name, "PUT", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", "application/json")

		response, err = client.Do(&request)
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
