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

func TestLogin(t *testing.T) {
	tokenResponse := `{
		"lastSaveDate": "2022-11-17T19:30:44Z",
		"createdDate": "2022-11-17T19:30:44Z",
		"restricted": false,
		"name": "fmeserver-cli-20221117135041",
		"description": "test token from REST API documentation",
		"type": "USER",
		"user": "admin",
		"enabled": true,
		"expirationDate": "2022-11-17T20:30:44Z",
		"token": "5ba5e0fd15c2403bc8b2e3aa1dfb975ca2197fbe"
	  }`

	infoResponse := `{
		"currentTime": "Mon-14-Nov-2022 07:20:24 PM",
		"licenseManagement": true,
		"build": "FME Server 2023.0 - Build 23166 - linux-x64",
		"timeZone": "+0000",
		"version": "FME Server"
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
		} else if strings.Contains(r.URL.Path, "info") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(infoResponse))
			require.NoError(t, err)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	}

	mainHttpServerLogin := httptest.NewServer(http.HandlerFunc(customHttpServerHandler))
	mainHttpServerToken := httptest.NewServer(http.HandlerFunc(customHttpServerHandler))
	//mainHttpServerExpiration := httptest.NewServer(http.HandlerFunc(customHttpServerHandler))

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"login", urlPlaceholder, "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"login", urlPlaceholder, "--user", "admin", "--password-file", passwordFile.Name()},
		},
		{
			name:        "422 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"login", urlPlaceholder, "--user", "admin", "--password-file", passwordFile.Name()},
		},
		{
			name:            "login with user and password",
			statusCode:      http.StatusOK,
			args:            []string{"login", mainHttpServerLogin.URL, "--user", "admin", "--password-file", passwordFile.Name(), "--config", f.Name()},
			httpServer:      mainHttpServerLogin,
			wantOutputRegex: "Credentials written to ",
			wantFileContents: fileContents{
				file: f.Name(),
				contents: fmt.Sprintf(`build: 23166
token: 5ba5e0fd15c2403bc8b2e3aa1dfb975ca2197fbe
url: %s
`, mainHttpServerLogin.URL),
			},
		},
		{
			name:            "login with token",
			statusCode:      http.StatusOK,
			args:            []string{"login", mainHttpServerToken.URL, "--token", "5ba5e0fd15c2403bc8b2e3aa1dfb975ca2197fbf", "--config", f.Name()},
			httpServer:      mainHttpServerToken,
			wantOutputRegex: "Credentials written to ",
			wantFileContents: fileContents{
				file: f.Name(),
				contents: fmt.Sprintf(`build: 23166
token: 5ba5e0fd15c2403bc8b2e3aa1dfb975ca2197fbf
url: %s
`, mainHttpServerToken.URL),
			},
		},
		{
			name:        "missing password flag",
			statusCode:  http.StatusOK,
			args:        []string{"login", urlPlaceholder, "--user", "admin"},
			wantErrText: "if any flags in the group [user password-file] are set they must all be set; missing [password-file]",
		},
		{
			name:        "missing user flag",
			statusCode:  http.StatusOK,
			args:        []string{"login", urlPlaceholder, "--password-file", passwordFile.Name()},
			wantErrText: "if any flags in the group [user password-file] are set they must all be set; missing [user]",
		},
		{
			name:        "token and password mutually exclusive",
			statusCode:  http.StatusOK,
			args:        []string{"login", urlPlaceholder, "--user", "admin", "--password-file", passwordFile.Name(), "--token", "5ba5e0fd15c2403bc8b2e3aa1dfb975ca2197fbf"},
			wantErrText: "if any flags in the group [token password-file] are set none of the others can be; [password-file token] were all set",
		},
	}

	runTests(cases, t)

}
