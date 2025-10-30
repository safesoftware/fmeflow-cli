package cmd

import (
	"net/http"
	"testing"
)

func TestRefreshStatus(t *testing.T) {
	// standard responses for v3 (uppercase status)
	statusV3 := `{
		"message":"License refresh completed. License file was updated.",
		"status":"SUCCESS"
	}`

	// standard responses for v4 (lowercase status)
	statusV4 := `{
		"message":"License refresh completed. License file was updated.",
		"status":"success"
	}`

	requestingV3 := `{
		"message":"License refresh in progress.",
		"status":"REQUESTING"
	}`

	requestingV4 := `{
		"message":"License refresh in progress.",
		"status":"requesting"
	}`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"license", "refresh", "status", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code v3",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "refresh", "status", "--api-version", "v3"},
		},
		{
			name:        "404 bad status code v3",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "refresh", "status", "--api-version", "v3"},
		},
		{
			name:            "get refresh status v3",
			statusCode:      http.StatusOK,
			args:            []string{"license", "refresh", "status", "--api-version", "v3"},
			wantOutputRegex: "^[\\s]*STATUS[\\s]*MESSAGE[\\s]*SUCCESS[\\s]*License refresh completed. License file was updated.[\\s]*$",
			body:            statusV3,
		},
		{
			name:           "get refresh status json v3",
			statusCode:     http.StatusOK,
			body:           statusV3,
			args:           []string{"license", "refresh", "status", "--json", "--api-version", "v3"},
			wantOutputJson: statusV3,
		},
		{
			name:            "get refresh requesting status v3",
			statusCode:      http.StatusOK,
			args:            []string{"license", "refresh", "status", "--api-version", "v3"},
			wantOutputRegex: "^[\\s]*STATUS[\\s]*MESSAGE[\\s]*REQUESTING[\\s]*License refresh in progress.[\\s]*$",
			body:            requestingV3,
		},
		{
			name:        "500 bad status code v4",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "refresh", "status", "--api-version", "v4"},
		},
		{
			name:        "404 bad status code v4",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "refresh", "status", "--api-version", "v4"},
		},
		{
			name:            "get refresh status v4",
			statusCode:      http.StatusOK,
			args:            []string{"license", "refresh", "status", "--api-version", "v4"},
			wantOutputRegex: "^[\\s]*STATUS[\\s]*MESSAGE[\\s]*success[\\s]*License refresh completed. License file was updated.[\\s]*$",
			body:            statusV4,
		},
		{
			name:           "get refresh status json v4",
			statusCode:     http.StatusOK,
			body:           statusV4,
			args:           []string{"license", "refresh", "status", "--json", "--api-version", "v4"},
			wantOutputJson: statusV4,
		},
		{
			name:            "get refresh requesting status v4",
			statusCode:      http.StatusOK,
			args:            []string{"license", "refresh", "status", "--api-version", "v4"},
			wantOutputRegex: "^[\\s]*STATUS[\\s]*MESSAGE[\\s]*requesting[\\s]*License refresh in progress.[\\s]*$",
			body:            requestingV4,
		},
		{
			name:            "get refresh status (no explicit version)",
			statusCode:      http.StatusOK,
			args:            []string{"license", "refresh", "status"},
			wantOutputRegex: "^[\\s]*STATUS[\\s]*MESSAGE[\\s]*SUCCESS[\\s]*License refresh completed. License file was updated.[\\s]*$",
			body:            statusV3,
		},
		{
			name:           "get refresh status json (no explicit version)",
			statusCode:     http.StatusOK,
			body:           statusV3,
			args:           []string{"license", "refresh", "status", "--json"},
			wantOutputJson: statusV3,
		},
	}

	runTests(cases, t)

}
