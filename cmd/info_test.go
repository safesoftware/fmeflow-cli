package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInfo(t *testing.T) {
	// standard responses for v3 and v4
	responseV3 := `{
		"currentTime": "Mon-14-Nov-2022 07:20:24 PM",
		"licenseManagement": true,
		"build": "FME Server 2023.0 - Build 23166 - linux-x64",
		"timeZone": "+0000",
		"version": "FME Server"
	  }`

	// generate random file to back up to
	f, err := os.CreateTemp("", "fmeflow-backup")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"info", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"info"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"info"},
		},
		{
			name:            "get info table output",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"info"},
			wantOutputRegex: "[\\s]*CURRENT TIME[\\s]*LICENSE MANAGEMENT[\\s]*BUILD[\\s]*TIME ZONE[\\s]*VERSION[\\s]*Mon-14-Nov-2022 07:20:24 PM[\\s]*true[\\s]*FME Server 2023.0 - Build 23166 - linux-x64[\\s]*\\+0000[\\s]*FME Server[\\s]*",
		},
		{
			name:            "get info no headers",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"info", "--no-headers"},
			wantOutputRegex: "[\\s]*Mon-14-Nov-2022 07:20:24 PM[\\s]*true[\\s]*FME Server 2023.0 - Build 23166 - linux-x64[\\s]*\\+0000[\\s]*FME Server[\\s]*",
		},
		{
			name:           "get info json",
			statusCode:     http.StatusOK,
			args:           []string{"info", "--json"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
		{
			name:           "get info json via output type",
			statusCode:     http.StatusOK,
			args:           []string{"info", "--output=json"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
		{
			name:            "get info custom columns",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"info", "--output=custom-columns=TIME:.currentTime,BUILD:.build"},
			wantOutputRegex: "[\\s]*TIME[\\s]*BUILD[\\s]*Mon-14-Nov-2022 07:20:24 PM[\\s]*FME Server 2023.0 - Build 23166 - linux-x64[\\s]*",
		},
	}

	runTests(cases, t)

}
