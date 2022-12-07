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

type licenseRequestFileFlags struct {
	firstName          string
	lastName           string
	email              string
	serialNumber       string
	company            string
	industry           string
	category           string
	salesSource        string
	subscribeToUpdates bool
	outputLicenseFile  string
}

func newLicenseRequestFileCmd() *cobra.Command {
	f := licenseRequestFileFlags{}
	cmd := &cobra.Command{
		Use:   "requestfile",
		Short: "Generates a JSON file for requesting a FME Server license file.",
		Long: `Generates a JSON file for requesting a FME Server license file.
		
	Example:
	
	# Generate a license request file and output to the console
	fmeserver license requestfile --first-name "Billy" --last-name "Bob" --email "billy.bob@example.com" --company "Example Company Inc."
	
	# Generate a license request file and output to a local file
	fmeserver license requestfile --first-name "Billy" --last-name "Bob" --email "billy.bob@example.com" --company "Example Company Inc." --file my-request-file.json`,
		Args: NoArgs,
		RunE: licenseRequestFileRun(&f),
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
	cmd.Flags().StringVar(&f.outputLicenseFile, "file", "", "Path to file to output to.")
	cmd.MarkFlagRequired("first-name")
	cmd.MarkFlagRequired("last-name")
	cmd.MarkFlagRequired("email")

	return cmd
}

func licenseRequestFileRun(f *licenseRequestFileFlags) func(cmd *cobra.Command, args []string) error {
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

		request, err := buildFmeServerRequest("/fmerest/v3/licensing/requestfile", "POST", strings.NewReader(data.Encode()))
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Accept", "application/octet-stream")

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 200 {
			return errors.New(response.Status)
		}

		// read the body which should be the contents of the file
		d, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if f.outputLicenseFile != "" {
			tmpfile, err := os.Create(f.outputLicenseFile)
			if err != nil {
				return err
			}
			defer tmpfile.Close()
			tmpfile.Write(d)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), string(d))
		}
		return nil
	}
}
