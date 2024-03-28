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

type projectDeleteFlags struct {
	id           string
	name         string
	all          bool
	dependencies bool
	noprompt     bool
}

func newProjectDeleteCmd() *cobra.Command {
	f := projectDeleteFlags{}
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes an FME Flow Project",
		Long:  `Deletes an FME Flow Project from the FME Server. Can optionally also delete the project contents and its dependencies.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if f.id == "" && f.name == "" {
				return errors.New("either id or name must be specified")
			}

			if f.dependencies && !f.all {
				return errors.New("dependencies can only be specified if all is also specified")
			}

			return nil
		},
		Example: `
  # Delete a project by id
  fmeflow projects delete --id 123

  # Delete a project by name
  fmeflow projects delete --name "My Project"
  
  # Delete a project by name and all its contents
  fmeflow projects delete --name "My Project" --all
  
  # Delete a project by name and all its contents and dependencies
  fmeflow projects delete --name "My Project" --all --dependencies`,
		Args: NoArgs,
		RunE: projectDeleteRun(&f),
	}

	cmd.Flags().StringVar(&f.id, "id", "", "The id of the project to delete. Either id or name must be specified")
	cmd.Flags().StringVar(&f.name, "name", "", "The name of the project to delete. Either id or name must be specified")
	cmd.Flags().BoolVar(&f.all, "all", false, "Delete the project and its contents")
	cmd.Flags().BoolVar(&f.dependencies, "dependencies", false, "Delete the project and its contents and dependencies. Can only be specified if all is also specified")
	cmd.Flags().BoolVarP(&f.noprompt, "no-prompt", "y", false, "Do not prompt for confirmation.")
	cmd.MarkFlagsMutuallyExclusive("id", "name")

	return cmd
}

func projectDeleteRun(f *projectDeleteFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := &http.Client{}

		// get project id if name was passed in
		if f.id == "" {
			projectID, err := getProjectId(f.name)
			if err != nil {
				return err
			}
			f.id = projectID
		} else if !f.noprompt {
			// check if the project exists if we are going to prompt to confirm deletion
			request, err := buildFmeFlowRequest("/fmeapiv4/projects/"+f.id, "GET", nil)
			if err != nil {
				return err
			}
			response, err := client.Do(&request)
			if err != nil {
				return err
			}
			if response.StatusCode != http.StatusOK {
				return errors.New(response.Status + ": check that the project id is correct")
			}
		}

		// the project exists. Confirm deletion.
		if !f.noprompt {
			// prompt to confirm deletion
			confirm := false
			promptUser := &survey.Confirm{
				Message: "Are you sure you want to delete the project " + f.name + "?",
			}
			survey.AskOne(promptUser, &confirm)
			if !confirm {
				return nil
			}
		}

		url := "/fmeapiv4/projects/" + f.id
		if f.all {
			url += "/delete-all"
		}

		request, err := buildFmeFlowRequest(url, "DELETE", nil)
		if err != nil {
			return err
		}

		if f.dependencies {
			q := request.URL.Query()
			q.Add("deleteDependencies", "true")
			request.URL.RawQuery = q.Encode()
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
			fmt.Fprintln(cmd.OutOrStdout(), "Project successfully deleted.")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "{}")
		}

		return nil
	}
}
