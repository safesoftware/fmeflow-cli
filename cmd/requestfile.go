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

var outputLicenseFile string

// requestfileCmd represents the requestfile command
var requestfileCmd = &cobra.Command{
	Use:   "requestfile",
	Short: "Generates a JSON file for requesting a FME Server license file.",
	Long: `Generates a JSON file for requesting a FME Server license file.
	
Example:

# Generate a license request file and output to the console
fmeserver license requestfile --first-name "Billy" --last-name "Bob" --email "billy.bob@example.com" --company "Example Company Inc."

# Generate a license request file and output to a local file
fmeserver license requestfile --first-name "Billy" --last-name "Bob" --email "billy.bob@example.com" --company "Example Company Inc." --file my-request-file.json`,
	Args: NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// add mandatory values
		data := url.Values{
			"firstName": {firstName},
			"lastName":  {lastName},
			"email":     {email},
		}

		// add optional values
		if serialNumber != "" {
			data.Add("serialNumber", serialNumber)
		}
		if company != "" {
			data.Add("company", company)
		}
		if industry != "" {
			data.Add("industry", industry)
		}
		if category != "" {
			data.Add("category", category)
		}
		if salesSource != "" {
			data.Add("salesSource", salesSource)
		}
		if subscribeToUpdates {
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

		if outputLicenseFile != "" {
			tmpfile, err := os.Create(outputLicenseFile)
			if err != nil {
				return err
			}
			defer tmpfile.Close()
			tmpfile.Write(d)
		} else {
			fmt.Println(string(d))
		}
		return nil
	},
}

func init() {
	licenseCmd.AddCommand(requestfileCmd)

	requestfileCmd.Flags().StringVar(&firstName, "first-name", "", "First name to use for license request.")
	requestfileCmd.Flags().StringVar(&lastName, "last-name", "", "Last name to use for license request.")
	requestfileCmd.Flags().StringVar(&email, "email", "", "Email address for license request.")
	requestfileCmd.Flags().StringVar(&serialNumber, "serial-number", "", "Serial Number for the license request.")
	requestfileCmd.Flags().StringVar(&company, "company", "", "Company for the licensing request")
	requestfileCmd.Flags().StringVar(&industry, "industry", "", "Industry for the licensing request")
	requestfileCmd.Flags().StringVar(&category, "category", "", "License Category")
	requestfileCmd.Flags().StringVar(&salesSource, "sales-source", "", "Sales source")
	requestfileCmd.Flags().BoolVar(&subscribeToUpdates, "subscribe-to-updates", false, "Subscribe to Updates")
	requestfileCmd.Flags().StringVar(&outputLicenseFile, "file", "", "Path to file to output to.")
	requestfileCmd.MarkFlagRequired("first-name")
	requestfileCmd.MarkFlagRequired("last-name")
	requestfileCmd.MarkFlagRequired("email")
}
