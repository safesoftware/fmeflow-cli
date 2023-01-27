package cmd

import (
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

type projectsDownloadFlags struct {
	file                 string
	name                 string
	excludeSensitiveInfo bool
}

// backupCmd represents the backup command
func newProjectDownloadCmd() *cobra.Command {
	f := projectsDownloadFlags{}
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Downloads an FME Server Project",
		Long:  "Backs up the FME Server configuration to a local file or to a shared resource location on the FME Server.",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		Example: `
  # download a project named "Test Project" to a local file with default name
  fmeserver projects download --name "Test Project"
	
  # download a project named "Test Project" to a local file named MyProject.fsproject
  fmeserver projects download --name "Test Project" -f MyProject.fsproject`,
		Args: NoArgs,
		RunE: projectDownloadRun(&f),
	}
	cmd.Flags().StringVarP(&f.file, "file", "f", "ProjectPackage.fsproject", "Path to file to download the backup to.")
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the project to download.")
	cmd.Flags().BoolVar(&f.excludeSensitiveInfo, "exclude-sensitive-info", false, "Whether to exclude sensitive information from the exported package. Sensitive information will be excluded from connections, subscriptions, publications, schedule tasks, S3 resources, and user accounts. Other items in the project may still contain sensitive data, especially workspaces. Please be careful before sharing the project export pacakge with others.")
	cmd.MarkFlagRequired("name")
	return cmd
}

func projectDownloadRun(f *projectsDownloadFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// add mandatory values
		data := url.Values{
			"exportPackageName":    {f.file},
			"excludeSensitiveInfo": {strconv.FormatBool(f.excludeSensitiveInfo)},
		}

		request, err := buildFmeServerRequest("/fmerest/v3/projects/projects/"+f.name+"/export/download", "POST", strings.NewReader(data.Encode()))
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Accept", "application/octet-stream")

		fmt.Fprintln(cmd.OutOrStdout(), "Downloading project file...")

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

		fmt.Fprintln(cmd.OutOrStdout(), "FME Server backed up to "+f.file)
		return nil
	}
}
