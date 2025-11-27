package cmd

import (
	"net/http"
	"testing"
)

func TestMigrationTasksV4(t *testing.T) {
	responseV4 := `{
		"offset": -1,
		"limit": -1,
		"totalCount": 2,
		"items": [
		  {
			"result": "Successful export",
			"id": 2,
			"packageName": "ServerConfigPackage.fsconfig",
			"type": "EXPORT",
			"username": "admin",
			"startDate": "2022-11-09T22:40:51Z",
			"finishedDate": "2022-11-09T22:40:52Z",
			"status": "SUCCESS"
		  },
		  {
			"result": "Successful export",
			"id": 1,
			"packageName": "ServerConfigPackage.fsconfig",
			"type": "EXPORT",
			"username": "admin",
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

	migrationLogParsedResponse := `{"items":[{"number":1,"time":"2025-11-19T20:24:33.836Z","status":"INFORM","message":"Getting migration package from shared resource \"FME_SHAREDRESOURCE_SYSTEM\" with path \"/temp/tomcat/5c2c13eac945440fb4c71ed1d1c51587/localhost_2025-11-19-T104055_b25606.fsconfig\"...","line":"Wed-19-Nov-2025 08:24:33.836 PM   INFORM   Thread-8   409520 : Getting migration package from shared resource \"FME_SHAREDRESOURCE_SYSTEM\" with path \"/temp/tomcat/5c2c13eac945440fb4c71ed1d1c51587/localhost_2025-11-19-T104055_b25606.fsconfig\"..."},{"number":2,"time":"2025-11-19T20:24:33.844Z","status":"INFORM","message":"Unzipping migration package.","line":"Wed-19-Nov-2025 08:24:33.844 PM   INFORM   Thread-8   409503 : Unzipping migration package."},{"number":3,"time":"2025-11-19T20:24:34.172Z","status":"INFORM","message":"Upgrading migration package schema version.","line":"Wed-19-Nov-2025 08:24:34.172 PM   INFORM   Thread-8   409504 : Upgrading migration package schema version."}],"totalCount":188,"limit":100,"offset":0}`

	responseV3One := `{
		"result": "Successful export",
		"id": 1,
		"packageName": "ServerConfigPackage.fsconfig",
		"type": "EXPORT",
		"username": "admin",
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
			fmeflowBuild:       26000,
		},
		{
			name:         "500 bad status code",
			statusCode:   http.StatusInternalServerError,
			wantErrText:  "500 Internal Server Error",
			args:         []string{"migration", "tasks"},
			fmeflowBuild: 26000,
		},
		{
			name:         "404 bad status code",
			statusCode:   http.StatusNotFound,
			wantErrText:  "404 Not Found",
			args:         []string{"migration", "tasks"},
			fmeflowBuild: 26000,
		},
		{
			name:            "get migration tasks table output",
			statusCode:      http.StatusOK,
			args:            []string{"migration", "tasks"},
			wantOutputRegex: "^[\\s]*ID[\\s]*TYPE[\\s]*USERNAME[\\s]*START TIME[\\s]*END TIME[\\s]*STATUS[\\s]*2[\\s]*EXPORT[\\s]*admin[\\s]*2022-11-09 22:40:51 \\+0000 UTC[\\s]*2022-11-09 22:40:52 \\+0000 UTC[\\s]*SUCCESS[\\s]*1[\\s]*EXPORT[\\s]*admin[\\s]*2022-11-09 21:59:22 \\+0000 UTC[\\s]*2022-11-09 21:59:24 \\+0000 UTC[\\s]*SUCCESS[\\s]*$",
			body:            responseV4,
			fmeflowBuild:    26000,
		},
		{
			name:            "get migration tasks no headers",
			statusCode:      http.StatusOK,
			body:            responseV4,
			args:            []string{"migration", "tasks", "--no-headers"},
			wantOutputRegex: "^[\\s]*2[\\s]*EXPORT[\\s]*admin[\\s]*2022-11-09 22:40:51 \\+0000 UTC[\\s]*2022-11-09 22:40:52 \\+0000 UTC[\\s]*SUCCESS[\\s]*1[\\s]*EXPORT[\\s]*admin[\\s]*2022-11-09 21:59:22 \\+0000 UTC[\\s]*2022-11-09 21:59:24 \\+0000 UTC[\\s]*SUCCESS[\\s]*$",
			fmeflowBuild:    26000,
		},
		{
			name:            "get migration by id",
			statusCode:      http.StatusOK,
			body:            responseV3One,
			args:            []string{"migration", "tasks", "--id", "1"},
			wantOutputRegex: "^[\\s]*ID[\\s]*TYPE[\\s]*USERNAME[\\s]*START TIME[\\s]*END TIME[\\s]*STATUS[\\s]*1[\\s]*EXPORT[\\s]*admin[\\s]*2022-11-09 21:59:22 \\+0000 UTC[\\s]*2022-11-09 21:59:24 \\+0000 UTC[\\s]*SUCCESS[\\s]*$",
			fmeflowBuild:    26000,
		},
		{
			name:            "get migration log",
			statusCode:      http.StatusOK,
			body:            migrationLogResponse,
			args:            []string{"migration", "tasks", "--id", "1", "--log"},
			wantOutputRegex: migrationLogResponse,
			fmeflowBuild:    26000,
		},
		{
			name:           "get migration tasks json",
			statusCode:     http.StatusOK,
			args:           []string{"migration", "tasks", "--json"},
			body:           responseV4,
			wantOutputJson: responseV4,
			fmeflowBuild:   26000,
		},
		{
			name:            "get migration by id custom columns",
			statusCode:      http.StatusOK,
			body:            responseV3One,
			args:            []string{"migration", "tasks", "--id", "1", "--output", "custom-columns=PACKAGE:.packageName,STATUS:.status"},
			wantOutputRegex: "^[\\s]*PACKAGE[\\s]*STATUS[\\s]*ServerConfigPackage.fsconfig[\\s]*SUCCESS[\\s]*$",
			fmeflowBuild:    26000,
		},
		{
			name:            "get migrations custom columns",
			statusCode:      http.StatusOK,
			body:            responseV4,
			args:            []string{"migration", "tasks", "--output", "custom-columns=PACKAGE:.packageName,STATUS:.status"},
			wantOutputRegex: "^[\\s]*PACKAGE[\\s]*STATUS[\\s]*ServerConfigPackage.fsconfig[\\s]*SUCCESS[\\s]*ServerConfigPackage.fsconfig[\\s]*SUCCESS[\\s]*$",
			fmeflowBuild:    26000,
		},
		{
			name:           "get migration log parsed json",
			statusCode:     http.StatusOK,
			body:           migrationLogParsedResponse,
			args:           []string{"migration", "tasks", "--id", "1", "--log", "--json"},
			wantOutputJson: migrationLogParsedResponse,
			fmeflowBuild:   26000,
		},
	}

	runTests(cases, t)

}
