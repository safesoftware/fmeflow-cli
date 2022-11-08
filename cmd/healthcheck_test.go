package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestHealthcheck(t *testing.T) {
	// standard responses for v3 and v4
	okResponseV3 := `{
		"status": "ok"
		}`
	okResponseV4 := `{
		"status": "ok",
		"message": "FME Server is healthy."
	  }`
	// random token to use for testing
	testToken := "57463e1b143db046ef3f4ae8ba1b0233e32ee9dd"
	cases := []struct {
		name           string // the name of the test
		statusCode     int    // the http status code the test server should return
		body           string // the body of the request that the test server should return
		wantErrText    string // the expected error text to be returned
		wantOutput     string // regex of the expected output to be returned
		fmeserverBuild int
		args           []string
	}{
		{
			name:        "unknown flag",
			statusCode:  http.StatusOK,
			args:        []string{"--badflag"},
			wantErrText: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
		},
		{
			name:       "v3 health check ok",
			statusCode: http.StatusOK,
			args:       []string{"--api-version", "v3"},
			body:       okResponseV3,
			wantOutput: "ok",
		},
		{
			name:       "v3 health check ready ok",
			statusCode: http.StatusOK,
			args:       []string{"--ready", "--api-version", "v3"},
			body:       okResponseV3,
			wantOutput: "ok",
		},
		{
			name:       "v4 health check ok",
			statusCode: http.StatusOK,
			body:       okResponseV4,
			wantOutput: "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
		},
		{
			name:       "v4 health check ready ok",
			statusCode: http.StatusOK,
			body:       okResponseV4,
			args:       []string{"--ready"},
			wantOutput: "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
		},
		{
			name:           "v3 health check used for 2022.2 build",
			statusCode:     http.StatusOK,
			body:           okResponseV3,
			wantOutput:     "^ok\n$",
			fmeserverBuild: 22765,
		},
		{
			name:           "v4 health check used for 2023.0 build",
			statusCode:     http.StatusOK,
			body:           okResponseV4,
			wantOutput:     "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
			fmeserverBuild: 23200,
		},
		{
			name:       "extra json fields",
			statusCode: http.StatusOK,
			body: `{
				"status": "ok",
				"message": "FME Server is healthy.",
				"extra": "Extra field"
				}`,
			wantOutput:     "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
			fmeserverBuild: 23200,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			cmd := newHealthcheckCmd()

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
				//require.Contains(t, stdOut.String(), tc.wantOutput)
				require.Regexp(t, regexp.MustCompile(tc.wantOutput), stdOut.String())
			}

			//require.ErrorIs(t, err, tc.wantErr)

			/*gotReleases, gotErr := rp.GetAvailableReleases()
			require.ErrorIs(t, gotErr, tc.wantErr)
			if gotErr == nil {
				require.Len(t, gotReleases, tc.wantReleaseCount)
			}*/
		})
	}
}
