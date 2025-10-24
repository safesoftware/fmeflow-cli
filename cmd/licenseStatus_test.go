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

	// standard responses for v4 (different field names)
	responseV4 := `{
		"licensed": true,
		"expiration": "PERMANENT",
		"maximumEngines": 10,
		"expired": false,
		"serialNumber": "AAAA-AAAA-AAAA",
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
			name:        "500 bad status code v3",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "status", "--api-version", "v3"},
		},
		{
			name:        "404 bad status code v3",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "status", "--api-version", "v3"},
		},
		{
			name:            "get license status table output v3",
			statusCode:      http.StatusOK,
			args:            []string{"license", "status", "--api-version", "v3"},
			wantOutputRegex: "^[\\s]*EXPIRY DATE[\\s]*MAXIMUM ENGINES[\\s]*SERIAL NUMBER[\\s]*IS LICENSE EXPIRED[\\s]*IS LICENSED[\\s]*MAXIMUM AUTHORS[\\s]*PERMANENT[\\s]*10[\\s]*AAAA-AAAA-AAAA[\\s]*false[\\s]*true[\\s]*10[\\s]*$",
			body:            responseV3,
		},
		{
			name:            "get license status no headers v3",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"license", "status", "--no-headers", "--api-version", "v3"},
			wantOutputRegex: "^[\\s]*PERMANENT[\\s]*10[\\s]*AAAA-AAAA-AAAA[\\s]*false[\\s]*true[\\s]*10[\\s]*$",
		},
		{
			name:           "get license status json v3",
			statusCode:     http.StatusOK,
			args:           []string{"license", "status", "--json", "--api-version", "v3"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
		{
			name:           "get license status json output flag v3",
			statusCode:     http.StatusOK,
			args:           []string{"license", "status", "--output=json", "--api-version", "v3"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
		{
			name:            "get license status custom columns v3",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"license", "status", "--output", "custom-columns=SERIAL:.serialNumber,LICENSED:.isLicensed", "--api-version", "v3"},
			wantOutputRegex: "^[\\s]*SERIAL[\\s]*LICENSED[\\s]*AAAA-AAAA-AAAA[\\s]*true[\\s]*$",
		},
		{
			name:        "500 bad status code v4",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "status", "--api-version", "v4"},
		},
		{
			name:        "404 bad status code v4",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "status", "--api-version", "v4"},
		},
		{
			name:            "get license status table output v4",
			statusCode:      http.StatusOK,
			args:            []string{"license", "status", "--api-version", "v4"},
			wantOutputRegex: "^[\\s]*LICENSED[\\s]*EXPIRATION[\\s]*MAXIMUM ENGINES[\\s]*EXPIRED[\\s]*SERIAL NUMBER[\\s]*MAXIMUM AUTHORS[\\s]*true[\\s]*PERMANENT[\\s]*10[\\s]*false[\\s]*AAAA-AAAA-AAAA[\\s]*10[\\s]*$",
			body:            responseV4,
		},
		{
			name:            "get license status no headers v4",
			statusCode:      http.StatusOK,
			body:            responseV4,
			args:            []string{"license", "status", "--no-headers", "--api-version", "v4"},
			wantOutputRegex: "^[\\s]*true[\\s]*PERMANENT[\\s]*10[\\s]*false[\\s]*AAAA-AAAA-AAAA[\\s]*10[\\s]*$",
		},
		{
			name:           "get license status json v4",
			statusCode:     http.StatusOK,
			args:           []string{"license", "status", "--json", "--api-version", "v4"},
			body:           responseV4,
			wantOutputJson: responseV4,
		},
		{
			name:           "get license status json output flag v4",
			statusCode:     http.StatusOK,
			args:           []string{"license", "status", "--output=json", "--api-version", "v4"},
			body:           responseV4,
			wantOutputJson: responseV4,
		},
		{
			name:            "get license status custom columns v4",
			statusCode:      http.StatusOK,
			body:            responseV4,
			args:            []string{"license", "status", "--output", "custom-columns=SERIAL:.serialNumber,LICENSED:.licensed", "--api-version", "v4"},
			wantOutputRegex: "^[\\s]*SERIAL[\\s]*LICENSED[\\s]*AAAA-AAAA-AAAA[\\s]*true[\\s]*$",
		},
		{
			name:            "get license status table output (no explicit version)",
			statusCode:      http.StatusOK,
			args:            []string{"license", "status"},
			wantOutputRegex: "^[\\s]*LICENSED[\\s]*EXPIRATION[\\s]*MAXIMUM ENGINES[\\s]*EXPIRED[\\s]*SERIAL NUMBER[\\s]*MAXIMUM AUTHORS[\\s]*true[\\s]*PERMANENT[\\s]*10[\\s]*false[\\s]*AAAA-AAAA-AAAA[\\s]*10[\\s]*$",
			body:            responseV4,
		},
		{
			name:           "get license status json (no explicit version)",
			statusCode:     http.StatusOK,
			args:           []string{"license", "status", "--json"},
			body:           responseV4,
			wantOutputJson: responseV4,
		},
	}

	runTests(cases, t)

}
