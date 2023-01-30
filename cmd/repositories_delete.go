package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type repositoryDeleteFlags struct {
	name       string
	noprompt   bool
	apiVersion apiVersionFlag
}

func newRepositoryDeleteCmd() *cobra.Command {
	f := repositoryDeleteFlags{}
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a repository",
		Long:  `Delete a repository.`,
		Example: `
	Examples:
	# Delete a repository with the name "myRepository"
	fmeserver repositories delete --name myRepository
	
	# Delete a repository with the name "myRepository" and no confirmation
	fmeserver repositories delete --name myRepository --no-prompt
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// get build to decide if we should use v3 or v4
			// FME Server 2023.0 and later can use v4. Otherwise fall back to v3
			if f.apiVersion == "" {
				fmeserverBuild := viper.GetInt("build")
				if fmeserverBuild < repositoriesV4BuildThreshold {
					f.apiVersion = apiVersionFlagV3
				} else {
					f.apiVersion = apiVersionFlagV4
				}
			}

			return nil
		},
		Args: NoArgs,
		RunE: repositoriesDeleteRun(&f),
	}

	cmd.Flags().BoolVarP(&f.noprompt, "no-prompt", "y", false, "Description of the new repository.")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the repository to create.")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	cmd.MarkFlagRequired("name")
	return cmd
}

func repositoriesDeleteRun(f *repositoryDeleteFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		if !f.noprompt {
			// prompt for a user and password
			confirm := false
			promptUser := &survey.Confirm{
				Message: "Are you sure you want to delete the repository " + f.name + "?",
			}
			survey.AskOne(promptUser, &confirm)
			if !confirm {
				return nil
			}
		}
		url := ""

		if f.apiVersion == "v4" {
			url = "/fmeapiv4/repositories/" + f.name
		} else if f.apiVersion == "v3" {
			url = "/fmerest/v3/repositories/" + f.name
		}

		request, err := buildFmeServerRequest(url, "DELETE", nil)
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
			fmt.Fprintln(cmd.OutOrStdout(), "Repository successfully deleted.")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "{}")
		}
		return nil
	}
}
