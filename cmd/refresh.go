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
)

type RefreshStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type refreshFlags struct {
	wait bool
}

func newRefreshCmd() *cobra.Command {
	f := refreshFlags{}
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refreshes the installed license file with a current license from Safe Software.",
		Long:  "Refreshes the installed license file with a current license from Safe Software.",
		Example: `
  # Refresh the license
  fmeserver license refresh`,
		Args: NoArgs,
		RunE: refreshRun(&f),
	}
	cmd.Flags().BoolVar(&f.wait, "wait", false, "Wait for licensing refresh to finish")
	cmd.AddCommand(newRefreshStatusCmd())
	return cmd

}

func refreshRun(f *refreshFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		request, err := buildFmeServerRequest("/fmerest/v3/licensing/refresh", "POST", nil)
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

		fmt.Println("License Refresh Successfully sent.")

		if f.wait {
			// check the license refresh status until it is finished
			complete := false
			for {
				fmt.Print(".")
				time.Sleep(1 * time.Second)
				// call the status endpoint to see if it is finished
				request, err := buildFmeServerRequest("/fmerest/v3/licensing/refresh/status", "GET", nil)
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
				} else if result.Status != "REQUESTING" {
					complete = true
					fmt.Println(result.Message)
				}

				if complete {
					break
				}
			}
		}

		return nil
	}
}
