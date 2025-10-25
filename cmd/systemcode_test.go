package cmd

import (
	"net/http"
	"testing"
)

func TestSystemCode(t *testing.T) {
	// standard responses for v3
	responseV3 := `{
		"systemCode": "fc1e6bdd-3ccd-4749-a9aa-7f4ef9039c06"
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"license", "systemcode", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
			fmeflowBuild:       25000, // Use older build that supports systemcode
		},
		{
			name:         "systemcode not available in newer builds",
			statusCode:   http.StatusOK,
			args:         []string{"license", "systemcode"},
			wantErrText:  "systemcode is not available in this version of FME Flow. The systemcode command was removed in FME Flow 2026.1+",
			fmeflowBuild: 26000, // Use build >= 26000 to trigger deprecation error
		},
		{
			name:         "500 bad status code",
			statusCode:   http.StatusInternalServerError,
			wantErrText:  "500 Internal Server Error",
			args:         []string{"license", "systemcode"},
			fmeflowBuild: 25000, // Use older build that supports systemcode
		},
		{
			name:         "404 bad status code",
			statusCode:   http.StatusNotFound,
			wantErrText:  "404 Not Found",
			args:         []string{"license", "systemcode"},
			fmeflowBuild: 25000, // Use older build that supports systemcode
		},
		{
			name:            "get license systemcode",
			statusCode:      http.StatusOK,
			args:            []string{"license", "systemcode"},
			wantOutputRegex: "^fc1e6bdd-3ccd-4749-a9aa-7f4ef9039c06[\\s]*$",
			body:            responseV3,
			fmeflowBuild:    25000, // Use older build that supports systemcode
		},
		{
			name:           "get license systemcode json",
			statusCode:     http.StatusOK,
			args:           []string{"license", "systemcode", "--json"},
			body:           responseV3,
			wantOutputJson: responseV3,
			fmeflowBuild:   25000, // Use older build that supports systemcode
		},
	}

	runTests(cases, t)

}
