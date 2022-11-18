package cmd

import (
	"net/http"
	"testing"
)

func TestRefreshStatus(t *testing.T) {
	// standard responses for v3
	statusV3 := `{
		"message":"License refresh completed. License file was updated.",
		"status":"SUCCESS"
	}`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"license", "refresh", "status", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "refresh", "status"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "refresh", "status"},
		},
		{
			name:            "get refresh status",
			statusCode:      http.StatusOK,
			args:            []string{"license", "refresh", "status"},
			wantOutputRegex: "^[\\s]*STATUS[\\s]*MESSAGE[\\s]*SUCCESS[\\s]*License refresh completed. License file was updated.[\\s]*$",
			body:            statusV3,
		},
		{
			name:           "get refresh status json",
			statusCode:     http.StatusOK,
			body:           statusV3,
			args:           []string{"license", "refresh", "status", "--json"},
			wantOutputJson: statusV3,
		},
	}

	runTests(cases, t)

}
