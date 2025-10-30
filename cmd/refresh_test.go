package cmd

import (
	"net/http"
	"testing"
)

func TestLicenseRefresh(t *testing.T) {
	// v3 responses (uppercase status)
	statusV3 := `{
		"message":"License refresh completed. License file was updated.",
		"status":"SUCCESS"
	}`

	// v4 responses (lowercase status)
	statusV4 := `{
		"message":"License refresh completed. License file was updated.",
		"status":"success"
	}`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"license", "refresh", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code v3",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "refresh", "--api-version", "v3"},
		},
		{
			name:        "404 bad status code v3",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "refresh", "--api-version", "v3"},
		},
		{
			name:            "refresh license v3",
			statusCode:      http.StatusAccepted,
			args:            []string{"license", "refresh", "--api-version", "v3"},
			wantOutputRegex: "^License Refresh Successfully sent\\.[\\s]*$",
			body:            "",
		},
		{
			name:            "refresh license wait v3",
			statusCode:      http.StatusAccepted,
			args:            []string{"license", "refresh", "--wait", "--api-version", "v3"},
			wantOutputRegex: "^License Refresh Successfully sent\\.[\\s]*License refresh completed. License file was updated.[\\s]*$",
			body:            statusV3,
		},
		{
			name:        "500 bad status code v4",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "refresh", "--api-version", "v4"},
		},
		{
			name:        "404 bad status code v4",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "refresh", "--api-version", "v4"},
		},
		{
			name:            "refresh license v4",
			statusCode:      http.StatusAccepted,
			args:            []string{"license", "refresh", "--api-version", "v4"},
			wantOutputRegex: "^License Refresh Successfully sent\\.[\\s]*$",
			body:            "",
		},
		{
			name:            "refresh license wait v4",
			statusCode:      http.StatusAccepted,
			args:            []string{"license", "refresh", "--wait", "--api-version", "v4"},
			wantOutputRegex: "^License Refresh Successfully sent\\.[\\s]*License refresh completed. License file was updated.[\\s]*$",
			body:            statusV4,
		},
		{
			name:            "refresh license (no explicit version)",
			statusCode:      http.StatusAccepted,
			args:            []string{"license", "refresh"},
			wantOutputRegex: "^License Refresh Successfully sent\\.[\\s]*$",
			body:            "",
		},
		{
			name:            "refresh license wait (no explicit version)",
			statusCode:      http.StatusAccepted,
			args:            []string{"license", "refresh", "--wait"},
			wantOutputRegex: "^License Refresh Successfully sent\\.[\\s]*License refresh completed. License file was updated.[\\s]*$",
			body:            statusV4,
		},
	}

	runTests(cases, t)

}
