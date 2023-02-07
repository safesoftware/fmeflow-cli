package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectDownload(t *testing.T) {
	// standard responses for v3 and v4
	okResponseV3 := `Random file contents`

	// generate random file to back up to
	f, err := os.CreateTemp("", "*fmeserver-project.fsproject")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"projects", "download", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"projects", "download", "--name", "TestProject"},
		},
		{
			name:        "422 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"projects", "download", "--name", "TestProject"},
		},
		{
			name:        "missing required flags",
			wantErrText: "required flag(s) \"name\" not set",
			args:        []string{"projects", "download"},
		},
		{
			name:             "download to file",
			statusCode:       http.StatusOK,
			args:             []string{"projects", "download", "--name", "TestProject", "--file", f.Name()},
			body:             okResponseV3,
			wantOutputRegex:  "FME Server backed up to",
			wantFileContents: fileContents{file: f.Name(), contents: okResponseV3},
		},
		{
			name:             "download to file exclude sensitive info",
			statusCode:       http.StatusOK,
			args:             []string{"projects", "download", "--name", "TestProject", "--file", f.Name(), "--exclude-sensitive-info"},
			body:             okResponseV3,
			wantOutputRegex:  "FME Server backed up to",
			wantFileContents: fileContents{file: f.Name(), contents: okResponseV3},
			wantFormParams:   map[string]string{"excludeSensitiveInfo": "true"},
		},
	}

	runTests(cases, t)

}
