package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type repositoryDeleteFlags struct {
	name     string
	noprompt bool
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
		Args: NoArgs,
		RunE: repositoriesDeleteRun(&f),
	}

	cmd.Flags().BoolVar(&f.noprompt, "no-prompt", false, "Description of the new repository.")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the repository to create.")
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

		request, err := buildFmeServerRequest("/fmerest/v3/repositories/"+f.name, "DELETE", nil)
		if err != nil {
			return err
		}

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("%w: The repository does not exist", errors.New(response.Status))
		} else if response.StatusCode != http.StatusNoContent {
			return errors.New(response.Status)
		}

		if !jsonOutput {
			fmt.Fprintln(cmd.OutOrStdout(), "Repository successfully deleted.")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "{}")
		}
		return nil
	}
}
