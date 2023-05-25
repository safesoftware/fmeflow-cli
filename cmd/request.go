package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type licenseRequestFlags struct {
	firstName          string
	lastName           string
	email              string
	serialNumber       string
	company            string
	industry           string
	category           string
	salesSource        string
	subscribeToUpdates bool
	wait               bool
}

type RequestStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func newLicenseRequestCmd() *cobra.Command {
	f := licenseRequestFlags{}
	cmd := &cobra.Command{
		Use:   "request",
		Short: "Request a license from the FME Server licensing server",
		Long: `Request a license file from the FME Server licensing server. First name, Last name and email are required for requesting a license file.
  If no serial number is passed in, a trial license will be requested.`,
		Example: `
  # Request a trial license and wait for it to be downloaded and installed
  fmeflow license request --first-name "Billy" --last-name "Bob" --email "billy.bob@example.com" --company "Example Company Inc." --wait
	
  # Request a license with a serial number
  fmeflow license request --first-name "Billy" --last-name "Bob" --email "billy.bob@example.com" --company "Example Company Inc." --serial-number "AAAA-BBBB-CCCC"
	`,
		Args: NoArgs,
		RunE: licenseRequestRun(&f),
	}
	cmd.Flags().StringVar(&f.firstName, "first-name", "", "First name to use for license request.")
	cmd.Flags().StringVar(&f.lastName, "last-name", "", "Last name to use for license request.")
	cmd.Flags().StringVar(&f.email, "email", "", "Email address for license request.")
	cmd.Flags().StringVar(&f.serialNumber, "serial-number", "", "Serial Number for the license request.")
	cmd.Flags().StringVar(&f.company, "company", "", "Company for the licensing request")
	cmd.Flags().StringVar(&f.industry, "industry", "", "Industry for the licensing request")
	cmd.Flags().StringVar(&f.category, "category", "", "License Category")
	cmd.Flags().StringVar(&f.salesSource, "sales-source", "", "Sales source")
	cmd.Flags().BoolVar(&f.subscribeToUpdates, "subscribe-to-updates", false, "Subscribe to Updates")
	cmd.Flags().BoolVar(&f.wait, "wait", false, "Wait for licensing request to finish")
	cmd.MarkFlagRequired("first-name")
	cmd.MarkFlagRequired("last-name")
	cmd.MarkFlagRequired("email")
	cmd.AddCommand(newLicenseRequestStatusCmd())
	return cmd

}

func licenseRequestRun(f *licenseRequestFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// add mandatory values
		data := url.Values{
			"firstName": {f.firstName},
			"lastName":  {f.lastName},
			"email":     {f.email},
		}

		// add optional values
		if f.serialNumber != "" {
			data.Add("serialNumber", f.serialNumber)
		}
		if f.company != "" {
			data.Add("company", f.company)
		}
		if f.industry != "" {
			data.Add("industry", f.industry)
		}
		if f.category != "" {
			data.Add("category", f.category)
		}
		if f.salesSource != "" {
			data.Add("salesSource", f.salesSource)
		}
		if f.subscribeToUpdates {
			data.Add("subscribeToUpdates", "true")
		}

		request, err := buildFmeFlowRequest("/fmerest/v3/licensing/request", "POST", strings.NewReader(data.Encode()))
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

		if !jsonOutput {
			fmt.Fprintln(cmd.OutOrStdout(), "License Request Successfully sent.")
		} else {
			if !f.wait {
				fmt.Fprintln(cmd.OutOrStdout(), "{}")
			}
		}

		if f.wait {
			// check the license status until it is finished
			complete := false
			for {
				if !jsonOutput {
					fmt.Print(".")
				}

				time.Sleep(1 * time.Second)
				// call the status endpoint to see if it is finished
				request, err := buildFmeFlowRequest("/fmerest/v3/licensing/request/status", "GET", nil)
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

				var result RequestStatus
				if err := json.Unmarshal(responseData, &result); err != nil {
					return err
				} else if result.Status != "REQUESTING" {
					complete = true
					if !jsonOutput {
						fmt.Fprintln(cmd.OutOrStdout(), result.Message)
					} else {
						prettyJSON, err := prettyPrintJSON(responseData)
						if err != nil {
							return err
						}
						fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
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
