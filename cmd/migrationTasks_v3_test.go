package cmd

import (
	"net/http"
	"testing"
)

func TestMigrationTasksV3(t *testing.T) {
	// standard responses for v3
	responseV3 := `{
		"offset": -1,
		"limit": -1,
		"totalCount": 2,
		"items": [
		  {
			"result": "Successful export",
			"id": 2,
			"packageName": "ServerConfigPackage.fsconfig",
			"excludeSensitiveInfo": false,
			"type": "EXPORT",
			"userName": "admin",
			"contentType": "SYSTEM",
			"startDate": "2022-11-09T22:40:51Z",
			"finishedDate": "2022-11-09T22:40:52Z",
			"status": "SUCCESS"
		  },
		  {
			"result": "Successful export",
			"id": 1,
			"packageName": "ServerConfigPackage.fsconfig",
			"excludeSensitiveInfo": false,
			"type": "EXPORT",
			"userName": "admin",
			"contentType": "SYSTEM",
			"startDate": "2022-11-09T21:59:22Z",
			"finishedDate": "2022-11-09T21:59:24Z",
			"status": "SUCCESS"
		  }
		]
	  }`

	migrationLogResponse := `Opening log file using character encoding UTF-8
	  Wed-09-Nov-2022 09:59:22.279 PM   INFORM   Thread-4   409511 : Exporting migration package.
	  Wed-09-Nov-2022 09:59:22.324 PM   INFORM   Thread-4   409574 : Exporting scheduled task Category: Utilities Name: Backup_Configuration.
	  Wed-09-Nov-2022 09:59:22.326 PM   INFORM   Thread-4   409574 : Exporting scheduled task Category: Dashboards Name: DashboardStatisticsGathering.`

	responseV3One := `{
		"result": "Successful export",
		"id": 1,
		"packageName": "ServerConfigPackage.fsconfig",
		"excludeSensitiveInfo": false,
		"type": "EXPORT",
		"userName": "admin",
		"contentType": "SYSTEM",
		"startDate": "2022-11-09T21:59:22Z",
		"finishedDate": "2022-11-09T21:59:24Z",
		"status": "SUCCESS"
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"migration", "tasks", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
			fmeflowBuild:       23166,
		},
		{
			name:         "500 bad status code",
			statusCode:   http.StatusInternalServerError,
			wantErrText:  "500 Internal Server Error",
			args:         []string{"migration", "tasks"},
			fmeflowBuild: 23166,
		},
		{
			name:         "404 bad status code",
			statusCode:   http.StatusNotFound,
			wantErrText:  "404 Not Found",
			args:         []string{"migration", "tasks"},
			fmeflowBuild: 23166,
		},
		{
			name:            "get migration tasks table output",
			statusCode:      http.StatusOK,
			args:            []string{"migration", "tasks"},
			wantOutputRegex: "^[\\s]*ID[\\s]*TYPE[\\s]*USERNAME[\\s]*START TIME[\\s]*END TIME[\\s]*STATUS[\\s]*2[\\s]*EXPORT[\\s]*admin[\\s]*2022-11-09 22:40:51 \\+0000 UTC[\\s]*2022-11-09 22:40:52 \\+0000 UTC[\\s]*SUCCESS[\\s]*1[\\s]*EXPORT[\\s]*admin[\\s]*2022-11-09 21:59:22 \\+0000 UTC[\\s]*2022-11-09 21:59:24 \\+0000 UTC[\\s]*SUCCESS[\\s]*$",
			body:            responseV3,
			fmeflowBuild:    23166,
		},
		{
			name:            "get migration tasks no headers",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"migration", "tasks", "--no-headers"},
			wantOutputRegex: "^[\\s]*2[\\s]*EXPORT[\\s]*admin[\\s]*2022-11-09 22:40:51 \\+0000 UTC[\\s]*2022-11-09 22:40:52 \\+0000 UTC[\\s]*SUCCESS[\\s]*1[\\s]*EXPORT[\\s]*admin[\\s]*2022-11-09 21:59:22 \\+0000 UTC[\\s]*2022-11-09 21:59:24 \\+0000 UTC[\\s]*SUCCESS[\\s]*$",
			fmeflowBuild:    23166,
		},
		{
			name:            "get migration by id",
			statusCode:      http.StatusOK,
			body:            responseV3One,
			args:            []string{"migration", "tasks", "--id", "1"},
			wantOutputRegex: "^[\\s]*ID[\\s]*TYPE[\\s]*USERNAME[\\s]*START TIME[\\s]*END TIME[\\s]*STATUS[\\s]*1[\\s]*EXPORT[\\s]*admin[\\s]*2022-11-09 21:59:22 \\+0000 UTC[\\s]*2022-11-09 21:59:24 \\+0000 UTC[\\s]*SUCCESS[\\s]*$",
			fmeflowBuild:    23166,
		},
		{
			name:            "get migration log",
			statusCode:      http.StatusOK,
			body:            migrationLogResponse,
			args:            []string{"migration", "tasks", "--id", "1", "--log"},
			wantOutputRegex: migrationLogResponse,
			fmeflowBuild:    23166,
		},
		{
			name:           "get migration tasks json",
			statusCode:     http.StatusOK,
			args:           []string{"migration", "tasks", "--json"},
			body:           responseV3,
			wantOutputJson: responseV3,
			fmeflowBuild:   23166,
		},
		{
			name:            "get migration by id custom columns",
			statusCode:      http.StatusOK,
			body:            responseV3One,
			args:            []string{"migration", "tasks", "--id", "1", "--output", "custom-columns=PACKAGE:.packageName,CONTENT:.contentType"},
			wantOutputRegex: "^[\\s]*PACKAGE[\\s]*CONTENT[\\s]*ServerConfigPackage.fsconfig[\\s]*SYSTEM[\\s]*$",
			fmeflowBuild:    23166,
		},
		{
			name:            "get migrations custom columns",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"migration", "tasks", "--output", "custom-columns=PACKAGE:.packageName,CONTENT:.contentType"},
			wantOutputRegex: "^[\\s]*PACKAGE[\\s]*CONTENT[\\s]*ServerConfigPackage.fsconfig[\\s]*SYSTEM[\\s]*ServerConfigPackage.fsconfig[\\s]*SYSTEM[\\s]*$",
			fmeflowBuild:    23166,
		},
	}

	runTests(cases, t)

}
