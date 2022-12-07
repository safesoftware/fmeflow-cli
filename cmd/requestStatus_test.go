package cmd

import (
	"net/http"
	"testing"
)

func TestRequestStatus(t *testing.T) {
	// standard responses for v3
	statusV3 := `{
		"message": "Success! Your FME Server has now been licensed.",
		"status": "SUCCESS"
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"license", "request", "status", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "request", "status"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "request", "status"},
		},
		{
			name:            "get request status",
			statusCode:      http.StatusOK,
			args:            []string{"license", "request", "status"},
			wantOutputRegex: "^[\\s]*STATUS[\\s]*MESSAGE[\\s]*SUCCESS[\\s]*Success! Your FME Server has now been licensed\\.[\\s]*$",
			body:            statusV3,
		},
		{
			name:           "get request status json",
			statusCode:     http.StatusOK,
			body:           statusV3,
			args:           []string{"license", "request", "status", "--json"},
			wantOutputJson: statusV3,
		},
	}

	runTests(cases, t)

}
