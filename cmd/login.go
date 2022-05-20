/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/term"
)

var token string
var user string
var password string
var expiration string
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

			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

			// make sure FME Server is up and ready
			response, err := http.Get(url + "/fmerest/v3/healthcheck?ready=false")
			if err != nil {
				return err
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result Healthcheck
			if err := json.Unmarshal(responseData, &result); err != nil { // Parse []byte to the go struct pointer
				return err
			} else if result.Status != "ok" {
				return errors.New("FME Server is not healthy: " + string(responseData))
			}

			// create a token to store in config file
			response, err = http.Get(url + "/fmetoken/generate?user=" + user + "&password=" + password + "&expiration=" + expiration + "&timeunit=" + timeunit)
			if err != nil {
				return err
			} else if response.StatusCode != 200 {
				return errors.New(response.Status)
			}

			responseData, err = ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}

			token = string(responseData)

		}

		// write to config file
		viper.Set("url", url)
		viper.Set("token", token)
		viper.WriteConfig()

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
	loginCmd.Flags().StringVar(&expiration, "expiration", "30", "The length of time to generate the token for.")
	loginCmd.Flags().StringVar(&timeunit, "timeunit", "day", "The timeunit for the expiration of the token.")

	// This isn't quite supported yet. Will work in next release of cobra
	//loginCmd.MarkFlagsRequiredTogether("user", "password")
	//loginCmd.MarkFlagsMutuallyExclusive("token", "user")
	//loginCmd.MarkFlagsMutuallyExclusive("token", "password")
}
