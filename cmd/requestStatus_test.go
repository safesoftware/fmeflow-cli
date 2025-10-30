package cmd

import (
	"net/http"
	"testing"
)

func TestRequestStatus(t *testing.T) {
	// standard responses for v3
	status := `{
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
			name:            "get request status v3",
			statusCode:      http.StatusOK,
			args:            []string{"license", "request", "status", "--api-version", "v3"},
			wantOutputRegex: "^[\\s]*STATUS[\\s]*MESSAGE[\\s]*SUCCESS[\\s]*Success! Your FME Server has now been licensed\\.[\\s]*$",
			body:            status,
			fmeflowBuild:    20000, // Force v3
		},
		{
			name:           "get request status v3 json",
			statusCode:     http.StatusOK,
			body:           status,
			args:           []string{"license", "request", "status", "--api-version", "v3", "--json"},
			wantOutputJson: status,
			fmeflowBuild:   20000, // Force v3
		},
		{
			name:            "get request status v4",
			statusCode:      http.StatusOK,
			args:            []string{"license", "request", "status", "--api-version", "v4"},
			wantOutputRegex: "^[\\s]*STATUS[\\s]*MESSAGE[\\s]*SUCCESS[\\s]*Success! Your FME Server has now been licensed\\.[\\s]*$",
			body:            status,
			fmeflowBuild:    25000, // Force v4
		},
		{
			name:           "get request status v4 json",
			statusCode:     http.StatusOK,
			body:           status,
			args:           []string{"license", "request", "status", "--api-version", "v4", "--json"},
			wantOutputJson: status,
			fmeflowBuild:   25000, // Force v4
		},
	}

	runTests(cases, t)

}
