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
	"github.com/spf13/viper"
)

type projectsFlags struct {
	id         string
	name       string
	owner      string
	outputType string
	noHeaders  bool
	apiVersion apiVersionFlag
}

type ProjectV4 struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	HubUID           string    `json:"hubUid"`
	HubPublisherUID  string    `json:"hubPublisherUid"`
	Description      string    `json:"description"`
	Readme           string    `json:"readme"`
	Version          string    `json:"version"`
	LastUpdated      time.Time `json:"lastUpdated"`
	Owner            string    `json:"owner"`
	OwnerID          string    `json:"ownerID"`
	Shareable        bool      `json:"shareable"`
	LastUpdateUser   string    `json:"lastUpdateUser"`
	LastUpdateUserID string    `json:"lastUpdateUserID"`
	HasIcon          bool      `json:"hasIcon"`
}

type FMEFlowProjectsV4 struct {
	Offset     int         `json:"offset"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"totalCount"`
	Items      []ProjectV4 `json:"items"`
}

type ProjectsResource struct {
	Id int `json:"id"`
}

type FMEFlowProjectsV3 struct {
	Offset     int         `json:"offset"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"totalCount"`
	Items      []ProjectV3 `json:"items"`
}

type ProjectV3 struct {
	Owner              string                     `json:"owner"`
	UID                string                     `json:"uid"`
	LastSaveDate       time.Time                  `json:"lastSaveDate"`
	HasIcon            bool                       `json:"hasIcon"`
	Name               string                     `json:"name"`
	Description        string                     `json:"description"`
	Sharable           bool                       `json:"sharable"`
	Readme             string                     `json:"readme"`
	UserName           string                     `json:"userName"`
	Version            string                     `json:"version"`
	FmeHubPublisherUID string                     `json:"fmeHubPublisherUid"`
	Accounts           []ProjectItemV3            `json:"accounts"`
	AppSuites          []ProjectItemV3            `json:"appSuites"`
	Apps               []ProjectItemV3            `json:"apps"`
	AutomationApps     []MutableProjectItemNameV3 `json:"automationApps"`
	Automations        []MutableProjectItemNameV3 `json:"automations"`
	CleanupTasks       []struct {
		Category string `json:"category"`
		Name     string `json:"name"`
	}
	Connections         []ProjectItemV3      `json:"connections"`
	CustomFormats       []RepositoryItemV3   `json:"customFormats"`
	CustomTransformers  []RepositoryItemV3   `json:"customTransformers"`
	Projects            []ProjectItemV3      `json:"projects"`
	Publications        []ProjectItemV3      `json:"publications"`
	Repositories        []ProjectItemV3      `json:"repositories"`
	ResourceConnections []ProjectItemV3      `json:"resourceConnections"`
	ResourcePaths       []ResourcePathItemV3 `json:"resourcePaths"`
	Roles               []ProjectItemV3      `json:"roles"`
	Schedules           []struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	Streams       []MutableProjectItemNameV3 `json:"streams"`
	Subscriptions []ProjectItemV3            `json:"subscriptions"`
	Templates     []RepositoryItemV3         `json:"templates"`
	Tokens        []struct {
		Name     string `json:"name"`
		UserName string `json:"userName"`
	}
	Topics     []ProjectItemV3    `json:"topics"`
	Workspaces []RepositoryItemV3 `json:"workspaces"`
}

type ProjectItemV3 struct {
	Name string `json:"name"`
}

type MutableProjectItemNameV3 struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type RepositoryItemV3 struct {
	Name           string `json:"name"`
	RepositoryName string `json:"repositoryName"`
}

type ResourcePathItemV3 struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

var projectsV4BuildThreshold = 23283

