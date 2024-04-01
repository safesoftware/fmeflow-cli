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

type connectionsFlags struct {
	name           string
	typeConnection []string
	excludedType   []string
	category       []string
	outputType     string
	noHeaders      bool
}

type FMEFlowConnections struct {
	Items      []Connection `json:"items"`
	TotalCount int          `json:"totalCount"`
	Limit      int          `json:"limit"`
	Offset     int          `json:"offset"`
}

type Connection struct {
	Name       string                 `json:"name"`
	Category   string                 `json:"category"`
	Type       string                 `json:"type"`
	Owner      string                 `json:"owner"`
	Shareable  bool                   `json:"shareable"`
	Parameters map[string]interface{} `json:"parameters"`
}

func newConnectionsCmd() *cobra.Command {
	f := connectionsFlags{}
	cmd := &cobra.Command{
		Use:   "connections",
		Short: "Lists connections on FME Flow",
		Long:  "Lists connections on FME Flow. Pass in a name to retrieve information on a single project.",
		Example: `
  # List all projects
  fmeflow projects

  # List all projects owned by the user admin
  fmeflow projects --owner admin`,
		Args: NoArgs,
		RunE: connectionsRun(&f),
	}

	cmd.Flags().StringVar(&f.name, "name", "", "Return a single project with the given name.")
	cmd.Flags().StringArrayVar(&f.typeConnection, "type", []string{}, "The types of connections to return. Can be passed in multiple times")
	cmd.Flags().StringArrayVar(&f.excludedType, "excluded-type", []string{}, "The types of connections to exclude. Can be passed in multiple times")
	cmd.Flags().StringArrayVar(&f.category, "category", []string{}, "The categories of connections to return. Can be passed in multiple times")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.AddCommand(newProjectDownloadCmd())
	cmd.AddCommand(newProjectUploadCmd())

	return cmd
}
func connectionsRun(f *connectionsFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		client := &http.Client{}

		url := "/fmeapiv4/connections"
		if f.name != "" {
			url = url + "/" + f.name
		}

		request, err := buildFmeFlowRequest(url, "GET", nil)
		if err != nil {
			return err
		}

		q := request.URL.Query()

		for _, t := range f.typeConnection {
			q.Add("types", t)
		}

		for _, t := range f.excludedType {
			q.Add("excludedTypes", t)
		}

		for _, c := range f.category {
			q.Add("categories", c)
		}

		//q.Add("filterString", f.owner)
		//q.Add("filterProperties", "owner")
		request.URL.RawQuery = q.Encode()

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
		var result FMEFlowConnections

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if f.name == "" {
			// if no name specified, request will return the full struct
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			}
		} else {
			var connectionStruct Connection

			if err := json.Unmarshal(responseData, &connectionStruct); err != nil {
				return err
			}
			result.TotalCount = 1
			result.Items = append(result.Items, connectionStruct)
		}

		if f.outputType == "table" {

			t := table.NewWriter()
			t.SetStyle(defaultStyle)

			t.AppendHeader(table.Row{"Name", "Type", "Category"})

			for _, element := range result.Items {
				t.AppendRow(table.Row{element.Name, element.Type, element.Category})
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
}
