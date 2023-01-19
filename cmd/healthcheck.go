package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type healthcheckFlags struct {
	ready      bool
	outputType string
	noHeaders  bool
	apiVersion apiVersionFlag
}

type HealthcheckV3 struct {
	Status string `json:"status"`
}

type HealthcheckV4 struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var healthcheckV4BuildThreshold = 23139

// healthcheckCmd represents the healthcheck command
func newHealthcheckCmd() *cobra.Command {
	f := healthcheckFlags{}
	cmd := &cobra.Command{
		Use:   "healthcheck",
		Short: "Retrieves the health status of FME Server",
		Long:  "Retrieves the health status of FME Server. The health status is normal if the FME Server REST API is responsive. Note that this endpoint does not require authentication. Load balancer or other systems can monitor FME Server using this endpoint without supplying token or password credentials.",
		Example: `
  # Check if the FME Server is healthy and accepting requests
  fmeserver healthcheck
		
  # Check if the FME Server is healthy and ready to process jobs
  fmeserver healthcheck --ready
		
  # Check if the FME Server is healthy and output in json
  fmeserver healthcheck --json
  
  # Check that the FME Server is healthy and output just the status
  fmeserver healthcheck --output=custom-columns=STATUS:.status`,
		Args: NoArgs,
		RunE: healthcheckRun(&f),
	}
	cmd.Flags().BoolVar(&f.ready, "ready", false, "The health check will report the status of FME Server if it is ready to process jobs.")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().StringVarP(&f.outputType, "output", "o", "table", "Specify the output type. Should be one of table, json, or custom-columns")
	cmd.Flags().BoolVar(&f.noHeaders, "no-headers", false, "Don't print column headers")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)
	return cmd
}

func healthcheckRun(f *healthcheckFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// --json overrides --output
		if jsonOutput {
			f.outputType = "json"
		}

		// set up http
		client := &http.Client{}

		// get build to decide if we should use v3 or v4
		// FME Server 2023.0 and later can use v4. Otherwise fall back to v3
		if f.apiVersion == "" {
			fmeserverBuild := viper.GetInt("build")
			if fmeserverBuild < healthcheckV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		endpoint := ""
		if f.apiVersion == "v4" {
			endpoint = "/fmeapiv4/healthcheck"
			if f.ready {
				endpoint += "/readiness"
			} else {
				endpoint += "/liveness"
			}

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

			var resultV4 HealthcheckV4
			if err := json.Unmarshal(responseData, &resultV4); err != nil {
				return err
			}
			if f.outputType == "table" {
				t := createTableWithDefaultColumns(resultV4)

				if noHeaders {
					t.ResetHeaders()
				}
				fmt.Fprintln(cmd.OutOrStdout(), t.Render())

			} else if f.outputType == "json" {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
			} else if strings.HasPrefix(f.outputType, "custom-columns") {
				// parse the columns and json queries
				columnsString := ""
				if strings.HasPrefix(f.outputType, "custom-columns=") {
					columnsString = f.outputType[len("custom-columns="):]
				}
				if len(columnsString) == 0 {
					return errors.New("custom-columns format specified but no custom columns given")
				}

				// we have to marshal the Items array, then create an array of marshalled items
				// to pass to the creation of the table.
				marshalledItems := [][]byte{}
				mJson, err := json.Marshal(resultV4)
				if err != nil {
					return err
				}
				marshalledItems = append(marshalledItems, mJson)

				columnsInput := strings.Split(columnsString, ",")
				t, err := createTableFromCustomColumns(marshalledItems, columnsInput)
				if err != nil {
					return err
				}
				if noHeaders {
					t.ResetHeaders()
				}
				fmt.Fprintln(cmd.OutOrStdout(), t.Render())

			} else {
				return errors.New("invalid output format specified")
			}
			if response.StatusCode == 503 {
				os.Exit(1)
			}
			return nil

		} else if f.apiVersion == "v3" {
			endpoint = "/fmerest/v3/healthcheck"
			if f.ready {
				endpoint += "?ready=true"
			}

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
			status := ""
			var resultV3 HealthcheckV3
			if err := json.Unmarshal(responseData, &resultV3); err != nil {
				return err
			}
			status = resultV3.Status
			if f.outputType == "table" {
				fmt.Fprintln(cmd.OutOrStdout(), status)
			} else if f.outputType == "json" {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
			} else if strings.HasPrefix(f.outputType, "custom-columns") {
				// since V3 only returns a single json parameter, we won't support the custom-columns output type
				return errors.New("custom-columns format not valid with V3 API")
			} else {
				return errors.New("invalid output format specified")
			}
			// if the server is unhealthy, make sure we exit with a non-zero error code
			if status != "ok" {
				os.Exit(1)
			}
			return nil
		} else {
			return fmt.Errorf("invalid apiVersion")
		}
	}
}
