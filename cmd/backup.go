package cmd

import (
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
)

type backupFlags struct {
	outputBackupFile    string
	backupResourceName  string
	backupExportPackage string
	backupFailureTopic  string
	backupSuccessTopic  string
	backupResource      bool
	suppressFileRename  bool
}

type BackupResource struct {
	Id int `json:"id"`
}

// backupCmd represents the backup command
func newBackupCmd() *cobra.Command {
	f := backupFlags{}
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backs up the FME Server configuration",
		Long:  "Backs up the FME Server configuration to a local file or to a shared resource location on the FME Server.",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		Example: `
  # back up to a local file
  fmeserver backup -f my_local_backup.fsconfig
	
  # back up to the "Backup" folder in the FME Server Shared Resources with the file name my_fme_backup.fsconfig
  fmeserver backup --resource --export-package my_fme_backup.fsconfig`,
		Args: NoArgs,
		RunE: backupRun(&f),
	}
	cmd.Flags().StringVarP(&f.outputBackupFile, "file", "f", "ServerConfigPackage.fsconfig", "Path to file to download the backup to.")
	cmd.Flags().BoolVar(&f.backupResource, "resource", false, "Backup to a shared resource instead of downloading.")
	cmd.Flags().StringVar(&f.backupResourceName, "resource-name", "FME_SHAREDRESOURCE_BACKUP", "Shared Resource Name where the exported package is saved.")
	cmd.Flags().StringVar(&f.backupExportPackage, "export-package", "ServerConfigPackage.fsconfig", "Path and name of the export package.")
	cmd.Flags().StringVar(&f.backupFailureTopic, "failure-topic", "", "Topic to notify on failure of the backup. Default is MIGRATION_ASYNC_JOB_FAILURE")
	cmd.Flags().StringVar(&f.backupSuccessTopic, "success-topic", "", "Topic to notify on success of the backup. Default is MIGRATION_ASYNC_JOB_SUCCESS")
	cmd.Flags().BoolVar(&f.suppressFileRename, "suppress-file-rename", false, "Specify this flag to not add .fsconfig to the output file automatically")
	cmd.MarkFlagsMutuallyExclusive("file", "resource")
	cmd.MarkFlagsMutuallyExclusive("file", "resource-name")
	cmd.MarkFlagsMutuallyExclusive("file", "export-package")
	cmd.MarkFlagsMutuallyExclusive("file", "failure-topic")
	cmd.MarkFlagsMutuallyExclusive("file", "success-topic")
	cmd.Flags().MarkHidden("suppress-file-rename")
	return cmd
}

func backupRun(f *backupFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// massage the backup file name
		if !f.suppressFileRename && f.outputBackupFile != "" {
			backupExtension := ".fsconfig"
			if !strings.HasSuffix(f.outputBackupFile, backupExtension) {
				f.outputBackupFile += backupExtension
			}
		}

		if !f.backupResource {

			// add mandatory values
			data := url.Values{
				"exportPackageName": {f.outputBackupFile},
			}

			request, err := buildFmeServerRequest("/fmerest/v3/migration/backup/download", "POST", strings.NewReader(data.Encode()))
			if err != nil {
				return err
			}

			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			request.Header.Add("Accept", "application/octet-stream")

			fmt.Fprintln(cmd.OutOrStdout(), "Downloading backup file...")

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != 200 {
				return errors.New(response.Status)
			}
			defer response.Body.Close()

			// Create the output file
			out, err := os.Create(f.outputBackupFile)
			if err != nil {
				return err
			}
			defer out.Close()

			// use Copy so that it doesn't store the entire file in memory
			_, err = io.Copy(out, response.Body)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "FME Server backed up to "+f.outputBackupFile)
		} else {
			// backup to a resource
			// add mandatory values
			data := url.Values{
				"exportPackage": {f.backupExportPackage},
				"resourceName":  {f.backupResourceName},
			}

			// add optional values
			if f.backupSuccessTopic != "" {
				data.Add("successTopic", f.backupSuccessTopic)
			}
			if f.backupFailureTopic != "" {
				data.Add("failureTopic", f.backupFailureTopic)
			}

			request, err := buildFmeServerRequest("/fmerest/v3/migration/backup/resource", "POST", strings.NewReader(data.Encode()))
			if err != nil {
				return err
			}

			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != http.StatusAccepted {
				return errors.New(response.Status)
			}

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result BackupResource
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			} else {
				if !jsonOutput {
					fmt.Fprintln(cmd.OutOrStdout(), "Backup task submitted with id: "+strconv.Itoa(result.Id))
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), string(responseData))
				}
			}
		}
		return nil
	}
}
