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
)

type FMEServerRepositories struct {
	Offset     int                   `json:"offset"`
	Limit      int                   `json:"limit"`
	TotalCount int                   `json:"totalCount"`
	Items      []FMEServerRepository `json:"items"`
}

type FMEServerRepository struct {
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sharable    bool   `json:"sharable"`
}

type repositoryFlags struct {
	owner      string
	name       string
	outputType string
	noHeaders  bool
}

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
	fmeserver repositories --output=custom-columns=NAME:$.name --no-headers
	
	# Output all repositories in json format
	fmeserver repositories --json`,
		Args: NoArgs,
		RunE: repositoriesRun(&f),
	}

	cmd.Flags().StringVar(&f.owner, "owner", "", "If specified, only repositories owned by the specified user will be returned.")
	cmd.Flags().StringVar(&f.name, "name", "", "If specified, only the repository with that name will be returned")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
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

		var result FMEServerRepositories

		if f.name == "" {
			// if no name specified, request will return the full struct
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			}
		} else {
			// else, we aree getting a single repository. We will just append this
			// to the Item list in the full struct for easier parsing
			var singleResult FMEServerRepository
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
		return nil
	}
}
