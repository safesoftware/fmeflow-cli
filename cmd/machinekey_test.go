package cmd

import (
	"net/http"
	"testing"
)

func TestMachineKey(t *testing.T) {
	// standard responses for v3
	responseV3 := `{
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
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "machinekey"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "machinekey"},
		},
		{
			name:            "get license machinekey",
			statusCode:      http.StatusOK,
			args:            []string{"license", "machinekey"},
			wantOutputRegex: "^3096247551[\\s]*$",
			body:            responseV3,
		},
		{
			name:           "get license machinekey json",
			statusCode:     http.StatusOK,
			args:           []string{"license", "machinekey", "--json"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
	}

	runTests(cases, t)

}
