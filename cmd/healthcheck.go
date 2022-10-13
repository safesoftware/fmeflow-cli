package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	Long:  `Retrieves the health status of FME Server. The health status is normal if the FME Server REST API is responsive. Note that this endpoint does not require authentication. Load balancer or other systems can monitor FME Server using this endpoint without supplying token or password credentials.`,
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

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result Healthcheck
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				// output all values returned by the JSON in a table
				if result.Status == "ok" {
					if ready {
						fmt.Println("FME Server is running and ready to run jobs.")
					} else {
						fmt.Println("FME Server is running and ready to accept requests.")
					}
				} else {
					fmt.Println("FME Server is not healthy.")
					fmt.Println(result.Status)
				}

			} else {
				fmt.Println(string(responseData))
			}

		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(healthcheckCmd)
	healthcheckCmd.Flags().BoolVar(&ready, "ready", false, "The health check will report the status of FME Server if it is ready to process jobs.")
}
