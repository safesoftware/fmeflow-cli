package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

type repositoryCreateFlags struct {
	description string
	name        string
	//outputType  string
	//noHeaders   bool
}

func newRepositoryCreateCmd() *cobra.Command {
	f := repositoryCreateFlags{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a repository",
		Long:  `Create a new repository.`,
		Example: `
	Examples:
	# Create a repository with the name "myRepository" and no description
	fmeserver repositories create --name myRepository
	
	# Output just the name of all the repositories
	fmeserver repositories create --name myRepository --description "This is my new repository"
`,
		Args: NoArgs,
		RunE: repositoriesCreateRun(&f),
	}

	cmd.Flags().StringVar(&f.description, "description", "", "Description of the new repository.")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the repository to create.")
	//cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	//cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.MarkFlagRequired("name")
	return cmd
}

func repositoriesCreateRun(f *repositoryCreateFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		// add mandatory values
		data := url.Values{
			"name": {f.name},
		}

		// add optional values
		if f.description != "" {
			data.Add("description", f.description)
		}

		request, err := buildFmeServerRequest("/fmerest/v3/repositories", "POST", strings.NewReader(data.Encode()))
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
		return nil
	}
}
