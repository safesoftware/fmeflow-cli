/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var outputBackupFile string
var backupResourceName string
var backupExportPackage string
var backupFailureTopic string
var backupSuccessTopic string
var backupResource bool

type BackupResource struct {
	Id int `json:"id"`
}

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backs up the FME Server configuration",
	Long:  `Backs up the FME Server configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		if !backupResource {

			// add mandatory values
			data := url.Values{
				"exportPackageName": {outputBackupFile},
			}

			request, err := buildFmeServerRequest("/fmerest/v3/migration/backup/download", "POST", strings.NewReader(data.Encode()))
			if err != nil {
				return err
			}

			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			request.Header.Add("Accept", "application/octet-stream")

			fmt.Println("Downloading backup file...")

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != 200 {
				return errors.New(response.Status)
			}
			defer response.Body.Close()

			// Create the output file
			out, err := os.Create(outputBackupFile)
			if err != nil {
				return err
			}
			defer out.Close()

			// use Copy so that it doesn't store the entire file in memory
			_, err = io.Copy(out, response.Body)
			if err != nil {
				return err
			}

			fmt.Println("FME Server backed up to " + outputBackupFile)
		} else {
			// backup to a resource
			// add mandatory values
			data := url.Values{
				"exportPackage": {backupExportPackage},
				"resourceName":  {backupResourceName},
			}

			// add optional values
			if backupSuccessTopic != "" {
				data.Add("successTopic", backupSuccessTopic)
			}
			if backupFailureTopic != "" {
				data.Add("failureTopic", backupFailureTopic)
			}

			request, err := buildFmeServerRequest("/fmerest/v3/migration/backup/resource", "POST", strings.NewReader(data.Encode()))
			if err != nil {
				return err
			}

			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != 202 {
				return errors.New(response.Status)
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result BackupResource
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			} else {
				if !jsonOutput {
					fmt.Println("Backup task submitted with id: " + strconv.Itoa(result.Id))
				} else {
					fmt.Println(string(responseData))
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVarP(&outputBackupFile, "file", "f", "ServerConfigPackage.fsconfig", "Path to file to download the backup to.")
	backupCmd.Flags().BoolVar(&backupResource, "resource", false, "Backup to a shared resource instead of downloading.")
	backupCmd.Flags().StringVar(&backupResourceName, "resource-name", "FME_SHAREDRESOURCE_BACKUP", "Shared Resource Name where the exported package is saved.")
	backupCmd.Flags().StringVar(&backupExportPackage, "export-package", "/ServerConfigPackage.fsconfig", "Path and name of the export package.")
	backupCmd.Flags().StringVar(&backupFailureTopic, "failure-topic", "", "Topic to notify on failure of the backup. Default is MIGRATION_ASYNC_JOB_FAILURE")
	backupCmd.Flags().StringVar(&backupSuccessTopic, "success-topic", "", "Topic to notify on success of the backup. Default is MIGRATION_ASYNC_JOB_SUCCESS")

}
