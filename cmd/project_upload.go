package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type projectUploadFlags struct {
	file                string
	importMode          string
	pauseNotifications  bool
	projectsImportMode  string
	disableProjectItems bool
}

type ProjectUploadTask struct {
	Id int `json:"id"`
}

func newProjectUploadCmd() *cobra.Command {
	f := projectUploadFlags{}
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Imports FME Server Projects from a downloaded import package.",
		Long:  "Imports FME Server Projects from a downloaded import package.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// verify import mode is valid
			if f.importMode != "UPDATE" && f.importMode != "INSERT" {
				return errors.New("invalid import-mode. Must be either UPDATE or INSERT")
			}

			// verify projects import mode is valid
			if f.projectsImportMode != "UPDATE" && f.projectsImportMode != "INSERT" && f.projectsImportMode != "" {
				return errors.New("invalid projects-import-mode. Must be either UPDATE or INSERT")
			}

			return nil
		},
		Example: `
  # Restore from a backup in a local file
  fmeserver restore --file ServerConfigPackage.fsconfig

  # Restore from a backup in a local file using UPDATE mode
  fmeserver restore --file ServerConfigPackage.fsconfig --import-mode UPDATE`,
		Args: NoArgs,
		RunE: projectUploadRun(&f),
	}

	cmd.Flags().StringVarP(&f.file, "file", "f", "", "Path to backup file to upload to restore. Can be a local file or the relative path inside the specified shared resource.")
	cmd.Flags().StringVar(&f.importMode, "import-mode", "INSERT", "To import only items in the import package that do not exist on the current instance, specify INSERT. To overwrite items on the current instance with those in the import package, specify UPDATE. Default is INSERT.")
	cmd.Flags().BoolVar(&f.pauseNotifications, "pause-notifications", true, "Disable notifications for the duration of the restore.")
	cmd.Flags().StringVar(&f.projectsImportMode, "projects-import-mode", "", "Import mode for projects. To import only projects in the import package that do not exist on the current instance, specify INSERT. To overwrite projects on the current instance with those in the import package, specify UPDATE. If not supplied, importMode will be used.")
	cmd.Flags().BoolVar(&f.disableProjectItems, "disable-project-items", false, "Whether to disable items in the imported FME Server Projects. If true, items that are new or overwritten will be imported but disabled. If false, project items are imported as defined in the import package.")
	cmd.MarkFlagRequired("file")

	return cmd
}
func projectUploadRun(f *projectUploadFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := &http.Client{}

		url := ""
		var request http.Request
		file, err := os.Open(f.file)
		if err != nil {
			return err
		}
		defer file.Close()

		url = "/fmerest/v3/projects/import/upload"
		request, err = buildFmeServerRequest(url, "POST", file)
		if err != nil {
			return err
		}
		request.Header.Set("Content-Type", "application/octet-stream")

		q := request.URL.Query()

		if f.pauseNotifications {
			q.Add("pauseNotifications", strconv.FormatBool(f.pauseNotifications))
		}

		if f.importMode != "" {
			q.Add("importMode", f.importMode)
		}

		if f.projectsImportMode != "" {
			q.Add("projectsImportMode", f.projectsImportMode)
		}

		if f.disableProjectItems {
			q.Add("disableProjectItems", strconv.FormatBool(f.disableProjectItems))
		}

		request.URL.RawQuery = q.Encode()

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

		var result ProjectUploadTask
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), "Project Upload task submitted with id: "+strconv.Itoa(result.Id))
			} else {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
			}
		}

		return nil
	}
}
