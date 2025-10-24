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

type SystemCode struct {
	SystemCode string `json:"systemCode"`
}

var systemCodeDeprecatedBuildThreshold = 26000

func newSystemCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "systemcode",
		Short: "Retrieves system code of the machine running FME Server.",
		Long:  `Retrieves system code of the machine running FME Server.`,
		Args:  NoArgs,
		RunE:  systemCodeRun(),
	}
	return cmd
}

func systemCodeRun() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Check if systemcode is available in this version of FME Flow
		fmeflowBuild := viper.GetInt("build")
		if fmeflowBuild >= systemCodeDeprecatedBuildThreshold {
			return errors.New("systemcode is not available in this version of FME Flow. The systemcode command was removed in FME Flow 2026.1+")
		}

		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeFlowRequest("/fmerest/v3/licensing/systemcode", "GET", nil)
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

		var result SystemCode
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), result.SystemCode)
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
