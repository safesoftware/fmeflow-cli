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

type TokenRequestV4 struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Enabled           bool   `json:"enabled"`
	CustomPermissions bool   `json:"customPermissions"`
	SecondsToExpiry   int    `json:"secondsToExpiry"`
}

type TokenResponseV4 struct {
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Owner             string    `json:"owner"`
	Type              string    `json:"type"`
	CustomPermissions bool      `json:"customPermissions"`
	Enabled           bool      `json:"enabled"`
	Created           time.Time `json:"created"`
	Updated           time.Time `json:"updated"`
	SecondsToExpiry   int       `json:"secondsToExpiry"`
	Expiration        time.Time `json:"expiration"`
	Token             string    `json:"token"`
}

type FMEFlowVersionInfo struct {
	BuildNumber   int    `json:"buildNumber"`
	BuildString   string `json:"buildString"`
	ReleaseYear   int    `json:"releaseYear"`
	MajorVersion  int    `json:"majorVersion"`
	MinorVersion  int    `json:"minorVersion"`
	HotfixVersion int    `json:"hotfixVersion"`
}

type TokenRequestV3 struct {
	Restricted        bool   `json:"restricted"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	ExpirationTimeout int    `json:"expirationTimeout"`
	User              string `json:"user"`
	Enabled           bool   `json:"enabled"`
}

type TokenResponseV3 struct {
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
	apiVersion   apiVersionFlag
}

var urlErrorMsg = "invalid FME Flow URL specified. URL should be of the form https://myfmeflowhostname.com"

var loginV4BuildThreshold = 25208

func newLoginCmd() *cobra.Command {
	f := loginFlags{}
	cmd := &cobra.Command{
		Use:   "login [URL]",
		Short: "Save credentials for an FME Server",
		Long: `Update the config file with the credentials to connect to FME Server. If just a URL is passed in, you will be prompted for a user and password for the FME Server. This will be used to generate an API token that will be saved to the config file for use connecting to FME Server.
	Use the --token flag to pass in an existing API token. To log in with a password on the command line without being prompted, place the password in a text file and pass that in using the --password-file flag.
	This will overwrite any existing credentials saved.`,

		Example: `
  # Prompt for user and password for the given FME Server URL  
  fmeflow login https://my-fmeflow.internal
	
  # Login to an FME Server using a pre-generated token
  fmeflow login https://my-fmeflow.internal --token 5937391ad3a87f19ba14dc6082867373087d031b
	
  # Login to an FME Server using a passed in user and password file (The password is contained in a file at the path /path/to/password-file)
  fmeflow login https://my-fmeflow.internal --user admin --password-file /path/to/password-file`,
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

			// strip an trailing forward slashes from args[0]
			args[0] = strings.TrimRight(args[0], "/")
			// strip the /fmeserver from the end of the URL
			args[0] = strings.TrimSuffix(args[0], "/fmeserver")

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
	cmd.Flags().StringVarP(&f.user, "user", "u", "", "The FME Server user to authenticate and generate an API token for.")
	cmd.Flags().StringVarP(&f.passwordFile, "password-file", "p", "", "A file containing the FME Server password for the user to generate an API token for.")
	cmd.Flags().IntVar(&f.expiration, "expiration", 2592000, "The length of time to generate the token for in seconds.")
	cmd.Flags().Var(&f.apiVersion, "api-version", "The api version to use when contacting FME Server. Must be one of v3 or v4")
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

		// call /fmeinfo/version to retrieve build number
		// we will call /fmeinfo/version in any case (even when the user passes in --api-version)
		// because we need to save the build number to the config file
		if f.token == "" {
			if f.user == "" || f.passwordFile == "" {
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
			viper.Set("url", url)
			// clear any existing token
			viper.Set("token", "")
		} else {
			viper.Set("url", url)
			viper.Set("token", f.token)
		}

		request, err := buildFmeFlowRequest("/fmeinfo/version", "GET", nil)
		if err != nil {
			return err
		}

		if f.token == "" {
			auth := base64.StdEncoding.EncodeToString([]byte(f.user + ":" + password))
			request.Header.Set("Authorization", "Basic "+auth)
		}

		response, err := client.Do(&request)
		if err != nil {
			return err
		} else if response.StatusCode != http.StatusOK {
			return errors.New(response.Status)
		}
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result FMEFlowVersionInfo
		if err := json.Unmarshal(responseData, &result); err != nil {
			return err
		}

		viper.Set("build", result.BuildNumber)

		// if an api version is not specified, determine whether to use v3 or v4 based on the build number
		// otherwise the api version passed in by the user has higher precedence
		if f.apiVersion == "" {
			if viper.GetInt("build") < loginV4BuildThreshold {
				f.apiVersion = apiVersionFlagV3
			} else {
				f.apiVersion = apiVersionFlagV4
			}
		}

		if f.apiVersion == "v4" {
			if f.token == "" {
				currentTime := time.Now()
				tokenRequest := TokenRequestV4{
					Name:              "fmeflow-cli-" + currentTime.Format("20060102150405"),
					Description:       "Token generated for use with the fmeflow-cli.",
					Enabled:           true,
					SecondsToExpiry:   f.expiration,
					CustomPermissions: false,
				}
				tokenJson, err := json.Marshal(tokenRequest)
				if err != nil {
					return err
				}

				auth := base64.StdEncoding.EncodeToString([]byte(f.user + ":" + password))

				req, err := http.NewRequest("POST", url+"/fmeapiv4/tokens", strings.NewReader(string(tokenJson)))
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
					// there was an error logging in. Return the error message
					responseData, err := io.ReadAll(response.Body)
					if err != nil {
						return err
					}
					var responseMessage Message
					if err := json.Unmarshal(responseData, &responseMessage); err != nil {
						// if we fail to unmarshal the response, the body is not json. Return the response status
						return errors.New(response.Status)
					}

					// if there is a message in the response, return that along with the response.Status
					return errors.New(response.Status + ": " + responseMessage.Message)

				}

				responseData, err := io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				var result TokenResponseV4
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

			err := os.MkdirAll(filepath.Dir(viper.ConfigFileUsed()), 0700)
			if err != nil {
				return err
			}
			err = viper.WriteConfig()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Credentials written to "+viper.ConfigFileUsed())

			return nil

		} else if f.apiVersion == "v3" {
			if f.token == "" {
				currentTime := time.Now()
				tokenRequest := TokenRequestV3{
					Restricted:        false,
					Name:              "fmeflow-cli-" + currentTime.Format("20060102150405"),
					Description:       "Token generated for use with the fmeflow-cli.",
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
					// there was an error logging in. Return the error message
					responseData, err := io.ReadAll(response.Body)
					if err != nil {
						return err
					}
					// Unmarshal responseData into a string that is the json text
					var responseMessage Message
					if err := json.Unmarshal(responseData, &responseMessage); err != nil {
						// if we fail to unmarshal the response, the body is not json. Return the response status
						return errors.New(response.Status)
					}

					// if there is a message in the response, return that along with the response.Status
					return errors.New(response.Status + ": " + responseMessage.Message)

				}

				responseData, err := io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				var result TokenResponseV3
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

			// ensure directory where config file is supposed to live exists
			err := os.MkdirAll(filepath.Dir(viper.ConfigFileUsed()), 0700)
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
