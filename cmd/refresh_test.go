package cmd

import (
	"net/http"
	"testing"
)

func TestLicenseRefresh(t *testing.T) {
	statusV3 := `{
		"message":"License refresh completed. License file was updated.",
		"status":"SUCCESS"
	}`
	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"license", "refresh", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "refresh"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "refresh"},
		},
		{
			name:            "refresh license",
			statusCode:      http.StatusAccepted,
			args:            []string{"license", "refresh"},
			wantOutputRegex: "^License Refresh Successfully sent\\.[\\s]*$",
			body:            "",
		},
		{
			name:            "refresh license wait",
			statusCode:      http.StatusAccepted,
			args:            []string{"license", "refresh", "--wait"},
			wantOutputRegex: "^License Refresh Successfully sent\\.[\\s]*License refresh completed. License file was updated.[\\s]*$",
			body:            statusV3,
		},
	}

	runTests(cases, t)

}
