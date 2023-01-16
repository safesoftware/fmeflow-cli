package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type projectsFlags struct {
	name       string
	owner      string
	outputType string
	noHeaders  bool
}

type ProjectsResource struct {
	Id int `json:"id"`
}

type FMEServerProjects struct {
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
	TotalCount int       `json:"totalCount"`
	Items      []Project `json:"items"`
}

type Project struct {
	Owner              string                   `json:"owner"`
	UID                string                   `json:"uid"`
	LastSaveDate       time.Time                `json:"lastSaveDate"`
	HasIcon            bool                     `json:"hasIcon"`
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	Sharable           bool                     `json:"sharable"`
	Readme             string                   `json:"readme"`
	UserName           string                   `json:"userName"`
	Version            string                   `json:"version"`
	FmeHubPublisherUID string                   `json:"fmeHubPublisherUid"`
	Accounts           []ProjectItem            `json:"accounts"`
	AppSuites          []ProjectItem            `json:"appSuites"`
	Apps               []ProjectItem            `json:"apps"`
	AutomationApps     []MutableProjectItemName `json:"automationApps"`
	Automations        []MutableProjectItemName `json:"automations"`
	CleanupTasks       []struct {
		Category string `json:"category"`
		Name     string `json:"name"`
	}
	Connections         []ProjectItem      `json:"connections"`
	CustomFormats       []RepositoryItem   `json:"customFormats"`
	CustomTransformers  []RepositoryItem   `json:"customTransformers"`
	Projects            []ProjectItem      `json:"projects"`
	Publications        []ProjectItem      `json:"publications"`
	Repositories        []ProjectItem      `json:"repositories"`
	ResourceConnections []ProjectItem      `json:"resourceConnections"`
	ResourcePaths       []ResourcePathItem `json:"resourcePaths"`
	Roles               []ProjectItem      `json:"roles"`
	Schedules           []struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	Streams       []MutableProjectItemName `json:"streams"`
	Subscriptions []ProjectItem            `json:"subscriptions"`
	Templates     []RepositoryItem         `json:"templates"`
	Tokens        []struct {
		Name     string `json:"name"`
		UserName string `json:"userName"`
	}
	Topics     []ProjectItem    `json:"topics"`
	Workspaces []RepositoryItem `json:"workspaces"`
}

type ProjectItem struct {
	Name string `json:"name"`
}

type MutableProjectItemName struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type RepositoryItem struct {
	Name           string `json:"name"`
	RepositoryName string `json:"repositoryName"`
}

type ResourcePathItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func newProjectsCmd() *cobra.Command {
	f := projectsFlags{}
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "Lists all projects on the FME Server",
		Long:  "Lists all projects on the FME Server",
		Example: `
  # List all projects
  fmeserver projects

  # List all projects owned by the user admin
  fmeserver projects --owner admin`,
		Args: NoArgs,
		RunE: projectsRun(&f),
	}

	cmd.Flags().StringVarP(&f.owner, "file", "f", "", "Path to backup file to upload to restore. Can be a local file or the relative path inside the specified shared resource.")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")

	return cmd
}
func projectsRun(f *projectsFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		client := &http.Client{}
		url := "/fmerest/v3/projects/projects"
		if f.name != "" {
			// add the project name to the request if specified
			url = url + "/" + f.name
		}

		// set up the URL to query
		request, err := buildFmeServerRequest(url, "GET", nil)
		if err != nil {
			return err
		}

		q := request.URL.Query()

		if f.owner != "" {
			// add the owner as a query parameter if specified
			q.Add("owner", f.owner)
			request.URL.RawQuery = q.Encode()
		}

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusOK {
			return errors.New(response.Status)
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result FMEServerProjects
		if f.name == "" {
			// if no name specified, request will return the full struct
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			}
		} else {
			// else, we aree getting a single repository. We will just append this
			// to the Item list in the full struct for easier parsing
			var singleResult Project
			if err := json.Unmarshal(responseData, &singleResult); err != nil {
				return err
			}
			result.TotalCount = 1
			result.Items = append(result.Items, singleResult)
		}

		if f.outputType == "table" {

			t := table.NewWriter()
			t.SetStyle(defaultStyle)

			t.AppendHeader(table.Row{"Name", "Owner", "Description", "Last Saved"})

			for _, element := range result.Items {
				t.AppendRow(table.Row{element.Name, element.Owner, element.Description, element.LastSaveDate})
			}
			if f.noHeaders {
				t.ResetHeaders()
			}
			fmt.Fprintln(cmd.OutOrStdout(), t.Render())

		} else if f.outputType == "json" {
			prettyJSON, err := prettyPrintJSON(responseData)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
		} else if strings.HasPrefix(f.outputType, "custom-columns") {
			// parse the columns and json queries
			columnsString := ""
			if strings.HasPrefix(f.outputType, "custom-columns=") {
				columnsString = f.outputType[len("custom-columns="):]
			}
			if len(columnsString) == 0 {
				return errors.New("custom-columns format specified but no custom columns given")
			}

			// we have to marshal the Items array, then create an array of marshalled items
			// to pass to the creation of the table.
			marshalledItems := [][]byte{}
			for _, element := range result.Items {
				mJson, err := json.Marshal(element)
				if err != nil {
					return err
				}
				marshalledItems = append(marshalledItems, mJson)
			}

			columnsInput := strings.Split(columnsString, ",")
			t, err := createTableFromCustomColumns(marshalledItems, columnsInput)
			if err != nil {
				return err
			}
			if noHeaders {
				t.ResetHeaders()
			}
			fmt.Fprintln(cmd.OutOrStdout(), t.Render())

		} else {
			return errors.New("invalid output format specified")
		}

		return nil
	}
}
