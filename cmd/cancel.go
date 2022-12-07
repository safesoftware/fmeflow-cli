/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

type cancelFlags struct {
	id string
}

func newCancelCmd() *cobra.Command {
	f := cancelFlags{}
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a running job on FME Server",
		Long:  `Cancels the job and marks it as aborted in the completed jobs section, but does not remove it from the database.`,
		Example: `
  # Cancel a job with id 42
  fmeserver cancel --id 42
	`,
		Args: NoArgs,
		RunE: runCancel(&f),
	}

	cmd.Flags().StringVar(&f.id, "id", "", "	The ID of the job to cancel.")
	cmd.MarkFlagRequired("id")

	return cmd
}

func runCancel(f *cancelFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/transformations/jobs/running/"+f.id, "DELETE", nil)
		if err != nil {
			return err
		}
		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode == 404 {
			return errors.New("the specified job ID was not found")
		} else if response.StatusCode != 204 {
			return errors.New(response.Status)
		}

		if jsonOutput {
			fmt.Fprintln(cmd.OutOrStdout(), "{}")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "Success. The job with id "+f.id+" was cancelled.")
		}

		return nil
	}
}
