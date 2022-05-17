/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var outputBackupFile string

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backs up the FME Server configuration",
	Long:  `Backs up the FME Server configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

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
		defer response.Body.Close()
		if err != nil {
			return err
		} else if response.StatusCode != 200 {
			return errors.New(response.Status)
		}

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
		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVarP(&outputBackupFile, "file", "f", "ServerConfigPackage.fsconfig", "Path to file to download the backup to.")

}
