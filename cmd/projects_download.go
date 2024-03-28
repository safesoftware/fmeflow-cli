package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectsDownloadFlags struct {
	id                     string
	file                   string
	name                   string
	excludeSensitiveInfo   bool
	suppressFileRename     bool
	excludeSelectableItems bool
	apiVersion             apiVersionFlag
}

type projectExport struct {
	IncludeSensitiveInfo      bool                        `json:"includeSensitiveInfo"`
	ExportPackageName         string                      `json:"exportPackageName"`
	SelectedItems             []projectExportSelectedItem `json:"selectedItems"`
	ExcludeAllSelectableItems bool                        `json:"excludeAllSelectableItems"`
}

type projectExportSelectedItem struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// backupCmd represents the backup command
func newProjectDownloadCmd() *cobra.Command {
	f := projectsDownloadFlags{}
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Downloads an FME Server Project",
		Long:  "Downloads an FME Server Project to a local file. Useful for backing up or moving a project to another FME Server.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// get build to decide if we should use v3 or v4
			if f.apiVersion == "" {
				fmeflowBuild := viper.GetInt("build")
				if fmeflowBuild < projectUploadV4BuildThreshold {
					f.apiVersion = apiVersionFlagV3
				} else {
					f.apiVersion = apiVersionFlagV4
				}
			}

			// ensure one of either id or name is specified
			if f.apiVersion == apiVersionFlagV4 {
				if f.id == "" && f.name == "" {
					return errors.New("required flag(s) \"id\" or \"name\" not set")
				}
			} else {
				if f.name == "" {
					return errors.New("required flag(s) \"name\" not set")
				}
			}

			// only specify id if using v4
			if f.apiVersion == apiVersionFlagV3 && f.id != "" {
				return errors.New("id flag is only supported with v4")
			}
			return nil
		},
		Example: `
  # download a project named "Test Project" to a local file with default name
  fmeflow projects download --name "Test Project"
	
  # download a project named "Test Project" to a local file named MyProject.fsproject
  fmeflow projects download --name "Test Project" -f MyProject.fsproject`,
		Args: NoArgs,
		RunE: projectDownloadRun(&f),
	}
	cmd.Flags().StringVarP(&f.file, "file", "f", "ProjectPackage.fsproject", "Path to file to download the backup to.")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the project to download.")
	cmd.Flags().StringVar(&f.id, "id", "", "ID of the project to download.")
	cmd.Flags().BoolVar(&f.excludeSensitiveInfo, "exclude-sensitive-info", false, "Whether to exclude sensitive information from the exported package. Sensitive information will be excluded from connections, subscriptions, publications, schedule tasks, S3 resources, and user accounts. Other items in the project may still contain sensitive data, especially workspaces. Please be careful before sharing the project export pacakge with others.")
	cmd.Flags().BoolVar(&f.suppressFileRename, "suppress-file-rename", false, "Specify this flag to not add .fsproject to the output file automatically")
	cmd.Flags().BoolVar(&f.excludeSelectableItems, "exclude-selectable-items", false, "Excludes all selectable item in this package. Default is false.")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4.")
	cmd.Flags().MarkHidden("suppress-file-rename")
	cmd.Flags().MarkHidden("api-version")
	cmd.MarkFlagsMutuallyExclusive("name", "id")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	return cmd
}

func projectDownloadRun(f *projectsDownloadFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}
		var request http.Request

		// massage the backup file name
		if !f.suppressFileRename && f.file != "" {
			backupExtension := ".fsproject"
			if !strings.HasSuffix(f.file, backupExtension) {
				f.file += backupExtension
			}
		}

		if f.apiVersion == "v4" {

			// if name isn't empty, we have to first get the id for this project
			if f.name != "" {
				id, err := getProjectId(f.name)
				if err != nil {
					return err
				}
				f.id = id
			}

			// create the export object
			var export projectExport
			export.ExcludeAllSelectableItems = f.excludeSelectableItems
			export.ExportPackageName = f.file
			export.IncludeSensitiveInfo = !f.excludeSensitiveInfo

			url := "/fmeapiv4/projects/" + f.id + "/export/download"
			jsonData, err := json.Marshal(export)
			if err != nil {
				return err
			}

			request, err = buildFmeFlowRequest(url, "POST", bytes.NewBuffer(jsonData))
			if err != nil {
				return err
			}

			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Accept", "application/octet-stream")

		} else {

			// add mandatory values
			data := url.Values{
				"exportPackageName":    {f.file},
				"excludeSensitiveInfo": {strconv.FormatBool(f.excludeSensitiveInfo)},
			}

			var err error
			request, err = buildFmeFlowRequest("/fmerest/v3/projects/projects/"+f.name+"/export/download", "POST", strings.NewReader(data.Encode()))
			if err != nil {
				return err
			}

			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			request.Header.Add("Accept", "application/octet-stream")

			fmt.Fprintln(cmd.OutOrStdout(), "Downloading project file...")

		}

		// setting up the request is different in v4 vs v3, but the execution from here on should be the same
		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 200 {
			if response.StatusCode == http.StatusUnprocessableEntity {
				return fmt.Errorf("%w: check that the specified project exists", errors.New(response.Status))
			}
			return errors.New(response.Status)
		}
		defer response.Body.Close()

		// Create the output file
		out, err := os.Create(f.file)
		if err != nil {
			return err
		}
		defer out.Close()

		// use Copy so that it doesn't store the entire file in memory
		_, err = io.Copy(out, response.Body)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Project exported to "+f.file)
		return nil
	}
}
