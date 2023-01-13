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

type restoreFlags struct {
	file               string
	importMode         string
	pauseNotifications bool
	projectsImportMode string
	resource           bool
	resourceName       string
	failureTopic       string
	successTopic       string
}

type RestoreResource struct {
	Id int `json:"id"`
}

func newRestoreCmd() *cobra.Command {
	f := restoreFlags{}
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restores the FME Server configuration from an import package",
		Long:  "Restores the FME Server configuration from an import package",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// ensure one of --file or --resource is set
			if f.file == "" && !f.resource {
				return errors.New("required flag \"file\" or \"resource\" not set")
			}

			// verify import mode is valid
			if f.importMode != "UPDATE" && f.importMode != "INSERT" {
				return errors.New("invalid import-mode. Must be either UPDATE or INSERT")
			}

			// verify projects import mode is valid
			if f.projectsImportMode != "UPDATE" && f.projectsImportMode != "INSERT" && f.projectsImportMode != "" {
				return errors.New("invalid projects-import-mode. Must be either UPDATE or INSERT")
			}

			// if restoring from the shared resource and no file is set, default to ServerConfigPackage.fsconfig
			if f.resource && f.file == "" {
				f.file = "ServerConfigPackage.fsconfig"
			}

			return nil
		},
		Example: `
  # Restore from a backup in a local file
  fmeserver restore --file ServerConfigPackage.fsconfig

  # Restore from a backup in a local file using UPDATE mode
  fmeserver restore --file ServerConfigPackage.fsconfig --import-mode UPDATE`,
		Args: NoArgs,
		RunE: restoreRun(&f),
	}

	cmd.Flags().StringVarP(&f.file, "file", "f", "", "Path to backup file to upload to restore. Can be a local file or the relative path inside the specified shared resource.")
	cmd.Flags().StringVar(&f.importMode, "import-mode", "INSERT", "To import only items in the import package that do not exist on the current instance, specify INSERT. To overwrite items on the current instance with those in the import package, specify UPDATE. Default is INSERT.")
	cmd.Flags().BoolVar(&f.pauseNotifications, "pause-notifications", true, "Disable notifications for the duration of the restore.")
	cmd.Flags().StringVar(&f.projectsImportMode, "projects-import-mode", "", "Import mode for projects. To import only projects in the import package that do not exist on the current instance, specify INSERT. To overwrite projects on the current instance with those in the import package, specify UPDATE. If not supplied, importMode will be used.")
	cmd.Flags().BoolVar(&f.resource, "resource", false, "Restore from a shared resource location instead of a local file.")
	cmd.Flags().StringVar(&f.resourceName, "resource-name", "FME_SHAREDRESOURCE_BACKUP", "Resource containing the import package.")
	cmd.Flags().StringVar(&f.failureTopic, "failure-topic", "", "Topic to notify on failure of the import. Default is MIGRATION_ASYNC_JOB_FAILURE")
	cmd.Flags().StringVar(&f.successTopic, "success-topic", "", "Topic to notify on success of the import. Default is MIGRATION_ASYNC_JOB_SUCCESS")

	return cmd
}
func restoreRun(f *restoreFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := &http.Client{}

		url := ""
		var request http.Request

		if !f.resource {
			file, err := os.Open(f.file)
			if err != nil {
				return err
			}
			defer file.Close()

			url = "/fmerest/v3/migration/restore/upload"
			request, err = buildFmeServerRequest(url, "POST", file)
			if err != nil {
				return err
			}
			request.Header.Set("Content-Type", "application/octet-stream")
		} else {
			url = "/fmerest/v3/migration/restore/resource"
			var err error
			request, err = buildFmeServerRequest(url, "POST", nil)
			if err != nil {
				return err
			}
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		}

		q := request.URL.Query()

		if f.resourceName != "" {
			q.Add("resourceName", f.resourceName)
		}

		if f.resource {
			q.Add("importPackage", f.file)
		}

		if f.pauseNotifications {
			q.Add("pauseNotifications", strconv.FormatBool(f.pauseNotifications))
		}

		if f.importMode != "" {
			q.Add("importMode", f.importMode)
		}

		if f.projectsImportMode != "" {
			q.Add("projectsImportMode", f.projectsImportMode)
		}

		if f.successTopic != "" {
			q.Add("successTopic", f.successTopic)
		}

		if f.failureTopic != "" {
			q.Add("failureTopic", f.failureTopic)
		}

		request.URL.RawQuery = q.Encode()

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if !f.resource && response.StatusCode != http.StatusOK {
			return errors.New(response.Status)
		} else if f.resource && response.StatusCode != http.StatusAccepted {
			return errors.New(response.Status)
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result RestoreResource
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), "Restore task submitted with id: "+strconv.Itoa(result.Id))
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
