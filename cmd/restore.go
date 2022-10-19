package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var restoreBackupFile string
var restoreImportMode string
var restorePauseNotifications bool
var restoreProjectsImportMode string

type RestoreResource struct {
	Id int `json:"id"`
}

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restores the FME Server configuration from an import package",
	Long:  `Restores the FME Server configuration from an import package`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := &http.Client{}

		file, err := os.Open(restoreBackupFile)
		if err != nil {
			return err
		}
		defer file.Close()

		// verify import mode is valid
		if restoreImportMode != "UPDATE" && restoreImportMode != "INSERT" {
			return errors.New("invalid import-mode. Must be either UPDATE or INSERT")
		}

		// verify projects import mode is valid
		if restoreProjectsImportMode != "UPDATE" && restoreProjectsImportMode != "INSERT" && restoreProjectsImportMode != "" {
			return errors.New("invalid projects-import-mode. Must be either UPDATE or INSERT")
		}

		endpoint := "/fmerest/v3/migration/restore/upload?pauseNotifications=" + strconv.FormatBool(restorePauseNotifications)
		if restoreImportMode != "" {
			endpoint += "&importMode=" + restoreImportMode
		}
		if restoreProjectsImportMode != "" {
			endpoint += "&projectsImportMode=" + restoreProjectsImportMode
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

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result RestoreResource
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Println("Restore task submitted with id: " + strconv.Itoa(result.Id))
			} else {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Println(prettyJSON)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringVarP(&restoreBackupFile, "file", "f", "", "Path to file to download the backup to.")
	restoreCmd.Flags().StringVar(&restoreImportMode, "import-mode", "INSERT", "To import only items in the import package that do not exist on the current instance, specify INSERT. To overwrite items on the current instance with those in the import package, specify UPDATE. Default is INSERT.")
	restoreCmd.Flags().BoolVar(&restorePauseNotifications, "pause-notifications", true, "Disable notifications for the duration of the restore.")
	restoreCmd.Flags().StringVar(&restoreProjectsImportMode, "projects-import-mode", "", "Import mode for projects. To import only projects in the import package that do not exist on the current instance, specify INSERT. To overwrite projects on the current instance with those in the import package, specify UPDATE. If not supplied, importMode will be used.")
}
