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

type FMEServerWorkspaces struct {
	Offset     int                  `json:"offset"`
	Limit      int                  `json:"limit"`
	TotalCount int                  `json:"totalCount"`
	Items      []FMEServerWorkspace `json:"items"`
}

type FMEServerWorkspace struct {
	LastSaveDate    time.Time `json:"lastSaveDate"`
	AvgCPUPct       float64   `json:"avgCpuPct"`
	AvgPeakMemUsage int       `json:"avgPeakMemUsage"`
	Description     string    `json:"description"`
	RepositoryName  string    `json:"repositoryName"`
	Title           string    `json:"title"`
	Type            string    `json:"type"`
	UserName        string    `json:"userName"`
	FileCount       int       `json:"fileCount"`
	AvgCPUTime      int       `json:"avgCpuTime"`
	LastPublishDate time.Time `json:"lastPublishDate"`
	Name            string    `json:"name"`
	TotalFileSize   int       `json:"totalFileSize"`
	TotalRuns       int       `json:"totalRuns"`
	AvgElapsedTime  int       `json:"avgElapsedTime"`
}

type FMEServerWorkspaceDetailed struct {
	LegalTermsConditions string  `json:"legalTermsConditions"`
	AvgCPUPct            float64 `json:"avgCpuPct"`
	Usage                string  `json:"usage"`
	AvgPeakMemUsage      int     `json:"avgPeakMemUsage"`
	Description          string  `json:"description"`
	Datasets             struct {
		Destination []struct {
			Format       string `json:"format"`
			Name         string `json:"name"`
			Location     string `json:"location"`
			Source       bool   `json:"source"`
			Featuretypes []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Attributes  []struct {
					Decimals int    `json:"decimals"`
					Name     string `json:"name"`
					Width    int    `json:"width"`
					Type     string `json:"type"`
				} `json:"attributes"`
				Properties []interface{} `json:"properties"`
			} `json:"featuretypes"`
			Properties []struct {
				Name       string `json:"name"`
				Attributes struct {
				} `json:"attributes"`
				Category string `json:"category"`
				Value    string `json:"value"`
			} `json:"properties"`
		} `json:"destination"`
		Source []struct {
			Format       string `json:"format"`
			Name         string `json:"name"`
			Location     string `json:"location"`
			Source       bool   `json:"source"`
			Featuretypes []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Attributes  []struct {
					Decimals int    `json:"decimals"`
					Name     string `json:"name"`
					Width    int    `json:"width"`
					Type     string `json:"type"`
				} `json:"attributes"`
				Properties []interface{} `json:"properties"`
			} `json:"featuretypes"`
			Properties []struct {
				Name       string `json:"name"`
				Attributes struct {
				} `json:"attributes"`
				Category string `json:"category"`
				Value    string `json:"value"`
			} `json:"properties"`
		} `json:"source"`
	} `json:"datasets"`
	Title           string    `json:"title"`
	Type            string    `json:"type"`
	BuildNumber     int       `json:"buildNumber"`
	Enabled         bool      `json:"enabled"`
	AvgCPUTime      int       `json:"avgCpuTime"`
	LastPublishDate time.Time `json:"lastPublishDate"`
	LastSaveBuild   string    `json:"lastSaveBuild"`
	AvgElapsedTime  int       `json:"avgElapsedTime"`
	LastSaveDate    time.Time `json:"lastSaveDate"`
	Requirements    string    `json:"requirements"`
	Resources       []struct {
		Size        int    `json:"size"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"resources"`
	History  string `json:"history"`
	Services []struct {
		DisplayName string `json:"displayName"`
		Name        string `json:"name"`
	} `json:"services"`
	UserName            string        `json:"userName"`
	RequirementsKeyword string        `json:"requirementsKeyword"`
	FileSize            int           `json:"fileSize"`
	Name                string        `json:"name"`
	TotalRuns           int           `json:"totalRuns"`
	Category            string        `json:"category"`
	Parameters          []interface{} `json:"parameters"`
	Properties          []struct {
		Name       string `json:"name"`
		Attributes struct {
		} `json:"attributes"`
		Category string `json:"category"`
		Value    string `json:"value"`
	} `json:"properties"`
}

type workspaceFlags struct {
	repository string
	name       string
	outputType string
	noHeaders  bool
}

func newWorkspaceCmd() *cobra.Command {
	f := workspaceFlags{}
	cmd := &cobra.Command{
		Use:   "workspaces",
		Short: "List workspaces by repository",
		Long:  `Lists workspaces on the given FME Server in the repository.`,
		Example: `
	Examples:
	# List all workspaces in Samples repository
	fmeserver workspaces --repository Samples
	
	# List all workspaces in the Samples repository and output it in json
	fmeserver workspaces --repository Samples --json
	
	# List all workspaces in the Samples repository with custom columns showing the last publish date and number of times run
	fmeserver workspaces --repository Samples --output="custom-columns=NAME:$.name,PUBLISH DATE:$.lastPublishDate,TOTAL RUNS:$.totalRuns"
	
	# Get information on a single workspace 
	fmeserver workspaces --repository Samples --name austinApartments.fmw
	
	# Get the name, source format, destination format, and the services this workspace is assigned to
	fmeserver workspaces --repository Samples --name austinApartments.fmw --output=custom-columns=NAME:$.name,SOURCE:$.datasets.source[*].format,DEST:$.datasets.destination[*].format,SERVICES:$.services[*].name`,
		Args: NoArgs,
		RunE: workspacesRun(&f),
	}

	cmd.Flags().StringVar(&f.repository, "repository", "", "Name of repository to list workspaces in.")
	cmd.Flags().StringVar(&f.name, "name", "", "If specified, get details about a specific workspace")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.MarkFlagRequired("repository")
	return cmd
}

func workspacesRun(f *workspaceFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		// set up http
		client := &http.Client{}

		// set up the URL to query
		url := "/fmerest/v3/repositories/" + f.repository + "/items"

		if f.name != "" {
			url += "/" + f.name
		}

		request, err := buildFmeServerRequest(url, "GET", nil)
		if err != nil {
			return err
		}

		q := request.URL.Query()
		q.Add("type", "WORKSPACE")

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusOK {
			if response.StatusCode == http.StatusNotFound {
				if f.name == "" {
					return fmt.Errorf("%w: check that the specified repository exists", errors.New(response.Status))
				} else {
					return fmt.Errorf("%w: check that the specified repository and workspace exists", errors.New(response.Status))
				}
			}
			return errors.New(response.Status)
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var result FMEServerWorkspaces
		var resultDetailed FMEServerWorkspaceDetailed
		if f.name == "" {
			// if no name specified, request will return the full struct
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			}
		} else {
			// if no name specified, request will return the full struct
			if err := json.Unmarshal(responseData, &resultDetailed); err != nil {
				return err
			}
		}

		if f.outputType == "table" {

			t := table.NewWriter()
			t.SetStyle(defaultStyle)

			t.AppendHeader(table.Row{"Name", "Title", "Last Save Date"})

			if f.name == "" {
				for _, element := range result.Items {
					t.AppendRow(table.Row{element.Name, element.Title, element.LastSaveDate})
				}
			} else {
				t.AppendRow(table.Row{resultDetailed.Name, resultDetailed.Title, resultDetailed.LastSaveDate})
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
			if f.name == "" {
				for _, element := range result.Items {
					mJson, err := json.Marshal(element)
					if err != nil {
						return err
					}
					marshalledItems = append(marshalledItems, mJson)
				}
			} else {
				mJson, err := json.Marshal(resultDetailed)
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
