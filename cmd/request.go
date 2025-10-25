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
	"github.com/spf13/viper"
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
	apiVersion         apiVersionFlag
}

// v3 response structure
type RequestStatusV3 struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// v4 request structure
type LicenseRequestV4 struct {
	JobTitle           string `json:"jobTitle,omitempty"`
	Company            string `json:"company,omitempty"`
	Industry           string `json:"industry,omitempty"`
	Email              string `json:"email"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	SerialNumber       string `json:"serialNumber,omitempty"`
	SubscribeToUpdates bool   `json:"subscribeToUpdates"`
}

// v4 response structure
type RequestStatusV4 struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var licenseRequestV4BuildThreshold = 23319

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
	cmd.Flags().StringVar(&f.category, "category", "", "License Category (v3 API only)")
	cmd.Flags().StringVar(&f.salesSource, "sales-source", "", "Sales source (v3 API only)")
	cmd.Flags().BoolVar(&f.subscribeToUpdates, "subscribe-to-updates", false, "Subscribe to Updates")
	cmd.Flags().BoolVar(&f.wait, "wait", false, "Wait for licensing request to finish")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	cmd.MarkFlagRequired("first-name")
	cmd.MarkFlagRequired("last-name")
	cmd.MarkFlagRequired("email")
	cmd.AddCommand(newLicenseRequestStatusCmd())
	return cmd

}

func licenseRequestRun(f *licenseRequestFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// get build to decide if we should use v3 or v4
		// FME Server 2023.0+ and later can use v4. Otherwise fall back to v3
		if f.apiVersion == "" {
			fmeflowBuild := viper.GetInt("build")
			if fmeflowBuild < licenseRequestV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		// Validate v3-only flags are not used with v4
		if f.apiVersion == "v4" {
			if f.category != "" {
				return fmt.Errorf("the --category flag is only supported with v3 API. Use --api-version v3 or remove the --category flag")
			}
			if f.salesSource != "" {
				return fmt.Errorf("the --sales-source flag is only supported with v3 API. Use --api-version v3 or remove the --sales-source flag")
			}
		}

		if f.apiVersion == "v4" {
			return licenseRequestRunV4(f, cmd)
		} else {
			return licenseRequestRunV3(f, cmd)
		}
	}
}

func licenseRequestRunV3(f *licenseRequestFlags, cmd *cobra.Command) error {
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

			var result RequestStatusV3
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

func licenseRequestRunV4(f *licenseRequestFlags, cmd *cobra.Command) error {
	// set up http
	client := &http.Client{}

	// create request body
	requestBody := LicenseRequestV4{
		FirstName:          f.firstName,
		LastName:           f.lastName,
		Email:              f.email,
		SubscribeToUpdates: f.subscribeToUpdates,
	}

	// add optional values
	if f.serialNumber != "" {
		requestBody.SerialNumber = f.serialNumber
	}
	if f.company != "" {
		requestBody.Company = f.company
	}
	if f.industry != "" {
		requestBody.Industry = f.industry
	}

	// marshal JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	request, err := buildFmeFlowRequest("/fmeapiv4/license/request", "POST", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")

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
			request, err := buildFmeFlowRequest("/fmeapiv4/license/request/status", "GET", nil)
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

			var result RequestStatusV4
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			} else if result.Status != "requesting" {
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
