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
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "systemcode"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "systemcode"},
		},
		{
			name:            "get license systemcode",
			statusCode:      http.StatusOK,
			args:            []string{"license", "systemcode"},
			wantOutputRegex: "^fc1e6bdd-3ccd-4749-a9aa-7f4ef9039c06[\\s]*$",
			body:            responseV3,
		},
		{
			name:           "get license systemcode json",
			statusCode:     http.StatusOK,
			args:           []string{"license", "systemcode", "--json"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
	}

	runTests(cases, t)

}
