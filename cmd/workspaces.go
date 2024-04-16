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

type FMEFlowWorkspacesV4 struct {
	Items []struct {
		AverageCPUPercent      float64   `json:"averageCpuPercent"`
		AverageCPUTime         float64   `json:"averageCpuTime"`
		AverageElapsedTime     float64   `json:"averageElapsedTime"`
		AveragePeakMemoryUsage int       `json:"averagePeakMemoryUsage"`
		Description            string    `json:"description"`
		Favorite               bool      `json:"favorite"`
		FileCount              int       `json:"fileCount"`
		LastPublishDate        time.Time `json:"lastPublishDate"`
		LastPublishUser        string    `json:"lastPublishUser"`
		LastPublishUserID      string    `json:"lastPublishUserId"`
		LastSaveDate           time.Time `json:"lastSaveDate"`
		Name                   string    `json:"name"`
		RepositoryName         string    `json:"repositoryName"`
		Title                  string    `json:"title"`
		TotalFileSize          int       `json:"totalFileSize"`
		TotalRuns              int       `json:"totalRuns"`
		Type                   string    `json:"type"`
	} `json:"items"`
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
	TotalCount int `json:"totalCount"`
}

type FMEFlowWorkspaceDetailedV4 struct {
	AverageCPUPercent      float64 `json:"averageCpuPercent"`
	AverageCPUTime         float64 `json:"averageCpuTime"`
	AverageElapsedTime     float64 `json:"averageElapsedTime"`
	AveragePeakMemoryUsage int     `json:"averagePeakMemoryUsage"`
	BuildNumber            int     `json:"buildNumber"`
	Category               string  `json:"category"`
	Datasets               struct {
		Destination []struct {
			FeatureTypes []struct {
				Attributes []struct {
					Decimals int    `json:"decimals"`
					Name     string `json:"name"`
					Type     string `json:"type"`
					Width    int    `json:"width"`
				} `json:"attributes"`
				Description string `json:"description"`
				Name        string `json:"name"`
				Properties  []struct {
					Attributes struct {
						AdditionalProp1 string `json:"additionalProp1"`
						AdditionalProp2 string `json:"additionalProp2"`
						AdditionalProp3 string `json:"additionalProp3"`
					} `json:"attributes"`
					Category string `json:"category"`
					Name     string `json:"name"`
					Value    string `json:"value"`
				} `json:"properties"`
			} `json:"featureTypes"`
			Format     string `json:"format"`
			Location   string `json:"location"`
			Name       string `json:"name"`
			Properties []struct {
				Attributes struct {
					AdditionalProp1 string `json:"additionalProp1"`
					AdditionalProp2 string `json:"additionalProp2"`
					AdditionalProp3 string `json:"additionalProp3"`
				} `json:"attributes"`
				Category string `json:"category"`
				Name     string `json:"name"`
				Value    string `json:"value"`
			} `json:"properties"`
			Source bool `json:"source"`
		} `json:"destination"`
		Source []struct {
			FeatureTypes []struct {
				Attributes []struct {
					Decimals int    `json:"decimals"`
					Name     string `json:"name"`
					Type     string `json:"type"`
					Width    int    `json:"width"`
				} `json:"attributes"`
				Description string `json:"description"`
				Name        string `json:"name"`
				Properties  []struct {
					Attributes struct {
						AdditionalProp1 string `json:"additionalProp1"`
						AdditionalProp2 string `json:"additionalProp2"`
						AdditionalProp3 string `json:"additionalProp3"`
					} `json:"attributes"`
					Category string `json:"category"`
					Name     string `json:"name"`
					Value    string `json:"value"`
				} `json:"properties"`
			} `json:"featureTypes"`
			Format     string `json:"format"`
			Location   string `json:"location"`
			Name       string `json:"name"`
			Properties []struct {
				Attributes struct {
					AdditionalProp1 string `json:"additionalProp1"`
					AdditionalProp2 string `json:"additionalProp2"`
					AdditionalProp3 string `json:"additionalProp3"`
				} `json:"attributes"`
				Category string `json:"category"`
				Name     string `json:"name"`
				Value    string `json:"value"`
			} `json:"properties"`
			Source bool `json:"source"`
		} `json:"source"`
	} `json:"datasets"`
	Description          string    `json:"description"`
	Favorite             bool      `json:"favorite"`
	FileSize             int       `json:"fileSize"`
	History              string    `json:"history"`
	LastPublishDate      time.Time `json:"lastPublishDate"`
	LastSaveBuild        string    `json:"lastSaveBuild"`
	LastSaveDate         time.Time `json:"lastSaveDate"`
	LegalTermsConditions string    `json:"legalTermsConditions"`
	Name                 string    `json:"name"`
	Parameters           []struct {
		AdditionalProp1 struct {
		} `json:"additionalProp1"`
		AdditionalProp2 struct {
		} `json:"additionalProp2"`
		AdditionalProp3 struct {
		} `json:"additionalProp3"`
	} `json:"parameters"`
	Properties []struct {
		Attributes struct {
			AdditionalProp1 string `json:"additionalProp1"`
			AdditionalProp2 string `json:"additionalProp2"`
			AdditionalProp3 string `json:"additionalProp3"`
		} `json:"attributes"`
		Category string `json:"category"`
		Name     string `json:"name"`
		Value    string `json:"value"`
	} `json:"properties"`
	Requirements        string `json:"requirements"`
	RequirementsKeyword string `json:"requirementsKeyword"`
	Resources           []struct {
		Name string `json:"name"`
		Size int    `json:"size"`
	} `json:"resources"`
	Services struct {
		DataDownload struct {
			Reader     string   `json:"reader"`
			Registered bool     `json:"registered"`
			Writers    []string `json:"writers"`
			ZipLayout  struct {
				AdditionalProp1 string `json:"additionalProp1"`
				AdditionalProp2 string `json:"additionalProp2"`
				AdditionalProp3 string `json:"additionalProp3"`
			} `json:"zipLayout"`
		} `json:"dataDownload"`
		DataStreaming struct {
			Reader     string   `json:"reader"`
			Registered bool     `json:"registered"`
			Writers    []string `json:"writers"`
		} `json:"dataStreaming"`
		JobSubmitter struct {
			Reader     string `json:"reader"`
			Registered bool   `json:"registered"`
		} `json:"jobSubmitter"`
		KmlNetworkLink struct {
			Description  string `json:"description"`
			LatLonAltBox struct {
				East  int `json:"east"`
				North int `json:"north"`
				South int `json:"south"`
				West  int `json:"west"`
			} `json:"latLonAltBox"`
			Link struct {
				RefreshMode         string `json:"refreshMode"`
				ViewFormat          string `json:"viewFormat"`
				ViewRefreshInterval int    `json:"viewRefreshInterval"`
				ViewRefreshMode     string `json:"viewRefreshMode"`
				ViewRefreshTime     int    `json:"viewRefreshTime"`
			} `json:"link"`
			Lod struct {
				MaxLodPixels int `json:"maxLodPixels"`
				MinLodPixels int `json:"minLodPixels"`
			} `json:"lod"`
			Name       string   `json:"name"`
			Registered bool     `json:"registered"`
			Visibility string   `json:"visibility"`
			Writers    []string `json:"writers"`
		} `json:"kmlNetworkLink"`
	} `json:"services"`
	Title     string `json:"title"`
	TotalRuns int    `json:"totalRuns"`
	Type      string `json:"type"`
	Usage     string `json:"usage"`
	UserName  string `json:"userName"`
}

