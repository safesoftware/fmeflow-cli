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
	restoreBackupFile         string
	restoreImportMode         string
	restorePauseNotifications bool
	restoreProjectsImportMode string
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
		Example: `
  # Restore from a backup in a local file
  fmeserver restore --file .\ServerConfigPackage.fsconfig

  # Restore from a backup in a local file using UPDATE mode
  fmeserver restore --file .\ServerConfigPackage.fsconfig --import-mode UPDATE`,
		Args: NoArgs,
		RunE: restoreRun(&f),
	}

	cmd.Flags().StringVarP(&f.restoreBackupFile, "file", "f", "", "Path to backup file to upload to restore.")
	cmd.Flags().StringVar(&f.restoreImportMode, "import-mode", "INSERT", "To import only items in the import package that do not exist on the current instance, specify INSERT. To overwrite items on the current instance with those in the import package, specify UPDATE. Default is INSERT.")
	cmd.Flags().BoolVar(&f.restorePauseNotifications, "pause-notifications", true, "Disable notifications for the duration of the restore.")
	cmd.Flags().StringVar(&f.restoreProjectsImportMode, "projects-import-mode", "", "Import mode for projects. To import only projects in the import package that do not exist on the current instance, specify INSERT. To overwrite projects on the current instance with those in the import package, specify UPDATE. If not supplied, importMode will be used.")
	cmd.MarkFlagRequired("file")

	return cmd
}
func restoreRun(f *restoreFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := &http.Client{}

		file, err := os.Open(f.restoreBackupFile)
		if err != nil {
			return err
		}
		defer file.Close()

		// verify import mode is valid
		if f.restoreImportMode != "UPDATE" && f.restoreImportMode != "INSERT" {
			return errors.New("invalid import-mode. Must be either UPDATE or INSERT")
		}

		// verify projects import mode is valid
		if f.restoreProjectsImportMode != "UPDATE" && f.restoreProjectsImportMode != "INSERT" && f.restoreProjectsImportMode != "" {
			return errors.New("invalid projects-import-mode. Must be either UPDATE or INSERT")
		}

		endpoint := "/fmerest/v3/migration/restore/upload?pauseNotifications=" + strconv.FormatBool(f.restorePauseNotifications)
		if f.restoreImportMode != "" {
			endpoint += "&importMode=" + f.restoreImportMode
		}
		if f.restoreProjectsImportMode != "" {
			endpoint += "&projectsImportMode=" + f.restoreProjectsImportMode
		}

		request, err := buildFmeServerRequest(endpoint, "POST", file)
		if err != nil {
			return err
		}

		request.Header.Set("Content-Type", "application/octet-stream")

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 200 {
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
