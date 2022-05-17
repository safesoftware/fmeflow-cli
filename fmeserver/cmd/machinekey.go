/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

type MachineKey struct {
	MachineKey string `json:"machineKey"`
}

// machinekeyCmd represents the machinekey command
var machinekeyCmd = &cobra.Command{
	Use:   "machinekey",
	Short: "Retrieves machine key of the machine running FME Server.",
	Long:  `Retrieves machine key of the machine running FME Server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set up http
		client := &http.Client{}

		// call the status endpoint to see if it is finished
		request, err := buildFmeServerRequest("/fmerest/v3/licensing/machinekey", "GET", nil)
		if err != nil {
			return err
		}
		response, err := client.Do(&request)
		if err != nil {
			return err
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result MachineKey
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		} else {
			if !jsonOutput {
				fmt.Printf(result.MachineKey)
			} else {
				fmt.Println(string(responseData))
			}

		}
		return nil
	},
}

func init() {
	licenseCmd.AddCommand(machinekeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// machinekeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// machinekeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