type FMEFlowWorkspacesV3 struct {
	Offset     int                  `json:"offset"`
	Limit      int                  `json:"limit"`
	TotalCount int                  `json:"totalCount"`
	Items      []FMEFlowWorkspaceV3 `json:"items"`
}

type FMEFlowWorkspaceV3 struct {
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

type FMEFlowWorkspaceDetailedV3 struct {
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
	repository   string
	name         string
	outputType   string
	filterString string
	noHeaders    bool
	apiVersion   apiVersionFlag
}

func newWorkspaceCmd() *cobra.Command {
	f := workspaceFlags{}
	cmd := &cobra.Command{
		Use:   "workspaces",
		Short: "List workspaces.",
		Long:  `Lists workspaces that exist on the FME Server. Filter by repository, specify a name to retrieve a specific workspace, or specify a filter string to narrow down by name or title.`,
		Example: `
  # List all workspaces on the FME Server
  fmeflow workspaces
	
  # List all workspaces in Samples repository
  fmeflow workspaces --repository Samples
	
  # List all workspaces in the Samples repository and output it in json
  fmeflow workspaces --repository Samples --json
	
  # List all workspaces in the Samples repository with custom columns showing the last publish date and number of times run
  fmeflow workspaces --repository Samples --output="custom-columns=NAME:.name,PUBLISH DATE:.lastPublishDate,TOTAL RUNS:.totalRuns"
	
  # Get information on a single workspace 
  fmeflow workspaces --repository Samples --name austinApartments.fmw
	
  # Get the name, source format, and destination format for this workspace
  fmeflow workspaces --repository Samples --name austinApartments.fmw --output=custom-columns=NAME:.name,SOURCE:.datasets.source[*].format,DEST:.datasets.destination[*].format`,
		Args: NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// get build to decide if we should use v3 or v4
			// FME Server 2023.0 and later can use v4. Otherwise fall back to v3
			if f.apiVersion == "" {
				fmeflowBuild := viper.GetInt("build")
				if fmeflowBuild < repositoriesV4BuildThreshold {
					f.apiVersion = apiVersionFlagV3
				} else {
					f.apiVersion = apiVersionFlagV4
				}
			}
			if f.apiVersion == apiVersionFlagV3 {
				if f.filterString != "" {
					return errors.New("cannot set the filter-string flag when using the V3 API")
				}
				cmd.MarkFlagRequired("repository")
			} else {
				if f.name != "" {
					cmd.MarkFlagRequired("repository")
				}
			}

			return nil
		},
		RunE: workspacesRun(&f),
	}

	cmd.Flags().StringVar(&f.repository, "repository", "", "Name of repository to list workspaces in.")
	cmd.Flags().StringVar(&f.name, "name", "", "If specified, get details about a specific workspace")
	cmd.Flags().StringVar(&f.filterString, "filter-string", "", "If specified, only workspaces with a matching name or title will be returned. Only usable with V4 API.")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
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

		if f.apiVersion == "v4" {
			// set up the URL to query
			url := "/fmeapiv4/workspaces"
			if f.name != "" {
				// add the repository name to the request if specified
				url = url + "/" + f.repository + "/" + f.name
			}
			request, err := buildFmeFlowRequest(url, "GET", nil)
			if err != nil {
				return err
			}

			q := request.URL.Query()

			if f.filterString != "" {
				// add the owner as a query parameter if specified
				q.Add("filterString", f.filterString)
			}

			if f.repository != "" {
				// add the owner as a query parameter if specified
				q.Add("repository", f.repository)
			}

			request.URL.RawQuery = q.Encode()

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

			var result FMEFlowWorkspacesV4
			var resultDetailed FMEFlowWorkspaceDetailedV4

			if f.name == "" {
				// if no name specified, request will return the full struct
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				}
			} else {
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
				if f.noHeaders {
					t.ResetHeaders()
				}
				fmt.Fprintln(cmd.OutOrStdout(), t.Render())

			} else {
				return errors.New("invalid output format specified")
			}
		} else if f.apiVersion == "v3" {

			// set up the URL to query
			url := "/fmerest/v3/repositories/" + f.repository + "/items"

			if f.name != "" {
				url += "/" + f.name
			}

			request, err := buildFmeFlowRequest(url, "GET", nil)
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
			var result FMEFlowWorkspacesV3
			var resultDetailed FMEFlowWorkspaceDetailedV3
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
				if f.noHeaders {
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
