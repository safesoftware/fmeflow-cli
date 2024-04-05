package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnectionsDelete(t *testing.T) {
	customHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(""))
			require.NoError(t, err)

		}
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
			_, err := w.Write([]byte(""))
			require.NoError(t, err)

		}

	}

	paramMissingBody := `{
		"message": "Unauthorized request by user admin due to lack of proper permissions or the object does not exist."
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"connections", "delete", "--name", "myConn", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"connections", "delete", "--name", "myConn", "--no-prompt"},
		},
		{
			name:        "missing flag",
			wantErrText: "required flag(s) \"name\" not set",
			args:        []string{"connections", "delete"},
		},
		{
			name:            "delete parameter",
			statusCode:      http.StatusNoContent,
			args:            []string{"connections", "delete", "--name", "myConn", "--no-prompt"},
			wantOutputRegex: "^Connection successfully deleted.[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:        "delete parameter not found",
			statusCode:  http.StatusNotFound,
			body:        paramMissingBody,
			args:        []string{"connections", "delete", "--name", "myConn", "--no-prompt"},
			wantErrText: "Unauthorized request by user admin due to lack of proper permissions or the object does not exist.",
		},
	}

	runTests(cases, t)

}
