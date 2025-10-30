package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type MachineKey struct {
	MachineKey string `json:"machineKey"`
}

type machineKeyFlags struct {
	apiVersion apiVersionFlag
}

var machineKeyV4BuildThreshold = 23319

func newMachineKeyCmd() *cobra.Command {
	f := machineKeyFlags{}
	cmd := &cobra.Command{
		Use:   "machinekey",
		Short: "Retrieves machine key of the machine running FME Flow.",
		Long:  `Retrieves machine key of the machine running FME Flow.`,
		Args:  NoArgs,
		RunE:  machineKeyRun(&f),
	}

	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
	cmd.Flags().MarkHidden("api-version")
	cmd.RegisterFlagCompletionFunc("api-version", apiVersionFlagCompletion)

	return cmd
}

func machineKeyRun(f *machineKeyFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// get build to decide if we should use v3 or v4
		// FME Server 2023.0+ and later can use v4. Otherwise fall back to v3
		if f.apiVersion == "" {
			fmeflowBuild := viper.GetInt("build")
			if fmeflowBuild < machineKeyV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		// v3 and v4 work exactly the same, so we can just change the endpoint
		var endpoint string
		if f.apiVersion == "v4" {
			endpoint = "/fmeapiv4/license/machinekey"
		} else {
			endpoint = "/fmerest/v3/licensing/machinekey"
		}

		// call the status endpoint to see if it is finished
		request, err := buildFmeFlowRequest(endpoint, "GET", nil)
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

		var result MachineKey
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), result.MachineKey)
			} else {
				prettyJSON, err := prettyPrintJSON(responseData)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
			}
		}
		return nil
	}
}
