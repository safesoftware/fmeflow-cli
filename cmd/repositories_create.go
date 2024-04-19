package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type NewRepository struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

type repositoryCreateFlags struct {
	description string
	name        string
	apiVersion  apiVersionFlag
}

func newRepositoryCreateCmd() *cobra.Command {
	f := repositoryCreateFlags{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new repository.",
		Long:  `Create a new repository.`,
		Example: `
  # Create a repository with the name "myRepository" and no description
  fmeflow repositories create --name myRepository
	
  # Output just the name of all the repositories
  fmeflow repositories create --name myRepository --description "This is my new repository"
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// get build to decide if we should use v3 or v4
			// FME Server 2023.0 and later can use v4. Otherwise fall back to v3
			if f.apiVersion == "" {
				fmeflowBuild := viper.GetInt("build")
				if fmeflowBuild < repositoriesV4BuildThreshold {
					f.apiVersion = apiVersionFlagV3
				} else {
					f.apiVersion = apiVersionFlagV4
				}
			}

			return nil
		},
		Args: NoArgs,
		RunE: repositoriesCreateRun(&f),
	}

	cmd.Flags().StringVar(&f.description, "description", "", "Description of the new repository.")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the repository to create.")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	cmd.MarkFlagRequired("name")
	return cmd
}

func repositoriesCreateRun(f *repositoryCreateFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		if f.apiVersion == "v4" {

			var newRepo NewRepository
			newRepo.Name = f.name
			newRepo.Description = f.description
			jsonData, err := json.Marshal(newRepo)
			if err != nil {
				return err
			}

			request, err := buildFmeFlowRequest("/fmeapiv4/repositories", "POST", bytes.NewBuffer(jsonData))
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
				} else if response.StatusCode == http.StatusNotFound {
					return fmt.Errorf("%w: The repository does not exist", errors.New(response.Status))
				} else {
					return errors.New(response.Status)
				}
			}

			if !jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), "Repository successfully created.")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "{}")
			}
		} else if f.apiVersion == "v3" {

			// add mandatory values
			data := url.Values{
				"name": {f.name},
			}

			// add optional values
			if f.description != "" {
				data.Add("description", f.description)
			}

			request, err := buildFmeFlowRequest("/fmerest/v3/repositories", "POST", strings.NewReader(data.Encode()))
			if err != nil {
				return err
			}

			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode == http.StatusUnprocessableEntity {
				return fmt.Errorf("%w: Some or all of the input parameters are invalid", errors.New(response.Status))
			} else if response.StatusCode == http.StatusConflict {
				return fmt.Errorf("%w: The repository already exists", errors.New(response.Status))
			} else if response.StatusCode != http.StatusCreated {
				return errors.New(response.Status)
			}

			if !jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), "Repository successfully created.")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "{}")
			}

		}
		return nil
	}
}
