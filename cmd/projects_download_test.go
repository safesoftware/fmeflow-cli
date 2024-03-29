package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectDownload(t *testing.T) {
	// standard responses for v3 and v4
	okResponse := `Random file contents`

	testProjectJson := `{
		"items": [
		  {
			"id": "a64297e7-a119-4e10-ac37-5d0bba12194b",
			"name": "test",
			"hubUid": "",
			"hubPublisherUid": "",
			"description": "test1",
			"readme": "",
			"version": "1.0.0",
			"lastUpdated": "2024-03-26T18:44:30.713Z",
			"owner": "admin",
			"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"shareable": true,
			"lastUpdateUser": "admin",
			"lastUpdateUserID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"hasIcon": false
		  }
		],
		"totalCount": 1,
		"limit": 100,
		"offset": 0
	  }`

	customHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {

		// send the file if we are downloading
		if strings.Contains(r.URL.Path, "export/download") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(okResponse))
			require.NoError(t, err)
		} else {
			// otherwise we are getting the project by name
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(testProjectJson))
			require.NoError(t, err)
		}

	}

	// generate random file to back up to
	f, err := os.CreateTemp("", "*fmeflow-project.fsproject")
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
			name:             "download to file v4 by id",
			statusCode:       http.StatusOK,
			args:             []string{"projects", "download", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--file", f.Name()},
			body:             okResponse,
			wantOutputRegex:  "Project exported to",
			wantFileContents: fileContents{file: f.Name(), contents: okResponse},
		},
		{
			name:             "download to file exclude sensitive info v4 by id",
			statusCode:       http.StatusOK,
			args:             []string{"projects", "download", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--file", f.Name(), "--exclude-sensitive-info"},
			body:             okResponse,
			wantOutputRegex:  "Project exported to",
			wantFileContents: fileContents{file: f.Name(), contents: okResponse},
			wantBodyRegEx:    ".*\"includeSensitiveInfo\":false.*",
		},
		{
			name:             "download to file exclude selectable items v4 by id",
			statusCode:       http.StatusOK,
			args:             []string{"projects", "download", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--file", f.Name(), "--exclude-selectable-items"},
			body:             okResponse,
			wantOutputRegex:  "Project exported to",
			wantFileContents: fileContents{file: f.Name(), contents: okResponse},
			wantBodyRegEx:    ".*\"excludeAllSelectableItems\":true.*",
		},
		{
			name:             "download to file v4 by name",
			args:             []string{"projects", "download", "--name", "test", "--file", f.Name()},
			wantOutputRegex:  "Project exported to",
			wantFileContents: fileContents{file: f.Name(), contents: okResponse},
			httpServer:       httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:        "422 bad status code V3",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"projects", "download", "--name", "TestProject", "--api-version=v3"},
		},
		{
			name:        "missing required flags",
			wantErrText: "required flag(s) \"name\" not set",
			args:        []string{"projects", "download", "--api-version=v3"},
		},
		{
			name:             "download to file V3",
			statusCode:       http.StatusOK,
			args:             []string{"projects", "download", "--name", "TestProject", "--file", f.Name(), "--api-version=v3"},
			body:             okResponse,
			wantOutputRegex:  "Project exported to",
			wantFileContents: fileContents{file: f.Name(), contents: okResponse},
		},
		{
			name:             "download to file exclude sensitive info V3",
			statusCode:       http.StatusOK,
			args:             []string{"projects", "download", "--name", "TestProject", "--file", f.Name(), "--exclude-sensitive-info", "--api-version=v3"},
			body:             okResponse,
			wantOutputRegex:  "Project exported to",
			wantFileContents: fileContents{file: f.Name(), contents: okResponse},
			wantFormParams:   map[string]string{"excludeSensitiveInfo": "true"},
		},
	}

	runTests(cases, t)

}
