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
	"github.com/spf13/viper"
)

type backupDownloadV4 struct {
	PackageName string `json:"packageName"`
}

type backupResourceV4 struct {
	ResourceName string `json:"resourceName"`
	PackagePath  string `json:"packagePath"`
	SuccessTopic string `json:"successTopic"`
	FailureTopic string `json:"failureTopic"`
}

type backupFlags struct {
	outputBackupFile    string
	backupResourceName  string
	backupExportPackage string
	backupFailureTopic  string
	backupSuccessTopic  string
	backupResource      bool
	suppressFileRename  bool
	apiVersion          apiVersionFlag
}

type BackupResource struct {
	Id int `json:"id"`
}

var backupV4BuildThreshold = 25208

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
  fmeflow backup -f my_local_backup.fsconfig
	
  # back up to the "Backup" folder in the FME Server Shared Resources with the file name my_fme_backup.fsconfig
  fmeflow backup --resource --export-package my_fme_backup.fsconfig`,
		Args: NoArgs,
		RunE: backupRun(&f),
	}
	cmd.Flags().StringVarP(&f.outputBackupFile, "file", "f", "ServerConfigPackage.fsconfig", "Path to file to download the backup to.")
	cmd.Flags().BoolVar(&f.backupResource, "resource", false, "Backup to a shared resource instead of downloading.")
	cmd.Flags().StringVar(&f.backupResourceName, "resource-name", "FME_SHAREDRESOURCE_BACKUP", "Shared Resource Name where the exported package is saved.")
	cmd.Flags().StringVar(&f.backupExportPackage, "export-package", "ServerConfigPackage.fsconfig", "Path and name of the export package when backing up to a shared resource. Must be used with --resource.")
	cmd.Flags().StringVar(&f.backupFailureTopic, "failure-topic", "", "Topic to notify on failure of the backup. In V3, default is MIGRATION_ASYNC_JOB_FAILURE")
	cmd.Flags().StringVar(&f.backupSuccessTopic, "success-topic", "", "Topic to notify on success of the backup. In V3, default is MIGRATION_ASYNC_JOB_SUCCESS")
	cmd.Flags().BoolVar(&f.suppressFileRename, "suppress-file-rename", false, "Specify this flag to not add .fsconfig to the output file automatically")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
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

		if f.apiVersion == "" {
			if viper.GetInt("build") < backupV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		if f.apiVersion == apiVersionFlagV4 {

			if !f.backupResource {
				var backupRequest backupDownloadV4
				backupRequest.PackageName = f.outputBackupFile

				requestBody, err := json.Marshal(backupRequest)
				if err != nil {
					return err
				}

				request, err := buildFmeFlowRequest("/fmeapiv4/migrations/backup/download", "POST", strings.NewReader(string(requestBody)))
				if err != nil {
					return err
				}

				request.Header.Add("Content-Type", "application/json")

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
				var backupRequest backupResourceV4
				backupRequest.PackagePath = f.backupExportPackage
				backupRequest.ResourceName = f.backupResourceName
				if f.backupSuccessTopic != "" {
					backupRequest.SuccessTopic = f.backupSuccessTopic
				}
				if f.backupFailureTopic != "" {
					backupRequest.FailureTopic = f.backupFailureTopic
				}
				requestBody, err := json.Marshal(backupRequest)
				if err != nil {
					return err
				}
				request, err := buildFmeFlowRequest("/fmeapiv4/migrations/backup/resource", "POST", strings.NewReader(string(requestBody)))
				if err != nil {
					return err
				}

				request.Header.Add("Content-Type", "application/json")

				response, err := client.Do(&request)
				if err != nil {
					return err
				} else if response.StatusCode != 202 && response.StatusCode != 200 {
					if response.StatusCode == 401 {
						return errors.New("failed to login")
					} else {
						return errors.New(response.Status)
					}
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

		} else if f.apiVersion == apiVersionFlagV3 {
			if !f.backupResource {

				// add mandatory values
				data := url.Values{
					"exportPackageName": {f.outputBackupFile},
				}

				request, err := buildFmeFlowRequest("/fmerest/v3/migration/backup/download", "POST", strings.NewReader(data.Encode()))
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

				request, err := buildFmeFlowRequest("/fmerest/v3/migration/backup/resource", "POST", strings.NewReader(data.Encode()))
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
		}

		return nil
	}
}
