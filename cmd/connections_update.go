package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

type UpdateConnection struct {
	Category             string                 `json:"category"`
	AuthenticationMethod string                 `json:"authenticationMethod,omitempty"`
	Username             string                 `json:"username"`
	Password             string                 `json:"password"`
	Parameters           map[string]interface{} `json:"parameters,omitempty"`
}

type ConnectionUpdateFlags struct {
	name                 string
	authenticationMethod string
	username             string
	password             string
	parameter            []string
}

func newConnectionUpdateCmd() *cobra.Command {
	f := ConnectionUpdateFlags{}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a connection",
		Long:  `Update a connection. Only things that need to be modified need to be specified.`,
		Example: `
  # Update a PostgreSQL connection with the name "myPGSQLConnection" and modify the host to "myDBHost"
  fmeflow connections update --name myPGSQLConnection --parameter HOST=myDBHost
`,

		Args: NoArgs,
		RunE: connectionUpdateRun(&f),
	}

	cmd.Flags().StringVar(&f.name, "name", "", "Name of the connection to update.")
	cmd.Flags().StringVar(&f.authenticationMethod, "authentication-method", "", "Authentication method of the connection to update.")
	cmd.Flags().StringVar(&f.username, "username", "", "Username of the connection to update.")
	cmd.Flags().StringVar(&f.password, "password", "", "Password of the connection to update.")
	cmd.Flags().StringArrayVar(&f.parameter, "parameter", []string{}, "Parameters of the connection to update. Must be of the form name=value. Can be specified multiple times.")

	cmd.MarkFlagRequired("name")
	return cmd
}

func connectionUpdateRun(f *ConnectionUpdateFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		// get the current values of the connection we are going to update
		url := "/fmeapiv4/connections/" + f.name
		request, err := buildFmeFlowRequest(url, "GET", nil)
		if err != nil {
			return err
		}

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusOK {
			responseData, err := io.ReadAll(response.Body)
			if err == nil {
				var responseMessage Message
				if err := json.Unmarshal(responseData, &responseMessage); err == nil {

					// if json output is requested, output the JSON to stdout before erroring
					if jsonOutput {
						prettyJSON, err := prettyPrintJSON(responseData)
						if err == nil {
							fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
						} else {
							return errors.New(response.Status)
						}
					}
					return errors.New(responseMessage.Message)
				} else {
					return errors.New(response.Status)
				}
			} else {
				return errors.New(response.Status)
			}
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var existingConnectionStruct Connection

		if err := json.Unmarshal(responseData, &existingConnectionStruct); err != nil {
			return err
		}

		// build the struct for the update, filling in missing fields with the existing connection's values

		var updateConnectionStruct UpdateConnection

		// category can't be updated
		updateConnectionStruct.Category = existingConnectionStruct.Category

		if f.authenticationMethod != "" {
			updateConnectionStruct.AuthenticationMethod = f.authenticationMethod
		}
		if f.username != "" {
			updateConnectionStruct.Username = f.username
		}
		if f.password != "" {
			updateConnectionStruct.Password = f.password
		}

		// copy over existing parameters
		updateConnectionStruct.Parameters = existingConnectionStruct.Parameters

		// set new parameters
		for _, param := range f.parameter {
			parts := strings.Split(param, "=")
			if len(parts) != 2 {
				return errors.New("parameter must be in the format name=value")
			}
			updateConnectionStruct.Parameters[parts[0]] = parts[1]
		}

		jsonData, err := json.Marshal(updateConnectionStruct)
		if err != nil {
			return err
		}

		request, err = buildFmeFlowRequest(url, "PUT", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", "application/json")

		response, err = client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusNoContent {
			// attempt to parse the body into JSON as there could be a valuable message in there
			// if fail, just output the status code
			responseData, err := io.ReadAll(response.Body)
			if err == nil {
				var responseMessage Message
				if err := json.Unmarshal(responseData, &responseMessage); err == nil {

					// if json output is requested, output the JSON to stdout before erroring
					if jsonOutput {
						prettyJSON, err := prettyPrintJSON(responseData)
						if err == nil {
							fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
						} else {
							return errors.New(response.Status)
						}
					}
					errorMessage := responseMessage.Message
					if response.StatusCode == http.StatusBadRequest {
						errorMessage = errorMessage + ". Check that the specified category is correct."
					}
					return errors.New(errorMessage)
				} else {
					return errors.New(response.Status)
				}
			} else {
				return errors.New(response.Status)
			}
		}

		if !jsonOutput {
			fmt.Fprintln(cmd.OutOrStdout(), "Connection successfully updated.")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "{}")
		}

		return nil
	}
}
