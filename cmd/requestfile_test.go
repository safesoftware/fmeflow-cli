package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLicenseRequestFile(t *testing.T) {
	// standard responses for v3
	responseV3Status := `{
		"salesSource": "Source",
		"lastName": "Bob",
		"machineKey": "3096247551",
		"fmeBuild": "FME Server 2023.0 - Build 23166 - linux-x64",
		"code": "fc1e6bdd-3ccd-4749-a9aa-7f4ef9039c06",
		"serialNumber": "AAAA-AAAA-AAAA",
		"fmeServerUri": "https://somehost/fmerest/v3/",
		"industry": "Industry",
		"publicIp": "127.0.0.1",
		"subscribeToUpdates": false,
		"productName": "FME Server Standard Edition",
		"firstName": "Billy",
		"requestId": "393d4bf4-81ec-4ffc-9a54-1bdc0676b43b",
		"company": "Example Inc.",
		"category": "Category",
		"email": "billy.bob@example.com"
	 }`

	// generate random file to back up to
	f, err := os.CreateTemp("", "license-request-file")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"license", "requestfile", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "requestfile", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "requestfile", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com"},
		},
		{
			name:        "request license missing email flag",
			statusCode:  http.StatusOK,
			args:        []string{"license", "requestfile", "--first-name", "Billy", "--last-name", "Bob"},
			wantErrText: "required flag(s) \"email\" not set",
		},
		{
			name:        "request license missing last name flag",
			statusCode:  http.StatusOK,
			args:        []string{"license", "requestfile", "--first-name", "Billy", "--email", "billy.bob@example.com"},
			wantErrText: "required flag(s) \"last-name\" not set",
		},
		{
			name:        "request license missing first name flag",
			statusCode:  http.StatusOK,
			args:        []string{"license", "requestfile", "--last-name", "Bob", "--email", "billy.bob@example.com"},
			wantErrText: "required flag(s) \"first-name\" not set",
		},
		{
			name:           "request licensefile",
			statusCode:     http.StatusOK,
			args:           []string{"license", "requestfile", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com"},
			body:           responseV3Status,
			wantOutputJson: responseV3Status,
		},
		{
			name:               "request licensefile into file no filename",
			statusCode:         http.StatusOK,
			args:               []string{"license", "requestfile", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--file"},
			wantErrOutputRegex: "flag needs an argument: --file",
		},
		{
			name:             "request licensefile into file",
			statusCode:       http.StatusOK,
			args:             []string{"license", "requestfile", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--file", f.Name()},
			wantFileContents: fileContents{file: f.Name(), contents: responseV3Status},
			body:             responseV3Status,
		},
		{
			name:           "request license check form params",
			statusCode:     http.StatusOK,
			args:           []string{"license", "requestfile", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--serial-number", "AAAA-AAAA-AAAA", "--company", "Example Inc.", "--industry", "Industry", "--sales-source", "source", "--subscribe-to-updates", "--category", "Category"},
			wantFormParams: map[string]string{"firstName": "Billy", "lastName": "Bob", "email": "billy.bob@example.com", "serialNumber": "AAAA-AAAA-AAAA", "company": "Example Inc.", "category": "Category", "industry": "Industry", "salesSource": "source", "subscribeToUpdates": "true"},
		},
	}

	runTests(cases, t)

}
