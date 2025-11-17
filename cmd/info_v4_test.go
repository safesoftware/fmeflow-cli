package cmd

import (
	"net/http"
	"testing"
)

func TestInfoV4(t *testing.T) {
	responseV4 := `{
		"buildNumber":25606,
		"buildString":"FME Flow 2025.1 - Build 25606 - linux-x64",
		"releaseYear":2025,
		"majorVersion":1,
		"minorVersion":0,
		"hotfixVersion":0
	}`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"info", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
			fmeflowBuild:       25606,
		},
		{
			name:         "500 bad status code",
			statusCode:   http.StatusInternalServerError,
			wantErrText:  "500 Internal Server Error",
			args:         []string{"info"},
			fmeflowBuild: 25606,
		},
		{
			name:         "404 bad status code",
			statusCode:   http.StatusNotFound,
			wantErrText:  "404 Not Found",
			args:         []string{"info"},
			fmeflowBuild: 25606,
		},
		{
			name:            "get info table output",
			statusCode:      http.StatusOK,
			body:            responseV4,
			args:            []string{"info"},
			wantOutputRegex: "[\\s]*BUILD NUMBER[\\s]*BUILD STRING[\\s]*RELEASE YEAR[\\s]*MAJOR VERSION[\\s]*MINOR VERSION[\\s]*HOTFIX VERSION[\\s]*25606[\\s]*FME Flow 2025.1 - Build 25606 - linux-x64[\\s]*2025[\\s]*1[\\s]*0[\\s]*0[\\s]*",
			fmeflowBuild:    25606,
		},
		{
			name:            "get info no headers",
			statusCode:      http.StatusOK,
			body:            responseV4,
			args:            []string{"info", "--no-headers"},
			wantOutputRegex: "[\\s]*25606[\\s]*FME Flow 2025.1 - Build 25606 - linux-x64[\\s]*2025[\\s]*1[\\s]*0[\\s]*0[\\s]*",
			fmeflowBuild:    25606,
		},
		{
			name:           "get info json",
			statusCode:     http.StatusOK,
			args:           []string{"info", "--json"},
			body:           responseV4,
			wantOutputJson: responseV4,
			fmeflowBuild:   25606,
		},
		{
			name:           "get info json via output type",
			statusCode:     http.StatusOK,
			args:           []string{"info", "--output=json"},
			body:           responseV4,
			wantOutputJson: responseV4,
			fmeflowBuild:   25606,
		},
		{
			name:            "get info custom columns",
			statusCode:      http.StatusOK,
			body:            responseV4,
			args:            []string{"info", "--output=custom-columns=BUILD NUMBER:.buildNumber,BUILD STRING:.buildString"},
			wantOutputRegex: "[\\s]*BUILD NUMBER[\\s]*BUILD STRING[\\s]*25606[\\s]*FME Flow 2025.1 - Build 25606 - linux-x64[\\s]*",
			fmeflowBuild:    25606,
		},
	}
	runTests(cases, t)
}
