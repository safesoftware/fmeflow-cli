package cmd

import (
	"net/http"
	"testing"
)

func TestLicenseStatus(t *testing.T) {
	// standard responses for v3
	responseV3 := `{
		"expiryDate": "PERMANENT",
		"maximumEngines": 10,
		"serialNumber": "AAAA-AAAA-AAAA",
		"isLicenseExpired": false,
		"isLicensed": true,
		"maximumAuthors": 10
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"license", "status", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "status"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "status"},
		},
		{
			name:            "get license status table output",
			statusCode:      http.StatusOK,
			args:            []string{"license", "status"},
			wantOutputRegex: "^[\\s]*EXPIRY DATE[\\s]*MAXIMUM ENGINES[\\s]*SERIAL NUMBER[\\s]*IS LICENSE EXPIRED[\\s]*IS LICENSED[\\s]*MAXIMUM AUTHORS[\\s]*PERMANENT[\\s]*10[\\s]*AAAA-AAAA-AAAA[\\s]*false[\\s]*true[\\s]*10[\\s]*$",
			body:            responseV3,
		},
		{
			name:            "get license status no headers",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"license", "status", "--no-headers"},
			wantOutputRegex: "^[\\s]*PERMANENT[\\s]*10[\\s]*AAAA-AAAA-AAAA[\\s]*false[\\s]*true[\\s]*10[\\s]*$",
		},
		{
			name:           "get license status json",
			statusCode:     http.StatusOK,
			args:           []string{"license", "status", "--json"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
		{
			name:           "get license status json",
			statusCode:     http.StatusOK,
			args:           []string{"license", "status", "--output=json"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
		{
			name:            "get license status custom columns",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"license", "status", "--output", "custom-columns=SERIAL:.serialNumber,LICENSED:.isLicensed"},
			wantOutputRegex: "^[\\s]*SERIAL[\\s]*LICENSED[\\s]*AAAA-AAAA-AAAA[\\s]*true[\\s]*$",
		},
	}

	runTests(cases, t)

}
