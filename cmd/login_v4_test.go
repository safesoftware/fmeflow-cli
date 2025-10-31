package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoginV4(t *testing.T) {
	tokenResponse := `{
        "name": "fmeflow-cli-20221117135041",
		"description": "test token from REST API documentation",
      	"owner": "admin",
		"type": "USER",
      	"customPermissions": false,
		"enabled": true,
      	"created": "2025-09-19T01:22:18.996Z",
      	"updated": "2025-09-19T01:22:18.996Z",
      	"secondsToExpiry": 0,
      	"expiration": "2025-09-19T01:22:18.996Z",
		"token": "5ba5e0fd15c2403bc8b2e3aa1dfb975ca2197fbe"
	  }`

	versionInfoResponse := `{
		"buildNumber": 25300,
		"buildString": "FME Server 2025.1 - Build 25300 - linux-x64",
		"releaseYear": 2025,
		"majorVersion": 1,
		"minorVersion": 0,
		"hotfixVersion": 0
	  }`

	// generate random file for config file
	f, err := os.CreateTemp("", "config-file*.yaml")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up

	// generate random file for password
	passwordFile, err := os.CreateTemp("", "password")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up
	// write out a password to the password file
	passwordFile.Write([]byte("passw0rd"))

	customHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "token") {
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte(tokenResponse))
			require.NoError(t, err)
		} else if strings.Contains(r.URL.Path, "version") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(versionInfoResponse))
			require.NoError(t, err)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	}

	mainHttpServerLogin := httptest.NewServer(http.HandlerFunc(customHttpServerHandler))
	mainHttpServerToken := httptest.NewServer(http.HandlerFunc(customHttpServerHandler))

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			fmeflowBuild:       25300,
			args:               []string{"login", urlPlaceholder, "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:         "500 bad status code",
			statusCode:   http.StatusInternalServerError,
			fmeflowBuild: 25300,
			wantErrText:  "500 Internal Server Error",
			args:         []string{"login", urlPlaceholder, "--user", "admin", "--password-file", passwordFile.Name()},
		},
		{
			name:         "422 bad status code",
			statusCode:   http.StatusNotFound,
			fmeflowBuild: 25300,
			wantErrText:  "404 Not Found",
			args:         []string{"login", urlPlaceholder, "--user", "admin", "--password-file", passwordFile.Name()},
		},
		{
			name:            "login with user and password",
			statusCode:      http.StatusOK,
			args:            []string{"login", mainHttpServerLogin.URL, "--user", "admin", "--password-file", passwordFile.Name(), "--config", f.Name()},
			fmeflowBuild:    25300,
			httpServer:      mainHttpServerLogin,
			wantOutputRegex: "Credentials written to ",
			wantFileContents: fileContents{
				file: f.Name(),
				contents: fmt.Sprintf(`build: 25300
token: 5ba5e0fd15c2403bc8b2e3aa1dfb975ca2197fbe
url: %s
`, mainHttpServerLogin.URL),
			},
			omitConfig: true,
		},
		{
			name:            "login with token",
			statusCode:      http.StatusOK,
			args:            []string{"login", mainHttpServerToken.URL, "--token", "5ba5e0fd15c2403bc8b2e3aa1dfb975ca2197fbf", "--config", f.Name()},
			fmeflowBuild:    25300,
			httpServer:      mainHttpServerToken,
			wantOutputRegex: "Credentials written to ",
			wantFileContents: fileContents{
				file: f.Name(),
				contents: fmt.Sprintf(`build: 25300
token: 5ba5e0fd15c2403bc8b2e3aa1dfb975ca2197fbf
url: %s
`, mainHttpServerToken.URL),
			},
		},
		{
			name:         "missing password flag",
			statusCode:   http.StatusOK,
			fmeflowBuild: 25300,
			args:         []string{"login", urlPlaceholder, "--user", "admin"},
			wantErrText:  "if any flags in the group [user password-file] are set they must all be set; missing [password-file]",
		},
		{
			name:         "missing user flag",
			statusCode:   http.StatusOK,
			fmeflowBuild: 25300,
			args:         []string{"login", urlPlaceholder, "--password-file", passwordFile.Name()},
			wantErrText:  "if any flags in the group [user password-file] are set they must all be set; missing [user]",
		},
		{
			name:         "token and password mutually exclusive",
			statusCode:   http.StatusOK,
			fmeflowBuild: 25300,
			args:         []string{"login", urlPlaceholder, "--user", "admin", "--password-file", passwordFile.Name(), "--token", "5ba5e0fd15c2403bc8b2e3aa1dfb975ca2197fbf"},
			wantErrText:  "if any flags in the group [token password-file] are set none of the others can be; [password-file token] were all set",
		},
	}

	runTests(cases, t)

}
