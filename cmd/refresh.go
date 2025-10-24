package cmd

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RefreshStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type refreshFlags struct {
	wait       bool
	apiVersion apiVersionFlag
}

var refreshV4BuildThreshold = 23319

func newRefreshCmd() *cobra.Command {
	f := refreshFlags{}
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refreshes the installed license file with a current license from Safe Software.",
		Long:  "Refreshes the installed license file with a current license from Safe Software.",
		Example: `
  # Refresh the license
  fmeflow license refresh`,
		Args: NoArgs,
		RunE: refreshRun(&f),
	}
	cmd.Flags().BoolVar(&f.wait, "wait", false, "Wait for licensing refresh to finish")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	cmd.AddCommand(newRefreshStatusCmd())
	return cmd

}

func refreshRun(f *refreshFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		// get build to decide if we should use v3 or v4
		// FME Server 2023.0+ and later can use v4. Otherwise fall back to v3
		if f.apiVersion == "" {
			fmeflowBuild := viper.GetInt("build")
			if fmeflowBuild < refreshV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		var refreshEndpoint, statusEndpoint string
		if f.apiVersion == "v4" {
			refreshEndpoint = "/fmeapiv4/license/refresh"
			statusEndpoint = "/fmeapiv4/license/refresh/status"
		} else {
			refreshEndpoint = "/fmerest/v3/licensing/refresh"
			statusEndpoint = "/fmerest/v3/licensing/refresh/status"
		}

		request, err := buildFmeFlowRequest(refreshEndpoint, "POST", nil)
		if err != nil {
			return err
		}

		//request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 202 {
			return errors.New(response.Status)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "License Refresh Successfully sent.")

		if f.wait {
			// check the license refresh status until it is finished
			complete := false
			for {
				fmt.Print(".")
				time.Sleep(1 * time.Second)
				// call the status endpoint to see if it is finished
				request, err := buildFmeFlowRequest(statusEndpoint, "GET", nil)
				if err != nil {
					return err
				}
				response, err := client.Do(&request)
				if err != nil {
					return err
				}

				responseData, err := io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				var result RefreshStatus
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				} else {
					// v3 uses uppercase "REQUESTING", v4 uses lowercase "requesting"
					isRequesting := result.Status == "REQUESTING" || result.Status == "requesting"
					if !isRequesting {
						complete = true
						fmt.Fprintln(cmd.OutOrStdout(), result.Message)
					}
				}

				if complete {
					break
				}
			}
		}

		return nil
	}
}
