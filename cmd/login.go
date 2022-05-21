/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"io/ioutil"
	"net/http"

	"golang.org/x/term"
)

type TokenRequest struct {
	Restricted        bool   `json:"restricted"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	ExpirationTimeout int    `json:"expirationTimeout"`
	User              string `json:"user"`
	Enabled           bool   `json:"enabled"`
}

type TokenResponse struct {
	LastSaveDate   time.Time `json:"lastSaveDate"`
	CreatedDate    time.Time `json:"createdDate"`
	Restricted     bool      `json:"restricted"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Type           string    `json:"type"`
	User           string    `json:"user"`
	Enabled        bool      `json:"enabled"`
	ExpirationDate time.Time `json:"expirationDate"`
	Token          string    `json:"token"`
}

var token string
var user string
var password string
var expiration int
var timeunit string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save credentials for an FME Server",
	Long: `Update the config file with the credentials to connect to FME Server. If just a URL is passed in, you will be prompted for a user and password for the FME Server. This will be used to generate an API token that will be saved to the config file for use connecting to FME Server.
Use the --token flag to pass in an existing API token. It is not recommended to pass the password in on the command line in plaintext.
This will overwrite any existing credentials saved.
Example:
  fmeserver login <URL>
  fmeserver login <URL> --token 5937391ad3a87f19ba14dc6082867373087d031b
  fmeserver login <URL> --user admin --password passw0rd`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a URL")
		}
		if !(strings.HasPrefix(args[0], "http") || strings.HasPrefix(args[0], "https")) {
			return errors.New("url must start with http or https")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		if token == "" {
			if user == "" && password == "" {
				// prompt for a user and password
				reader := bufio.NewReader(os.Stdin)

				fmt.Print("Enter Username: ")
				username, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				fmt.Print("Enter Password: ")
				bytePassword, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return err
				}

				user = strings.TrimSpace(username)
				password = string(bytePassword)

			}

			currentTime := time.Now()

			tokenRequest := TokenRequest{
				Restricted:        false,
				Name:              "fmeserver-cli-" + currentTime.Format("20060102150405"),
				Description:       "Token generated for use with the fmeserver-cli.",
				ExpirationTimeout: expiration,
				User:              user,
				Enabled:           true,
			}

			tokenJson, err := json.Marshal(tokenRequest)
			if err != nil {
				return err
			}

			auth := base64.StdEncoding.EncodeToString([]byte(user + ":" + password))

			req, err := http.NewRequest("POST", url+"/fmerest/v3/tokens", strings.NewReader(string(tokenJson)))
			if err != nil {
				return err
			}
			req.Header.Set("Authorization", "Basic "+auth)
			req.Header.Set("Content-Type", "application/json")

			// create a token to store in config file
			client := &http.Client{}
			response, err := client.Do(req)
			if err != nil {
				return err
			} else if response.StatusCode != 201 {
				return errors.New(response.Status)
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result TokenResponse
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			} else {
				token = result.Token
				fmt.Println("Successfully generated new token.")
			}

		}

		// write to config file
		viper.Set("url", url)
		viper.Set("token", token)
		viper.WriteConfig()
		fmt.Println("Credentials written to " + viper.ConfigFileUsed())

		return nil

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	loginCmd.Flags().StringVarP(&token, "token", "t", "", "The existing API token to use to connect to FME Server")
	loginCmd.Flags().StringVarP(&user, "user", "u", "", "The FME Server user to generate an API token for.")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "The FME Server password for the user to generate an API token for.")
	loginCmd.Flags().IntVar(&expiration, "expiration", 2592000, "The length of time to generate the token for in seconds.")

	// This isn't quite supported yet. Will work in next release of cobra
	//loginCmd.MarkFlagsRequiredTogether("user", "password")
	//loginCmd.MarkFlagsMutuallyExclusive("token", "user")
	//loginCmd.MarkFlagsMutuallyExclusive("token", "password")
}
