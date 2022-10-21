package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var ready bool

type Healthcheck struct {
	Status string `json:"status"`
}

// healthcheckCmd represents the healthcheck command
var healthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Retrieves the health status of FME Server",
	Long: `Retrieves the health status of FME Server. The health status is normal if the FME Server REST API is responsive. Note that this endpoint does not require authentication. Load balancer or other systems can monitor FME Server using this endpoint without supplying token or password credentials.
	
Examples:
# Check if the FME Server is healthy and accepting requests
fmeserver healthcheck

# Check if the FME Server is healthy and ready to run jobs
fmeserver healthcheck --ready

# Check if the FME Server is healthy and output in json
fmeserver healthcheck --json`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		endpoint := "/fmerest/v3/healthcheck"
		if ready {
			endpoint += "?ready=true"
		}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest(endpoint, "GET", nil)
		if err != nil {
			return err
		}
		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 200 {
			return errors.New(response.Status)
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result Healthcheck
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Println(result.Status)
			} else if outputType == "json" {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Println(prettyJSON)
			} else {
				return errors.New("invalid output format specified")
			}

		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(healthcheckCmd)
	healthcheckCmd.Flags().BoolVar(&ready, "ready", false, "The health check will report the status of FME Server if it is ready to process jobs.")
}
