package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FMEServerRepositoriesV3 struct {
	Offset     int                     `json:"offset"`
	Limit      int                     `json:"limit"`
	TotalCount int                     `json:"totalCount"`
	Items      []FMEServerRepositoryV3 `json:"items"`
}

type FMEServerRepositoryV3 struct {
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sharable    bool   `json:"sharable"`
}

type FMEServerRepositoriesV4 struct {
	Items      []FMEServerRepositoryV4 `json:"items"`
	Limit      int                     `json:"limit"`
	Offset     int                     `json:"offset"`
	TotalCount int                     `json:"totalCount"`
}

type FMEServerRepositoryV4 struct {
	CustomFormatCount      int    `json:"customFormatCount"`
	CustomTransformerCount int    `json:"customTransformerCount"`
	Description            string `json:"description"`
	FileCount              int    `json:"fileCount"`
	Name                   string `json:"name"`
	Owner                  string `json:"owner"`
	OwnerID                string `json:"ownerID"`
	Sharable               bool   `json:"sharable"`
	TemplateCount          int    `json:"templateCount"`
	TotalFileSize          int    `json:"totalFileSize"`
	WorkspaceCount         int    `json:"workspaceCount"`
}

type repositoryFlags struct {
	owner        string
	name         string
	filterString string
	outputType   string
	noHeaders    bool
	apiVersion   apiVersionFlag
}

var repositoriesV4BuildThreshold = 22337

func newRepositoryCmd() *cobra.Command {
	f := repositoryFlags{}
	cmd := &cobra.Command{
		Use:   "repositories",
		Short: "List repositories",
		Long:  `Lists repositories on the given FME Server.`,
		Example: `
	Examples:
	# List all repositories
	fmeserver repositories
	
	# List all repositories owned by the admin user
	fmeserver repositories --owner admin
	
	# List a single repository with the name "Samples"
	fmeserver repositories --name Samples
	
	# Output just the name of all the repositories
	fmeserver repositories --output=custom-columns=NAME:.name --no-headers
	
	# Output all repositories in json format
	fmeserver repositories --json`,
		Args: NoArgs,
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
			if f.apiVersion == apiVersionFlagV3 {
				if f.filterString != "" {
					return errors.New("cannot set the filter-string flag when using the V3 API")
				}
			}

			return nil
		},
		RunE: repositoriesRun(&f),
	}

	cmd.Flags().StringVar(&f.owner, "owner", "", "If specified, only repositories owned by the specified user uuid will be returned. With the V3 API, set this to the user name.")
	cmd.Flags().StringVar(&f.name, "name", "", "If specified, only the repository with that name will be returned")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().StringVar(&f.filterString, "filter-string", "", "Specify the output type. Should be one of table, json, or custom-columns. Only usable with V4 API.")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	cmd.AddCommand(newRepositoryCreateCmd())
	cmd.AddCommand(newRepositoryDeleteCmd())

	return cmd
}

func repositoriesRun(f *repositoryFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		// set up http
		client := &http.Client{}

		if f.apiVersion == "v4" {
			// set up the URL to query
			url := "/fmeapiv4/repositories"
			if f.name != "" {
				// add the repository name to the request if specified
				url = url + "/" + f.name
			}
			request, err := buildFmeServerRequest(url, "GET", nil)
			if err != nil {
				return err
			}

			q := request.URL.Query()

			if f.filterString != "" {
				// add the owner as a query parameter if specified
				q.Add("filterString", f.filterString)
				request.URL.RawQuery = q.Encode()
			}

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != http.StatusOK {
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

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result FMEServerRepositoriesV4

			if f.name == "" {
				// if no name specified, request will return the full struct
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				}
			} else {
				// else, we aree getting a single repository. We will just append this
				// to the Item list in the full struct for easier parsing
				var singleResult FMEServerRepositoryV4
				if err := json.Unmarshal(responseData, &singleResult); err != nil {
					return err
				}
				result.TotalCount = 1
				result.Items = append(result.Items, singleResult)
			}

			if f.outputType == "table" {

				t := table.NewWriter()
				t.SetStyle(defaultStyle)

				t.AppendHeader(table.Row{"Name", "Owner", "Description", "Workspaces"})

				for _, element := range result.Items {
					t.AppendRow(table.Row{element.Name, element.Owner, element.Description, element.WorkspaceCount})
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
		} else if f.apiVersion == "v3" {
			// set up the URL to query
			url := "/fmerest/v3/repositories"
			if f.name != "" {
				// add the repository name to the request if specified
				url = url + "/" + f.name
			}
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
				if response.StatusCode == http.StatusNotFound {
					return fmt.Errorf("%w: check that the specified repository exists", errors.New(response.Status))
				}
				return errors.New(response.Status)
			}

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result FMEServerRepositoriesV3

			if f.name == "" {
				// if no name specified, request will return the full struct
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				}
			} else {
				// else, we are getting a single repository. We will just append this
				// to the Item list in the full struct for easier parsing
				var singleResult FMEServerRepositoryV3
				if err := json.Unmarshal(responseData, &singleResult); err != nil {
					return err
				}
				result.TotalCount = 1
				result.Items = append(result.Items, singleResult)

			}

			if f.outputType == "table" {

				t := table.NewWriter()
				t.SetStyle(defaultStyle)

				t.AppendHeader(table.Row{"Name", "Owner", "Description", "Sharable"})

				for _, element := range result.Items {
					t.AppendRow(table.Row{element.Name, element.Owner, element.Description, element.Sharable})
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
		}
		return nil
	}
}
