package cmd

import (
	"net/http"
	"testing"
)

func TestMachineKey(t *testing.T) {
	// standard response for both v3 and v4 (both use the same JSON format)
	response := `{
		"machineKey": "3096247551"
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"license", "machinekey", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code v3",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "machinekey", "--api-version", "v3"},
		},
		{
			name:        "404 bad status code v3",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "machinekey", "--api-version", "v3"},
		},
		{
			name:            "get license machinekey v3",
			statusCode:      http.StatusOK,
			args:            []string{"license", "machinekey", "--api-version", "v3"},
			wantOutputRegex: "^3096247551[\\s]*$",
			body:            response,
		},
		{
			name:           "get license machinekey json v3",
			statusCode:     http.StatusOK,
			args:           []string{"license", "machinekey", "--json", "--api-version", "v3"},
			body:           response,
			wantOutputJson: response,
		},
		{
			name:        "500 bad status code v4",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "machinekey", "--api-version", "v4"},
		},
		{
			name:        "404 bad status code v4",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "machinekey", "--api-version", "v4"},
		},
		{
			name:            "get license machinekey v4",
			statusCode:      http.StatusOK,
			args:            []string{"license", "machinekey", "--api-version", "v4"},
			wantOutputRegex: "^3096247551[\\s]*$",
			body:            response,
		},
		{
			name:           "get license machinekey json v4",
			statusCode:     http.StatusOK,
			args:           []string{"license", "machinekey", "--json", "--api-version", "v4"},
			body:           response,
			wantOutputJson: response,
		},
		{
			name:            "get license machinekey (no explicit version)",
			statusCode:      http.StatusOK,
			args:            []string{"license", "machinekey"},
			wantOutputRegex: "^3096247551[\\s]*$",
			body:            response,
		},
		{
			name:           "get license machinekey json (no explicit version)",
			statusCode:     http.StatusOK,
			args:           []string{"license", "machinekey", "--json"},
			body:           response,
			wantOutputJson: response,
		},
	}

	runTests(cases, t)

}
