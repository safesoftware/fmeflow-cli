package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type fileContents struct {
	file     string
	contents string
}

type testCase struct {
	name               string            // the name of the test
	statusCode         int               // the http status code the test server should return
	body               string            // the body of the request that the test server should return
	wantErrText        string            // the expected text in the error object to be returned
	wantOutputRegex    string            // regex of the expected stdout to be returned
	wantOutputJson     string            // regex of the expected stdout to be returned
	wantErrOutputRegex string            // regex of the expected stderr to be returned
	wantFormParams     map[string]string // array to ensure that all required URL form parameters exist
	wantFileContents   fileContents      // check file contents
	wantBodyRegEx      string            // check the contents of the body sent
	fmeserverBuild     int               // build to pretend we are contacting
	args               []string          // flags to pass into the command

	httpServer *httptest.Server // custom http test server if needed
}

// random token to use for testing
var testToken = "57463e1b143db046ef3f4ae8ba1b0233e32ee9dd"

func runTests(tcs []testCase, t *testing.T) {
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.httpServer == nil {
				tc.httpServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseForm()
					urlParams := r.Form
					for param, value := range tc.wantFormParams {
						require.True(t, urlParams.Has(param))
						require.EqualValues(t, urlParams.Get(param), value)
					}
					if tc.wantBodyRegEx != "" {
						body, err := io.ReadAll(r.Body)
						require.NoError(t, err)
						require.Regexp(t, regexp.MustCompile(tc.wantBodyRegEx), string(body))
					}
					w.WriteHeader(tc.statusCode)
					_, err := w.Write([]byte(tc.body))
					require.NoError(t, err)
				}))
			}

			defer tc.httpServer.Close()

			// set up the config file
			viper.Set("url", tc.httpServer.URL)
			viper.Set("token", testToken)
			if tc.fmeserverBuild != 0 {
				viper.Set("build", tc.fmeserverBuild)
			} else {
				viper.Set("build", 23159)
			}

			// create a new copy of the command for each test
			cmd := NewRootCommand()

			// override the stdout and stderr
			stdOut := bytes.NewBufferString("")
			stdErr := bytes.NewBufferString("")
			cmd.SetOut(stdOut)
			cmd.SetErr(stdErr)

			// a bit of a hack to make login work
			// requires that URL is passed in first in testing before flags
			if tc.args[0] == "login" {
				tc.args[1] = tc.httpServer.URL
			}

			// if a config file isn't specified, generate a random file and set the config file flag
			if !configFlagExists(tc.args) {
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
