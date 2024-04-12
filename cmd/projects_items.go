package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type projectItemFlags struct {
	id                  string
	name                string
	typeFlag            []string
	includeDependencies bool
	filterString        string
	filterProperty      []string
	outputType          string
	noHeaders           bool
}

type ProjectItemV4 struct {
	ID           string                    `json:"id"`
	Name         string                    `json:"name"`
	Type         string                    `json:"type"`
	Owner        string                    `json:"owner"`
	LastUpdated  time.Time                 `json:"lastUpdated"`
	Dependencies []ProjectItemDependencyV4 `json:"dependencies"`
}

type ProjectItemDependencyV4 struct {
	Name         string   `json:"name"`
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Dependencies []string `json:"dependencies"`
}
type ProjectItemsV4 struct {
	Items      []ProjectItemV4 `json:"items"`
	TotalCount int             `json:"totalCount"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
}

func newProjectItemsCmd() *cobra.Command {
	f := projectItemFlags{}
	cmd := &cobra.Command{
		Use:   "items",
		Short: "Lists the items for the specified project",
		Long:  `Lists the items contained in the specified project.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// can't specify both id and name
			if f.id == "" && f.name == "" {
				return errors.New("required flag(s) \"id\" or \"name\" not set")
			}

			// can't specify filter-property without filter-string
			if f.filterString == "" && len(f.filterProperty) > 0 {
				return errors.New("flag \"filter-property\" specified without flag \"filter-string\"")
			}

			return nil
		},
		Example: `
  # Get all items for a project via id
  fmeflow projects items --id a64297e7-a119-4e10-ac37-5d0bba12194b

  # Get all items for a project via name
  fmeflow projects items --name test_project

  # Get items with type workspace for a project via name
  fmeflow projects items --name test_project --type workspace
  
  # Get all items for a project via name without dependencies
  fmeflow projects items --name test_project --include-dependencies=false
  
  # Get all items for a project via name with a filter on name
  fmeflow projects items --name test_project --filter-string "test_name" --filter-properties "name"`,
		Args: NoArgs,
		RunE: projectItemsRun(&f),
	}

	cmd.Flags().StringVar(&f.id, "id", "", "Id of project to get items for ")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of project to get items for")
	cmd.Flags().StringArrayVar(&f.typeFlag, "type", []string{}, "Type of items to get")
	cmd.Flags().BoolVar(&f.includeDependencies, "include-dependencies", true, "Include dependencies in the output")
	cmd.Flags().StringVar(&f.filterString, "filter-string", "", "String to filter items by")
	cmd.Flags().StringArrayVar(&f.filterProperty, "filter-property", []string{}, "Property to filter by. Should be one of \"name\" or \"owner\". Can only be set if filter-string is also set")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")

	cmd.MarkFlagsMutuallyExclusive("id", "name")

	return cmd
}

func projectItemsRun(f *projectItemFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		// if name isn't specified, just get the id
		if f.name != "" {
			id, err := getProjectId(f.name)
			if err != nil {
				return err
			}
			f.id = id
		}

		client := &http.Client{}

		url := "/fmeapiv4/projects/" + f.id + "/items"
		request, err := buildFmeFlowRequest(url, "GET", nil)
		if err != nil {
			return err
		}

		q := request.URL.Query()
		// loop through type flags and add them to the query
		for _, t := range f.typeFlag {
			q.Add("type", t)
		}

		// add the include dependencies flag
		q.Add("includeDependencies", strconv.FormatBool(f.includeDependencies))

		// add the filter string and filter properties
		for _, t := range f.filterProperty {
			q.Add("filterProperties", t)
		}

		if f.filterString != "" {
			q.Add("filterString", f.filterString)
		}

		request.URL.RawQuery = q.Encode()

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

		// marshal into struct
		var projectItems ProjectItemsV4

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(responseData, &projectItems)
		if err != nil {
			return err
		}

		if f.outputType == "table" {

			t := table.NewWriter()
			t.SetStyle(defaultStyle)

			t.AppendHeader(table.Row{"ID", "Name", "Type", "Owner", "Last Updated"})

			for _, element := range projectItems.Items {
				t.AppendRow(table.Row{element.ID, element.Name, element.Type, element.Owner, element.LastUpdated})
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
			for _, element := range projectItems.Items {
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
}
