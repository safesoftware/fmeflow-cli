package cmd

import (
	"net/http"
	"testing"
)

func TestRepositories(t *testing.T) {
	responseV3 := `{
		"offset": -1,
		"limit": -1,
		"totalCount": 4,
		"items": [
		  {
			"owner": "admin",
			"name": "Dashboards",
			"description": "FME Server Dashboards Repository",
			"sharable": true
		  },
		  {
			"owner": "admin",
			"name": "Samples",
			"description": "FME Server Samples Repository",
			"sharable": true
		  },
		  {
			"owner": "admin",
			"name": "test",
			"description": "",
			"sharable": true
		  },
		  {
			"owner": "admin",
			"name": "Utilities",
			"description": "FME Server Utilities Repository",
			"sharable": true
		  }
		]
	  }`

	responseV3SingleRepo := `{
		"owner": "admin",
		"name": "Samples",
		"description": "FME Server Samples Repository",
		"sharable": true
	  }`

	responseV4 := `{
		"items": [
		  {
			"name": "Dashboards",
			"description": "FME Server Dashboards Repository",
			"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"owner": "admin",
			"sharable": true,
			"totalFileSize": 0,
			"fileCount": 0,
			"workspaceCount": 6,
			"customTransformerCount": 0,
			"customFormatCount": 0,
			"templateCount": 0
		  },
		  {
			"name": "Samples",
			"description": "FME Server Samples Repository",
			"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"owner": "admin",
			"sharable": true,
			"totalFileSize": 30276633,
			"fileCount": 31,
			"workspaceCount": 4,
			"customTransformerCount": 0,
			"customFormatCount": 0,
			"templateCount": 0
		  },
		  {
			"name": "Utilities",
			"description": "FME Server Utilities Repository",
			"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"owner": "admin",
			"sharable": true,
			"totalFileSize": 0,
			"fileCount": 0,
			"workspaceCount": 1,
			"customTransformerCount": 0,
			"customFormatCount": 0,
			"templateCount": 0
		  }
		],
		"totalCount": 3,
		"limit": 100,
		"offset": 0
	  }`

	responseV4SingleRepo := `{
		"name": "Samples",
		"description": "FME Server Samples Repository",
		"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
		"owner": "admin",
		"sharable": true,
		"totalFileSize": 30276633,
		"fileCount": 31,
		"workspaceCount": 4,
		"customTransformerCount": 0,
		"customFormatCount": 0,
		"templateCount": 0
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"repositories", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"repositories"},
		},
		{
			name:            "get repositories table output V4",
			statusCode:      http.StatusOK,
			args:            []string{"repositories", "--api-version", "v4"},
			body:            responseV4,
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*WORKSPACES[\\s]*Dashboards[\\s]*admin[\\s]*FME Server Dashboards Repository[\\s]*6[\\s]*Samples[\\s]*admin[\\s]*FME Server Samples Repository[\\s]*4[\\s]*Utilities[\\s]*admin[\\s]*FME Server Utilities Repository[\\s]*1[\\s]*$",
		},
		{
			name:            "get single repository V4",
			statusCode:      http.StatusOK,
			body:            responseV4SingleRepo,
			args:            []string{"repositories", "--name", "Samples", "--api-version", "v4"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*WORKSPACES[\\s]*Samples[\\s]*admin[\\s]*FME Server Samples Repository[\\s]*4[\\s]*$",
		},
		{
			name:            "get repository with filter string V4",
			statusCode:      http.StatusOK,
			body:            responseV4,
			args:            []string{"repositories", "--filter-string", "admin", "--api-version", "v4"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*WORKSPACES[\\s]*Dashboards[\\s]*admin[\\s]*FME Server Dashboards Repository[\\s]*6[\\s]*Samples[\\s]*admin[\\s]*FME Server Samples Repository[\\s]*4[\\s]*Utilities[\\s]*admin[\\s]*FME Server Utilities Repository[\\s]*1[\\s]*$",
		},
		{
			name:            "get repositories custom columns V4",
			statusCode:      http.StatusOK,
			args:            []string{"repositories", "--output=custom-columns=NAME:$.name", "--api-version", "v4"},
			body:            responseV4,
			wantOutputRegex: "^[\\s]*NAME[\\s]*Dashboards[\\s]*Samples[\\s]*Utilities[\\s]*$",
		},
		{
			name:        "repository not found V3",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found: check that the specified repository exists",
			args:        []string{"repositories", "--name", "Samples123", "--api-version", "v3"},
		},
		{
			name:            "get repositories table output V3",
			statusCode:      http.StatusOK,
			args:            []string{"repositories", "--api-version", "v3"},
			body:            responseV3,
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*SHARABLE[\\s]*Dashboards[\\s]*admin[\\s]*FME Server Dashboards Repository[\\s]*true[\\s]*Samples[\\s]*admin[\\s]*FME Server Samples Repository[\\s]*true[\\s]*test[\\s]*admin[\\s]*true[\\s]*Utilities[\\s]*admin[\\s]*FME Server Utilities Repository[\\s]*true[\\s]*$",
		},
		{
			name:            "get single repository V3",
			statusCode:      http.StatusOK,
			body:            responseV3SingleRepo,
			args:            []string{"repositories", "--name", "Samples", "--api-version", "v3"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*SHARABLE[\\s]*Samples[\\s]*admin[\\s]*FME Server Samples Repository[\\s]*true[\\s]*$",
		},
		{
			name:            "get repository from owner V3",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"repositories", "--owner", "admin", "--api-version", "v3"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*SHARABLE[\\s]*Dashboards[\\s]*admin[\\s]*FME Server Dashboards Repository[\\s]*true[\\s]*Samples[\\s]*admin[\\s]*FME Server Samples Repository[\\s]*true[\\s]*test[\\s]*admin[\\s]*true[\\s]*Utilities[\\s]*admin[\\s]*FME Server Utilities Repository[\\s]*true[\\s]*$",
		},
		{
			name:            "get repositories custom columns V3",
			statusCode:      http.StatusOK,
			args:            []string{"repositories", "--output=custom-columns=NAME:$.name", "--api-version", "v3"},
			body:            responseV3,
			wantOutputRegex: "^[\\s]*NAME[\\s]*Dashboards[\\s]*Samples[\\s]*test[\\s]*Utilities[\\s]*$",
		},
	}

	runTests(cases, t)

}
