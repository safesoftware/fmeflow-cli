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

type NewConnection struct {
	Category             string                 `json:"category"`
	Name                 string                 `json:"name"`
	Type                 string                 `json:"type"`
	AuthenticationMethod string                 `json:"authenticationMethod,omitempty"`
	Username             string                 `json:"username"`
	Password             string                 `json:"password"`
	Parameters           map[string]interface{} `json:"parameters,omitempty"`
}

type ConnectionCreateFlags struct {
	connectionType       string
	name                 string
	category             string
	authenticationMethod string
	username             string
	password             string
	parameter            []string
}

type ConnectionCreateMessage struct {
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

func newConnectionCreateCmd() *cobra.Command {
	f := ConnectionCreateFlags{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a connection",
		Long:  `Create a connection.`,
		Example: `
  # Create a PostgreSQL connection
  fmeflow connections create --name myPGSQLConnection --category database --type PostgreSQL --parameter HOST=myDBHost --parameter PORT=5432 --parameter DATASET=dbname --parameter USER_NAME=dbuser --parameter SSL_OPTIONS="" --parameter SSLMODE=prefer

  # Create a Google Drive connection (web service must already exist on FME Flow)
  fmeflow connections create --name googleDriveConn --category oauthV2 --type "Google Drive"
`,

		Args: NoArgs,
		RunE: connectionCreateRun(&f),
	}

	cmd.Flags().StringVar(&f.name, "name", "", "Name of the connection to create.")
	cmd.Flags().StringVar(&f.category, "category", "", "Category of the connection to create. Typically it is one of: \"basic\", \"database\", \"token\", \"oauthV1\", \"oauthV2\".")
	cmd.Flags().StringVar(&f.connectionType, "type", "", "Type of connection.")
	cmd.Flags().StringVar(&f.authenticationMethod, "authentication-method", "", "Authentication method of the connection to create.")
	cmd.Flags().StringVar(&f.username, "username", "", "Username of the connection to create.")
	cmd.Flags().StringVar(&f.password, "password", "", "Password of the connection to create.")
	cmd.Flags().StringArrayVar(&f.parameter, "parameter", []string{}, "Parameters of the connection to create. Must be of the form name=value. Can be specified multiple times.")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("category")
	return cmd
}

func connectionCreateRun(f *ConnectionCreateFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		// set up http
		client := &http.Client{}

		var newConnection NewConnection
		newConnection.Name = f.name
		newConnection.Category = f.category
		newConnection.Type = f.connectionType
		if f.authenticationMethod != "" {
			newConnection.AuthenticationMethod = f.authenticationMethod
		}
		newConnection.Username = f.username
		newConnection.Password = f.password

		if len(f.parameter) != 0 {
			newConnection.Parameters = make(map[string]interface{})
		}
		for _, param := range f.parameter {
			parts := strings.Split(param, "=")
			if len(parts) != 2 {
				return errors.New("parameter must be in the format name=value")
			}
			newConnection.Parameters[parts[0]] = parts[1]
		}

		jsonData, err := json.Marshal(newConnection)
		if err != nil {
			return err
		}

		request, err := buildFmeFlowRequest("/fmeapiv4/connections", "POST", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		q := request.URL.Query()
		q.Add("encoded", "false")

		request.Header.Add("Content-Type", "application/json")

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusCreated {
			// attempt to parse the body into JSON as there could be a valuable message in there
			// if fail, just output the status code
			responseData, err := io.ReadAll(response.Body)
			if err == nil {
				var responseMessage ConnectionCreateMessage
				if err := json.Unmarshal(responseData, &responseMessage); err == nil {

					// if json output is requested, output the JSON to stdout before erroring
					if jsonOutput {
						prettyJSON, err := prettyPrintJSON(responseData)
						if err == nil {
							fmt.Fprintln(cmd.OutOrStdout(), prettyJSON)
						} else {
							return errors.New(response.Status)
						}
					} else {
						errorMessage := responseMessage.Message
						for key, value := range responseMessage.Details {
							errorMessage += fmt.Sprintf("\n%s: %v", key, value)
						}
						return errors.New(errorMessage)
					}

				} else {
					return errors.New(response.Status)
				}
			} else {
				return errors.New(response.Status)
			}
		} else {
			if !jsonOutput {
				fmt.Fprintln(cmd.OutOrStdout(), "Connection successfully created.")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "{}")
			}
		}

		return nil
	}
}
