package cmd

import (
	"net/http"
	"testing"
)

func TestProjects(t *testing.T) {
	responsev4 := `{
		"items": [
		  {
			"id": "a64297e7-a119-4e10-ac37-5d0bba12194b",
			"name": "test1",
			"hubUid": "",
			"hubPublisherUid": "",
			"description": "test1",
			"readme": "",
			"version": "1.0.0",
			"lastUpdated": "2024-03-26T18:44:30.713Z",
			"owner": "admin",
			"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"shareable": true,
			"lastUpdateUser": "admin",
			"lastUpdateUserID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"hasIcon": false
		  },
		  {
			"id": "5a70afe2-db56-4a15-8ea8-559a10c326ec",
			"name": "test2",
			"hubUid": null,
			"hubPublisherUid": null,
			"description": "test2",
			"readme": null,
			"version": "1.0.0",
			"lastUpdated": "2024-03-26T18:45:32.991Z",
			"owner": "admin",
			"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"shareable": true,
			"lastUpdateUser": "admin",
			"lastUpdateUserID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"hasIcon": false
		  },
		  {
			"id": "0adc89c8-5eda-4cbb-a599-eefd66b1c8b9",
			"name": "testAutomation",
			"hubUid": null,
			"hubPublisherUid": null,
			"description": "test",
			"readme": null,
			"version": "1.0.0",
			"lastUpdated": "2024-03-26T22:52:38.092Z",
			"owner": "admin",
			"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"shareable": true,
			"lastUpdateUser": "admin",
			"lastUpdateUserID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"hasIcon": false
		  }
		],
		"totalCount": 3,
		"limit": 100,
		"offset": 0
	  }`

	singleResponseV4 := `{
		"id": "a64297e7-a119-4e10-ac37-5d0bba12194b",
		"name": "test1",
		"hubUid": "",
		"hubPublisherUid": "",
		"description": "test1",
		"readme": "",
		"version": "1.0.0",
		"lastUpdated": "2024-03-26T18:44:30.713Z",
		"owner": "admin",
		"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
		"shareable": true,
		"lastUpdateUser": "admin",
		"lastUpdateUserID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
		"hasIcon": false
	  }`

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

	paramMissingBody := `{
		"message": "Unauthorized request by user admin due to lack of proper permissions or the object does not exist."
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
			args:        []string{"projects", "--api-version=v3"},
		},
		{
			name:        "project not found id v4",
			args:        []string{"projects", "--id", "a64297e7-a119-4e10-ac37-5d0bba121940"},
			statusCode:  http.StatusForbidden,
			body:        paramMissingBody,
			wantErrText: "Unauthorized request by user admin due to lack of proper permissions or the object does not exist.",
		},
		{
			name:            "get project table output v4",
			args:            []string{"projects"},
			body:            responsev4,
			statusCode:      http.StatusOK,
			wantOutputRegex: "^[\\s]*ID[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*LAST UPDATED[\\s]*a64297e7-a119-4e10-ac37-5d0bba12194b[\\s]*test1[\\s]*admin[\\s]*test1[\\s]*2024-03-26 18:44:30.713 \\+0000 UTC[\\s]*5a70afe2-db56-4a15-8ea8-559a10c326ec[\\s]*test2[\\s]*admin[\\s]*test2[\\s]*2024-03-26 18:45:32.991 \\+0000 UTC[\\s]*0adc89c8-5eda-4cbb-a599-eefd66b1c8b9[\\s]*testAutomation[\\s]*admin[\\s]*test[\\s]*2024-03-26 22:52:38.092 \\+0000 UTC[\\s]*$",
		},
		{
			name:           "get project json output v4",
			args:           []string{"projects", "--json"},
			body:           responsev4,
			statusCode:     http.StatusOK,
			wantOutputJson: responsev4,
		},
		{
			name:            "get single project v4",
			args:            []string{"projects", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b"},
			body:            singleResponseV4,
			statusCode:      http.StatusOK,
			wantOutputRegex: "^[\\s]*ID[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*LAST UPDATED[\\s]*a64297e7-a119-4e10-ac37-5d0bba12194b[\\s]*test1[\\s]*admin[\\s]*test1[\\s]*2024-03-26 18:44:30.713 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get project from owner v4",
			args:            []string{"projects", "--owner", "admin"},
			body:            responsev4,
			statusCode:      http.StatusOK,
			wantOutputRegex: "^[\\s]*ID[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*LAST UPDATED[\\s]*a64297e7-a119-4e10-ac37-5d0bba12194b[\\s]*test1[\\s]*admin[\\s]*test1[\\s]*2024-03-26 18:44:30.713 \\+0000 UTC[\\s]*5a70afe2-db56-4a15-8ea8-559a10c326ec[\\s]*test2[\\s]*admin[\\s]*test2[\\s]*2024-03-26 18:45:32.991 \\+0000 UTC[\\s]*0adc89c8-5eda-4cbb-a599-eefd66b1c8b9[\\s]*testAutomation[\\s]*admin[\\s]*test[\\s]*2024-03-26 22:52:38.092 \\+0000 UTC[\\s]*$",
			wantFormParams:  map[string]string{"filterString": "admin", "filterProperties": "owner"},
		},
		{
			name:            "get project custom columns v4",
			args:            []string{"projects", "--output=custom-columns=ID:.id,NAME:.name"},
			body:            responsev4,
			statusCode:      http.StatusOK,
			wantOutputRegex: "^[\\s]*ID[\\s]*NAME[\\s]*a64297e7-a119-4e10-ac37-5d0bba12194b[\\s]*test1[\\s]*5a70afe2-db56-4a15-8ea8-559a10c326ec[\\s]*test2[\\s]*0adc89c8-5eda-4cbb-a599-eefd66b1c8b9[\\s]*testAutomation[\\s]*$",
		},
		{
			name:        "project not found v3",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found: check that the specified project exists",
			args:        []string{"projects", "--name", "Samples123", "--api-version=v3"},
		},
		{
			name:            "get project table output v3",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "--api-version=v3"},
			body:            responseV3,
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*LAST SAVED[\\s]*Test123[\\s]*admin[\\s]*My Description[\\s]*2023-01-17 16:34:45 \\+0000 UTC[\\s]*test2[\\s]*admin[\\s]*a[\\s]*2023-01-17 16:30:57 \\+0000 UTC[\\s]*$",
		},
		{
			name:           "get project json output v3",
			statusCode:     http.StatusOK,
			args:           []string{"projects", "--json", "--api-version=v3"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
		{
			name:            "get single project v3",
			statusCode:      http.StatusOK,
			body:            responseSingleProject,
			args:            []string{"projects", "--name", "Test123", "--api-version=v3"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*LAST SAVED[\\s]*Test123[\\s]*admin[\\s]*My Description[\\s]*2023-01-17 16:34:45 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get project from owner v3",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"projects", "--owner", "admin", "--api-version=v3"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*LAST SAVED[\\s]*Test123[\\s]*admin[\\s]*My Description[\\s]*2023-01-17 16:34:45 \\+0000 UTC[\\s]*test2[\\s]*admin[\\s]*a[\\s]*2023-01-17 16:30:57 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get project custom columns v3",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "--output=custom-columns=NAME:.name,WORKSPACE:.workspaces[*].name", "--api-version=v3"},
			body:            responseV3,
			wantOutputRegex: "^[\\s]*NAME[\\s]*WORKSPACE[\\s]*Test123[\\s]*austinApartments.fmw backupConfiguration.fmw[\\s]*test2[\\s]*austinDownload.fmw[\\s]*$",
		},
	}

	runTests(cases, t)

}