func newProjectsCmd() *cobra.Command {
	f := projectsFlags{}
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "Lists projects on the FME Server",
		Long:  "Lists projects on the FME Server. Pass in a name to retrieve information on a single project.",
		Example: `
  # List all projects
  fmeflow projects

  # List all projects owned by the user admin
  fmeflow projects --owner admin
  
  # Get a single project by name
  fmeflow projects --name "My Project"
  
  # Get a single project by id
  fmeflow projects --id a64297e7-a119-4e10-ac37-5d0bba12194b
  
  # Get a single project by name and output as JSON
  fmeflow projects --name "My Project" --output json
  
  # Get all projects and output as custom columns
  fmeflow projects --output=custom-columns=ID:.id,NAME:.name`,
		Args: NoArgs,
		RunE: projectsRun(&f),
	}

	cmd.Flags().StringVar(&f.owner, "owner", "", "If specified, only projects owned by the specified user will be returned.")
	cmd.Flags().StringVar(&f.name, "name", "", "Return a single project with the given name.")
	cmd.Flags().StringVar(&f.id, "id", "", "Return a single project with the given id. (v4 only)")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.AddCommand(newProjectDownloadCmd())
	cmd.AddCommand(newProjectUploadCmd())
	cmd.AddCommand(newProjectItemsCmd())
	cmd.AddCommand(newProjectDeleteCmd())
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)

	return cmd
}
func projectsRun(f *projectsFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		if f.apiVersion == "" {
			if viper.GetInt("build") < projectsV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		client := &http.Client{}

		if f.apiVersion == "v4" {

			if f.name != "" {
				id, err := getProjectId(f.name)
				if err != nil {
					return err
				}
				f.id = id
			}

			url := "/fmeapiv4/projects"
			if f.id != "" {
				url = url + "/" + f.id
			}

			request, err := buildFmeFlowRequest(url, "GET", nil)
			if err != nil {
				return err
			}

			if f.owner != "" {
				q := request.URL.Query()
				q.Add("filterString", f.owner)
				q.Add("filterProperties", "owner")
				request.URL.RawQuery = q.Encode()
			}

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != http.StatusOK {
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

			// marshal into struct
			var result FMEFlowProjectsV4

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			if f.id == "" {
				// if no id specified, request will return the full struct
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				}
			} else {
				var projectStruct ProjectV4

				if err := json.Unmarshal(responseData, &projectStruct); err != nil {
					return err
				}
				result.TotalCount = 1
				result.Items = append(result.Items, projectStruct)
			}

			if f.outputType == "table" {

				t := table.NewWriter()
				t.SetStyle(defaultStyle)

				t.AppendHeader(table.Row{"ID", "Name", "Owner", "Description", "Last Updated"})

				for _, element := range result.Items {
					t.AppendRow(table.Row{element.ID, element.Name, element.Owner, element.Description, element.LastUpdated})
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
				if f.noHeaders {
					t.ResetHeaders()
				}
				fmt.Fprintln(cmd.OutOrStdout(), t.Render())

			} else {
				return errors.New("invalid output format specified")
			}

			return nil

		} else if f.apiVersion == "v3" {

			url := "/fmerest/v3/projects/projects"
			if f.name != "" {
				// add the project name to the request if specified
				url = url + "/" + f.name
			}

			// set up the URL to query
			request, err := buildFmeFlowRequest(url, "GET", nil)
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
				if response.StatusCode == http.StatusNotFound {
					return fmt.Errorf("%w: check that the specified project exists", errors.New(response.Status))
				} else {
					return errors.New(response.Status)
				}
			}

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result FMEFlowProjectsV3
			if f.name == "" {
				// if no name specified, request will return the full struct
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				}
			} else {
				// else, we are getting a single project. We will just append this
				// to the Item list in the full struct for easier parsing
				var singleResult ProjectV3
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
				if f.noHeaders {
					t.ResetHeaders()
				}
				fmt.Fprintln(cmd.OutOrStdout(), t.Render())

			} else {
				return errors.New("invalid output format specified")
			}

			return nil
		}
		return nil
	}
}
