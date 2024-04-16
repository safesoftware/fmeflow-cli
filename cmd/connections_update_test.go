package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnectionsUpdate(t *testing.T) {
	responseGet := `{
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

	customHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseGet))
			require.NoError(t, err)

		}
		if r.Method == "PUT" {
			w.WriteHeader(http.StatusNoContent)
			_, err := w.Write([]byte(""))
			require.NoError(t, err)

		}

	}

	parameterDoesNotExistBody := `{
		"message": "Unauthorized request by user admin due to lack of proper permissions or the object does not exist."
	  }`
	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"connections", "update", "--name", "PostGIS 3.3 Testsuite", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"connections", "update", "--name", "PostGIS 3.3 Testsuite"},
		},
		{
			name:        "missing name flag",
			wantErrText: "required flag(s) \"name\" not set",
			args:        []string{"connections", "update"},
		},
		{
			name:            "delete parameter",
			statusCode:      http.StatusNoContent,
			args:            []string{"connections", "update", "--name", "PostGIS 3.3 Testsuite"},
			wantOutputRegex: "^Connection successfully updated.[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:        "parameter does not exist",
			statusCode:  http.StatusConflict,
			body:        parameterDoesNotExistBody,
			args:        []string{"connections", "update", "--name", "PostGIS 3.3 Testsuite"},
			wantErrText: "Unauthorized request by user admin due to lack of proper permissions or the object does not exist.",
		},
	}

	runTests(cases, t)

}
