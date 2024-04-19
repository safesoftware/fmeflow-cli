package cmd

import (
	"net/http"
	"testing"
)

func TestConnections(t *testing.T) {
	response := `{
		"items": [
		  {
			"name": "Google Drive Named Connection",
			"category": "oauthV2",
			"type": "Google Drive",
			"owner": "admin",
			"shareable": true
		  },
		  {
			"name": "PostGIS 3.3 Testsuite",
			"category": "database",
			"type": "PostgreSQL",
			"owner": "admin",
			"shareable": true,
			"parameters": {
			  "SSL_OPTIONS": "",
			  "PORT": "5434",
			  "CLIENT_PRIV_KEY": "",
			  "SSLMODE": "prefer",
			  "HOST": "somehost",
			  "CA_CERTIFICATE": "",
			  "DATASET": "testsuitegis",
			  "USER_NAME": "testsuite",
			  "CLIENT_CERTIFICATE": ""
			}
		  }
		],
		"totalCount": 2,
		"limit": 100,
		"offset": 0
	  }`

	responseSingle := `{
		"name": "PostGIS 3.3 Testsuite",
		"category": "database",
		"type": "PostgreSQL",
		"owner": "admin",
		"shareable": true,
		"parameters": {
		  "SSL_OPTIONS": "",
		  "PORT": "5434",
		  "CLIENT_PRIV_KEY": "",
		  "SSLMODE": "prefer",
		  "HOST": "somehost",
		  "CA_CERTIFICATE": "",
		  "DATASET": "testsuitegis",
		  "USER_NAME": "testsuite",
		  "CLIENT_CERTIFICATE": ""
		}
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"connections", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"connections"},
		},
		{
			name:            "get connections table output",
			statusCode:      http.StatusOK,
			args:            []string{"connections"},
			body:            response,
			wantOutputRegex: "^[\\s]*NAME[\\s]*TYPE[\\s]*CATEGORY[\\s]*Google Drive Named Connection[\\s]*Google Drive[\\s]*oauthV2[\\s]*PostGIS 3.3 Testsuite[\\s]*PostgreSQL[\\s]*database[\\s]*$",
		},
		{
			name:           "get connections json output",
			statusCode:     http.StatusOK,
			args:           []string{"connections", "--json"},
			body:           response,
			wantOutputJson: response,
		},
		{
			name:            "get single connection",
			statusCode:      http.StatusOK,
			body:            responseSingle,
			args:            []string{"connections", "--name", "PostGIS 3.3 Testsuite"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*TYPE[\\s]*CATEGORY[\\s]*PostGIS 3.3 Testsuite[\\s]*PostgreSQL[\\s]*database[\\s]*$",
		},
		{
			name:            "get connections custom columns",
			statusCode:      http.StatusOK,
			args:            []string{"repositories", "--output=custom-columns=NAME:.name"},
			body:            response,
			wantOutputRegex: "^[\\s]*NAME[\\s]*Google Drive Named Connection[\\s]*PostGIS 3.3 Testsuite[\\s]*$",
		},
	}

	runTests(cases, t)

}
