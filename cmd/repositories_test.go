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
			name:        "repository not found",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found: check that the specified repository exists",
			args:        []string{"repositories", "--name", "Samples123"},
		},
		{
			name:            "get repositories table output",
			statusCode:      http.StatusOK,
			args:            []string{"repositories"},
			body:            responseV3,
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*SHARABLE[\\s]*Dashboards[\\s]*admin[\\s]*FME Server Dashboards Repository[\\s]*true[\\s]*Samples[\\s]*admin[\\s]*FME Server Samples Repository[\\s]*true[\\s]*test[\\s]*admin[\\s]*true[\\s]*Utilities[\\s]*admin[\\s]*FME Server Utilities Repository[\\s]*true[\\s]*$",
		},
		{
			name:            "get single repository",
			statusCode:      http.StatusOK,
			body:            responseV3SingleRepo,
			args:            []string{"repositories", "--name", "Samples"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*SHARABLE[\\s]*Samples[\\s]*admin[\\s]*FME Server Samples Repository[\\s]*true[\\s]*$",
		},
		{
			name:            "get repository from owner",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"repositories", "--owner", "admin"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*DESCRIPTION[\\s]*SHARABLE[\\s]*Dashboards[\\s]*admin[\\s]*FME Server Dashboards Repository[\\s]*true[\\s]*Samples[\\s]*admin[\\s]*FME Server Samples Repository[\\s]*true[\\s]*test[\\s]*admin[\\s]*true[\\s]*Utilities[\\s]*admin[\\s]*FME Server Utilities Repository[\\s]*true[\\s]*$",
		},
		{
			name:            "get repositories custom columns",
			statusCode:      http.StatusOK,
			args:            []string{"repositories", "--output=custom-columns=NAME:.name"},
			body:            responseV3,
			wantOutputRegex: "^[\\s]*NAME[\\s]*Dashboards[\\s]*Samples[\\s]*test[\\s]*Utilities[\\s]*$",
		},
	}

	runTests(cases, t)

}
