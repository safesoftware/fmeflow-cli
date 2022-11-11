package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name           string            // the name of the test
	statusCode     int               // the http status code the test server should return
	body           string            // the body of the request that the test server should return
	wantErrText    string            // the expected error text to be returned
	wantOutput     string            // regex of the expected output to be returned
	fmeserverBuild int               // build to pretend we are contacting
	args           []string          // flags to pass into the command
	wantFormParams map[string]string // array to ensure that all required URL form parameters exist
}

// random token to use for testing
var testToken = "57463e1b143db046ef3f4ae8ba1b0233e32ee9dd"

func runTests(tcs []testCase, fnCmd func() *cobra.Command, t *testing.T) {
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.ParseForm()
				urlParams := r.Form
				for param, value := range tc.wantFormParams {

					require.True(t, urlParams.Has(param))
					require.EqualValues(t, urlParams.Get(param), value)
				}
				//bodyForms := r.ParseForm()

				w.WriteHeader(tc.statusCode)
				_, err := w.Write([]byte(tc.body))
				require.NoError(t, err)
			}))

			defer testServer.Close()
			// set up the config file
			viper.Set("url", testServer.URL)
			viper.Set("token", testToken)
			if tc.fmeserverBuild != 0 {
				viper.Set("build", tc.fmeserverBuild)
			} else {
				viper.Set("build", 23159)
			}

			// create a new copy of the command for each test
			cmd := fnCmd()

			// override the stdout and stderr
			stdOut := bytes.NewBufferString("")
			stdErr := bytes.NewBufferString("")
			cmd.SetOut(stdOut)
			cmd.SetErr(stdErr)

			// set the arguments on the command
			cmd.SetArgs(tc.args)

			// execute
			err := cmd.Execute()

			if err != nil {
				require.EqualValues(t, err.Error(), tc.wantErrText)
			} else {
				require.EqualValues(t, "", tc.wantErrText)
			}
			if tc.wantOutput != "" {

				require.Regexp(t, regexp.MustCompile(tc.wantOutput), stdOut.String())
			}
		})
	}
}
