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

type Healthcheck struct {
	Status string `json:"status"`
}

var token string
var user string
var password string
var expiration string
var timeunit string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save credentials for an FME Server",
	Long: `Log in to an FME Server using a user and password. This will test the credentials to ensure the work correctly and then generate a token
This will overwrite any existing credentials saved.
Example:
  fmeserver login <URL>`,
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

			// test if credentials work
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
	loginCmd.Flags().StringVarP(&expiration, "expiration", "e", "30", "The length of time to generate the token for.")
	loginCmd.Flags().StringVarP(&timeunit, "timeunit", "t", "day", "The timeunit for the expiration of the token.")

	// This isn't quite supported yet. Will work in next release of cobra
	//loginCmd.MarkFlagsRequiredTogether("user", "password")
	//loginCmd.MarkFlagsMutuallyExclusive("token", "user")
	//loginCmd.MarkFlagsMutuallyExclusive("token", "password")
}
