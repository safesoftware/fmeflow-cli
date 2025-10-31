package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type fileContents struct {
	file     string
	contents string
}

type testCase struct {
	name               string              // the name of the test
	statusCode         int                 // the http status code the test server should return
	body               string              // the body of the request that the test server should return
	wantErrText        string              // the expected text in the error object to be returned
	wantOutputRegex    string              // regex of the expected stdout to be returned
	wantOutputJson     string              // the expected json to be returned
	wantErrOutputRegex string              // regex of the expected stderr to be returned
	wantFormParams     map[string]string   // array to ensure that all required URL form parameters exist
	wantFormParamsList map[string][]string // for URL forms with multiple values
	wantURLContains    string              // check the URL contains a certain string
	wantFileContents   fileContents        // check file contents
	wantBodyRegEx      string              // check the contents of the body sent
	wantBodyJson       string              // check the JSON body sent
	fmeflowBuild       int                 // build to pretend we are contacting
	args               []string            // flags to pass into the command
	httpServer         *httptest.Server    // custom http test server if needed
	omitConfig         bool                // set this to true if testing a command with no config file set up
	omitConfigToken    bool                // set this to true if testing a command that reads from the config file but doesn't require a token
}

// random token to use for testing
var testToken = "57463e1b143db046ef3f4ae8ba1b0233e32ee9dd"

// string for when we need the test http server in a command
var urlPlaceholder = "[url]"

func runTests(tcs []testCase, t *testing.T) {
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.httpServer == nil {
				tc.httpServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if tc.wantURLContains != "" {
						require.Contains(t, r.URL.String(), tc.wantURLContains)
					}
					r.ParseForm()
					urlParams := r.Form
					for param, value := range tc.wantFormParams {
						require.True(t, urlParams.Has(param))
						require.EqualValues(t, urlParams.Get(param), value)
					}
					for param, value := range tc.wantFormParamsList {
						require.True(t, urlParams.Has(param))
						require.EqualValues(t, urlParams[param], value)
					}
					if tc.wantBodyRegEx != "" {
						body, err := io.ReadAll(r.Body)
						require.NoError(t, err)
						require.Regexp(t, regexp.MustCompile(tc.wantBodyRegEx), string(body))
					}
					if tc.wantBodyJson != "" {
						body, err := io.ReadAll(r.Body)
						require.NoError(t, err)
						require.JSONEq(t, tc.wantBodyJson, string(body))
					}
					w.WriteHeader(tc.statusCode)
					_, err := w.Write([]byte(tc.body))
					require.NoError(t, err)
				}))
			}

			defer tc.httpServer.Close()

			if !tc.omitConfig {
				// set up the config file
				viper.Set("url", tc.httpServer.URL)
				if !tc.omitConfigToken {
					viper.Set("token", testToken)
				}
				if tc.fmeflowBuild != 0 {
					viper.Set("build", tc.fmeflowBuild)
				} else {
					viper.Set("build", 25645)
				}
			}

			// create a new copy of the command for each test
			cmd := NewRootCommand()

			// override the stdout and stderr
			stdOut := bytes.NewBufferString("")
			stdErr := bytes.NewBufferString("")
			cmd.SetOut(stdOut)
			cmd.SetErr(stdErr)

			// a bit of a hack to make login work as it needs the URL of the test server
			for i, s := range tc.args {
				tc.args[i] = strings.Replace(s, urlPlaceholder, tc.httpServer.URL, -1)
			}

			// if a config file isn't specified, generate a random file and set the config file flag
			if !configFlagExists(tc.args) && !tc.omitConfig {
				f, err := os.CreateTemp("", "config-file*.yaml")
				require.NoError(t, err)
				defer os.Remove(f.Name()) // clean up
				// insert right after the command so we don't mess up tests that are testing missing arguments
				tc.args = insert(tc.args, 1, "--config")
				tc.args = insert(tc.args, 2, f.Name())
			}

			// set the arguments on the command
			cmd.SetArgs(tc.args)

			// execute
			err := cmd.Execute()

			if err != nil && err != ErrSilent {
				require.EqualValues(t, tc.wantErrText, err.Error())
			} else {
				require.EqualValues(t, tc.wantErrText, "")
			}
			if tc.wantOutputRegex != "" {

				require.Regexp(t, regexp.MustCompile(tc.wantOutputRegex), stdOut.String())
			}
			if tc.wantErrOutputRegex != "" {
				require.Regexp(t, regexp.MustCompile(tc.wantErrOutputRegex), stdErr.String())
			}
			if tc.wantOutputJson != "" {
				require.JSONEq(t, tc.wantOutputJson, stdOut.String())
			}

			if !isEmpty(tc.wantFileContents) {
				file, err := os.Open(tc.wantFileContents.file)
				require.NoError(t, err)
				buf := new(bytes.Buffer)
				buf.ReadFrom(file)
				contents := buf.String()
				require.EqualValues(t, tc.wantFileContents.contents, contents)
			}

		})
	}
}

// helper function to insert into the middle of a slice
func insert(s []string, index int, item string) []string {
	result := make([]string, len(s)+1)
	copy(result[:index], s[:index])
	result[index] = item
	copy(result[index+1:], s[index:])
	return result
}

// helper function to check if the config flag was already set by the test
func configFlagExists(args []string) bool {
	for _, s := range args {
		if s == "--config" {
			return true
		}
	}
	return false
}
