package cmd

import (
	"net/http"
	"testing"
)

func TestProjects(t *testing.T) {
	responseV3 := `{
		"offset": -1,
		"limit": -1,
		"totalCount": 2,
		"items": [
		  {
			"owner": "admin",
			"apps": [
			  {
				"name": "test_app"
			  }
			],
			"appsuites": [
			  {
				"name": "test_appsuite"
			  }
			],
			"automations": [
			  {
				"name": "TestAutomation",
				"uuid": "6f865528-fe20-4275-bcd4-10bbc5f7486a"
			  }
			],
			"automationApps": [
			  {
				"name": "TestAutomationApp",
				"uuid": "6f865528-fe20-4275-bcd4-10bbc5f7486a"
			  }
			],
			"lastSaveDate": "2023-01-17T16:34:45Z",
			
			"cleanupTasks": [
			  {
				"name": "Delete_Automation_Download_Logs",
				"category": "Utilities"
			  }
			],
			"customFormats": [
			  {
				"name": "testformat",
				"repositoryName": "Samples"
			  }
			],
			"customTransformers": [
			  {
				"name": "testTransformer.fmx",
				"repositoryName": "Samples"
			  }
			],
			"projects": [
			  {
				"name": "test2"
			  }
			],
			"publications": [
			  {
				"name": "testPublication"
			  }
			],
			"topics": [
			  {
				"name": "DATADOWNLOAD_ASYNC_JOB_FAILURE"
			  }
			],
			"resourceConnections": [
			  {
				"name": "resourceConnection"
			  }
			],
			"resourcePaths": [
			  {
				"path": "/",
				"name": "FME_SHAREDRESOURCE_DATA"
			  }
			],
			"roles": [
			  {
				"name": "fmeuser"
			  }
			],
			"subscriptions": [
			  {
				"name": "Dashboards_AverageRunningTime"
			  }
			],
			"streams": [
			  {
				"name": "TestStream",
				"uuid": "6f865528-fe20-4275-bcd4-10bbc5f7486a"
			  }
			],
			"templates": [
			  {
				"name": "testtemplate",
				"repositoryName": "Samples"
			  }
			],
			"tokens": [
			  {
				"name": "token",
				"userName": "admin"
			  }
			],
			"description": "My Description",
			"sharable": true,
			"readme": "A cool project",
			"userName": "admin",
			"version": "1.0.0",
			"fmeHubPublisherUid": "",
			"uid": "",
			"repositories": [
			  {
				"name": "Utilities"
			  }
			],
			
			"hasIcon": false,
			"schedules": [
			  {
				"name": "Backup_Configuration",
				"category": "Utilities"
			  }
			],
			"name": "Test123",
			"accounts": [
			  {
				"name": "admin"
			  }
			],
			"workspaces": [
			  {
				"name": "austinApartments.fmw",
				"repositoryName": "Samples"
			  },
			  {
				"name": "backupConfiguration.fmw",
				"repositoryName": "Utilities"
			  }
			],
			"connections": [
			  {
				"name": "test"
			  }
			]
		  },
		  {
			"owner": "admin",
			"uid": "",
			"lastSaveDate": "2023-01-17T16:30:57Z",
			"hasIcon": false,
			"name": "test2",
			"description": "a",
			"sharable": true,
			"readme": "a",
			"workspaces": [
			  {
				"name": "austinDownload.fmw",
				"repositoryName": "Samples"
			  }
			],
			"userName": "admin",
			"version": "1.0.0",
			"fmeHubPublisherUid": ""
		  }
		]
	  }`

	responseSingleProject := `{
		"owner": "admin",
		"apps": [
		  {
			"name": "test_app"
		  }
		],
		"appsuites": [
		  {
			"name": "test_appsuite"
		  }
		],
		"automations": [
		  {
			"name": "TestAutomation",
			"uuid": "6f865528-fe20-4275-bcd4-10bbc5f7486a"
		  }
		],
		"automationApps": [
		  {
			"name": "TestAutomationApp",
			"uuid": "6f865528-fe20-4275-bcd4-10bbc5f7486a"
		  }
		],
		"lastSaveDate": "2023-01-17T16:34:45Z",
		
		"cleanupTasks": [
		  {
			"name": "Delete_Automation_Download_Logs",
			"category": "Utilities"
		  }
		],
		"customFormats": [
		  {
			"name": "testformat",
			"repositoryName": "Samples"
		  }
		],
		"customTransformers": [
		  {
			"name": "testTransformer.fmx",
			"repositoryName": "Samples"
		  }
		],
		"projects": [
		  {
			"name": "test2"
		  }
		],
		"publications": [
		  {
			"name": "testPublication"
		  }
		],
		"topics": [
		  {
			"name": "DATADOWNLOAD_ASYNC_JOB_FAILURE"
		  }
		],
		"resourceConnections": [
		  {
			"name": "resourceConnection"
		  }
		],
		"resourcePaths": [
		  {
			"path": "/",
			"name": "FME_SHAREDRESOURCE_DATA"
		  }
		],
		"roles": [
		  {
			"name": "fmeuser"
		  }
		],
		"subscriptions": [
		  {
			"name": "Dashboards_AverageRunningTime"
		  }
		],
		"streams": [
		  {
			"name": "TestStream",
			"uuid": "6f865528-fe20-4275-bcd4-10bbc5f7486a"
		  }
		],
		"templates": [
		  {
			"name": "testtemplate",
			"repositoryName": "Samples"
		  }
		],
		"tokens": [
		  {
			"name": "token",
			"userName": "admin"
		  }
		],
		"description": "My Description",
		"sharable": true,
		"readme": "A cool project",
		"userName": "admin",
		"version": "1.0.0",
		"fmeHubPublisherUid": "",
		"uid": "",
		"repositories": [
		  {
			"name": "Utilities"
		  }
		],
		"hasIcon": false,
		"schedules": [
		  {
			"name": "Backup_Configuration",
			"category": "Utilities"
		  }
		],
		"name": "Test123",
		"accounts": [
		  {
			"name": "admin"
		  }
		],
		"workspaces": [
		  {
			"name": "austinApartments.fmw",
			"repositoryName": "Samples"
		  },
		  {
			"name": "backupConfiguration.fmw",
			"repositoryName": "Utilities"
		  }
		],
		"connections": [
		  {
			"name": "test"
		  }
		]
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"projects", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"projects"},
		},
		{
			name:        "project not found",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found: check that the specified project exists",
			args:        []string{"projects", "--name", "Samples123"},
		},
		{
			name:            "get project table output",
			statusCode:      http.StatusOK,
			args:            []string{"projects"},
			body:            responseV3,
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*LAST SAVED[\\s]*Test123[\\s]*admin[\\s]*My Description[\\s]*2023-01-17 16:34:45 \\+0000 UTC[\\s]*test2[\\s]*admin[\\s]*a[\\s]*2023-01-17 16:30:57 \\+0000 UTC[\\s]*$",
		},
		{
			name:           "get project json  output",
			statusCode:     http.StatusOK,
			args:           []string{"projects", "--json"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
		{
			name:            "get single project",
			statusCode:      http.StatusOK,
			body:            responseSingleProject,
			args:            []string{"projects", "--name", "Test123"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*LAST SAVED[\\s]*Test123[\\s]*admin[\\s]*My Description[\\s]*2023-01-17 16:34:45 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get project from owner",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"projects", "--owner", "admin"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*LAST SAVED[\\s]*Test123[\\s]*admin[\\s]*My Description[\\s]*2023-01-17 16:34:45 \\+0000 UTC[\\s]*test2[\\s]*admin[\\s]*a[\\s]*2023-01-17 16:30:57 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get project custom columns",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "--output=custom-columns=NAME:.name,WORKSPACE:.workspaces[*].name"},
			body:            responseV3,
			wantOutputRegex: "^[\\s]*NAME[\\s]*WORKSPACE[\\s]*Test123[\\s]*austinApartments.fmw backupConfiguration.fmw[\\s]*test2[\\s]*austinDownload.fmw[\\s]*$",
		},
	}

	runTests(cases, t)

}
