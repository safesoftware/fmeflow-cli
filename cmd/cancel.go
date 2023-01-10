/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Message struct {
	Message string `json:"message"`
}

type cancelFlags struct {
	id         string
	apiVersion apiVersionFlag
}

var cancelV4BuildThreshold = 22337

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
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	cmd.MarkFlagRequired("id")

	return cmd
}

func runCancel(f *cancelFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// get build to decide if we should use v3 or v4
		// FME Server 2022.0 and later can use v4. Otherwise fall back to v3
		if f.apiVersion == "" {
			fmeserverBuild := viper.GetInt("build")
			if fmeserverBuild < healthcheckV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		if f.apiVersion == "v4" {
			endpoint := "/fmeapiv4/jobs/" + f.id + "/cancel"

			request, err := buildFmeServerRequest(endpoint, "POST", nil)
			if err != nil {
				return err
			}
			response, err := client.Do(&request)
			if err != nil {
				return err
			} else if response.StatusCode != 204 {
				// attempt to parse the body into JSON as there could be a valuable message in there
				// if fail, just output the status code
				responseData, err := io.ReadAll(response.Body)
				if err == nil {

					var responseMessage Message
					if err := json.Unmarshal(responseData, &responseMessage); err == nil {

						// if json output is requested, output the JSON to stdout before erroring
						if jsonOutput {
							prettyJSON, err := prettyPrintJSON(responseData)
							if err == nil {
								fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
							} else {
								return errors.New(response.Status)
							}
						}
						return errors.New(responseMessage.Message)
					} else {
						return errors.New(response.Status)
					}
				} else {
					return errors.New(response.Status)
				}

			}

			if jsonOutput {
				// This endpoint returns no content if successful. Just output empty JSON if requested.
				fmt.Fprintln(cmd.OutOrStdout(), "{}")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "Success. The job with id "+f.id+" was cancelled.")
			}

			return nil

		} else if f.apiVersion == "v3" {

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
		return nil
	}
}
