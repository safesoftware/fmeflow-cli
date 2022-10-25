package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var ready bool

type HealthcheckV3 struct {
	Status string `json:"status"`
}

type HealthcheckV4 struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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
	Args: NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		// This should dynamically be set somehow
		useV4 := true
		endpoint := ""
		if useV4 {
			endpoint = "/fmeapiv4/healthcheck"
			if ready {
				endpoint += "/readiness"
			} else {
				endpoint += "/liveness"
			}
		} else {
			endpoint = "/fmerest/v3/healthcheck"
			if ready {
				endpoint += "?ready=true"
			}
		}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest(endpoint, "GET", nil)
		if err != nil {
			return err
		}
		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != 200 && response.StatusCode != 503 {
			return errors.New(response.Status)
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var resultV3 HealthcheckV3
		var resultV4 HealthcheckV4
		status := ""

		if useV4 {
			if err := json.Unmarshal(responseData, &resultV4); err != nil {
				return err
			}
			status = resultV4.Status
		} else {
			if err := json.Unmarshal(responseData, &resultV3); err != nil {
				return err
			}
			status = resultV3.Status
		}

		if !jsonOutput {
			if useV4 {
				t := createTableWithDefaultColumns(resultV4)

				if noHeaders {
					t.ResetHeaders()
				}
				fmt.Println(t.Render())
			} else {
				fmt.Println(status)
			}
		} else if outputType == "json" {
			prettyJSON, err := prettyPrintJSON(responseData)
			if err != nil {
				return err
			}
			fmt.Println(prettyJSON)
		} else {
			return errors.New("invalid output format specified")
		}
		// if the server is unhealthy, make sure we exit with a non-zero error code
		if useV4 {
			if response.StatusCode == 503 {
				os.Exit(1)
			}
		} else {
			if status != "ok" {
				os.Exit(1)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(healthcheckCmd)
	healthcheckCmd.Flags().BoolVar(&ready, "ready", false, "The health check will report the status of FME Server if it is ready to process jobs.")
}
