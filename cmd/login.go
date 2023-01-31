package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"net/http"
	"net/url"
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

type loginFlags struct {
	token        string
	user         string
	passwordFile string
	expiration   int
}

func newLoginCmd() *cobra.Command {
	f := loginFlags{}
	cmd := &cobra.Command{
		Use:   "login [URL]",
		Short: "Save credentials for an FME Server",
		Long: `Update the config file with the credentials to connect to FME Server. If just a URL is passed in, you will be prompted for a user and password for the FME Server. This will be used to generate an API token that will be saved to the config file for use connecting to FME Server.
	Use the --token flag to pass in an existing API token. To log in with a password on the command line without being prompted, place the password in a text file and pass that in using the --password-file flag.
	This will overwrite any existing credentials saved.
	
	Examples:
	
	# Prompt for user and password for the given FME Server URL  
	fmeserver login https://my-fmeserver.internal
	
	# Login to an FME Server using a pre-generated token
	fmeserver login https://my-fmeserver.internal --token 5937391ad3a87f19ba14dc6082867373087d031b
	
	# Login to an FME Server using a passed in user and password file
	fmeserver login https://my-fmeserver.internal --user admin --password passw0rd`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Usage()
				return fmt.Errorf("requires a URL")
			}
			if len(args) > 1 {
				cmd.Usage()
				return fmt.Errorf("accepts at most 1 argument, received %d", len(args))
			}
			urlErrorMsg := "invalid FME Server URL specified. URL should be of the form https://myfmeserverhostname.com"
			url, err := url.ParseRequestURI(args[0])
			if err != nil {
				return fmt.Errorf(urlErrorMsg)
			}
			if url.Path != "" {
				return fmt.Errorf(urlErrorMsg)
			}
			return nil
		},
		RunE: loginRun(&f),
	}

	cmd.Flags().StringVarP(&f.token, "token", "t", "", "The existing API token to use to connect to FME Server")
	cmd.Flags().StringVarP(&f.user, "user", "u", "", "The FME Server user to generate an API token for.")
	cmd.Flags().StringVarP(&f.passwordFile, "password-file", "p", "", "A file containing the FME Server password for the user to generate an API token for.")
	cmd.Flags().IntVar(&f.expiration, "expiration", 2592000, "The length of time to generate the token for in seconds.")

	cmd.MarkFlagsRequiredTogether("user", "password-file")
	cmd.MarkFlagsMutuallyExclusive("token", "user")
	cmd.MarkFlagsMutuallyExclusive("token", "password-file")

	return cmd

}
func loginRun(f *loginFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		url := args[0]
		client := &http.Client{}

		var password string

		if f.token == "" {
			if f.user == "" && f.passwordFile == "" {
				// prompt for a user and password
				promptUser := &survey.Input{
					Message: "Username:",
				}
				survey.AskOne(promptUser, &f.user)

				promptPassword := &survey.Password{
					Message: "Password:",
				}
				survey.AskOne(promptPassword, &password)
			} else {
				// get the password from the password-file
				file, err := os.Open(f.passwordFile)
				if err != nil {
					return err
				}
				defer file.Close()

				// just grab the first line of the text file to use as the password
				scanner := bufio.NewScanner(file)
				scanner.Scan()
				password = scanner.Text()

			}

			currentTime := time.Now()

			tokenRequest := TokenRequest{
				Restricted:        false,
				Name:              "fmeserver-cli-" + currentTime.Format("20060102150405"),
				Description:       "Token generated for use with the fmeserver-cli.",
				ExpirationTimeout: f.expiration,
				User:              f.user,
				Enabled:           true,
			}

			tokenJson, err := json.Marshal(tokenRequest)
			if err != nil {
				return err
			}

			auth := base64.StdEncoding.EncodeToString([]byte(f.user + ":" + password))

			req, err := http.NewRequest("POST", url+"/fmerest/v3/tokens", strings.NewReader(string(tokenJson)))
			if err != nil {
				return err
			}
			req.Header.Set("Authorization", "Basic "+auth)
			req.Header.Set("Content-Type", "application/json")

			// create a token to store in config file
			response, err := client.Do(req)
			if err != nil {
				return err
			} else if response.StatusCode != http.StatusCreated {
				return errors.New(response.Status)
			}

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			var result TokenResponse
			if err := json.Unmarshal(responseData, &result); err != nil {
				return err
			} else {
				f.token = result.Token
				fmt.Fprintln(cmd.OutOrStdout(), "Successfully generated new token.")
			}

		}

		// write token and url to config file
		viper.Set("url", url)
		viper.Set("token", f.token)

		request, err := buildFmeServerRequest("/fmerest/v3/info", "GET", nil)
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

		var result FMEServerInfo
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		}

		buildNum, err := parseFMEBuildString(result.Build)
		if err != nil {
			return err
		}
		viper.Set("build", buildNum)
		// ensure directory where config file is supposed to live exists
		err = os.MkdirAll(filepath.Dir(viper.ConfigFileUsed()), 0600)
		if err != nil {
			return err
		}
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Credentials written to "+viper.ConfigFileUsed())

		return nil

	}
}

func parseFMEBuildString(s string) (int, error) {
	first := strings.Split(s, "-")
	if len(first) < 3 {
		return 0, fmt.Errorf("unable to parse build string")
	}

	second := strings.Split(strings.TrimSpace(first[1]), " ")
	if len(second) < 2 {
		return 0, fmt.Errorf("unable to parse build string")
	}
	buildNum, err := strconv.Atoi(second[1])
	if err != nil {
		return 0, fmt.Errorf("unable to parse build string: %w", err)
	}
	return buildNum, nil

}
