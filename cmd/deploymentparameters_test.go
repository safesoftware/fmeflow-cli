package cmd

import (
	"net/http"
	"testing"
)

func TestDeploymentParameters(t *testing.T) {
	response := `{
		"items": [
		  {
			"name": "myDep",
			"value": "myVal",
			"type": "text",
			"owner": "admin",
			"resourceMissing": false,
			"updated": "2023-01-18T23:04:25.764Z"
		  },
		  {
			"type": "dropdown",
			"name": "testdb",
			"value": "db",
			"type": "dropdown",
			"owner": "admin",
			"updated": "2024-03-12T20:39:12.149Z",
			"resourceMissing": false,
			"choiceSettings": {
			  "choiceSet": "dbConnections",
			  "family": "PostgreSQL"
			}
		  },
		  {
			"type": "dropdown",
			"name": "testweb",
			"value": "aaa",
			"type": "dropdown",
			"owner": "admin",
			"updated": "2024-03-15T17:42:56.752Z",
			"resourceMissing": true,
			"choiceSettings": {
			  "choiceSet": "webConnections",
			  "services": [
				"Slack"
			  ]
			}
		 }
		],
		
		"totalCount": 3,
		"limit": 100,
		"offset": 0
	  }`

	responseSingleText := `{
		"name": "myDep",
		"value": "myVal",
		"type": "text",
		"owner": "admin",
		"updated": "2023-01-18T23:04:25.764Z"
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"deploymentparameters", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"deploymentparameters"},
		},
		{
			name:            "get deployment parameters table output",
			statusCode:      http.StatusOK,
			args:            []string{"deploymentparameters"},
			body:            response,
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*TYPE[\\s]*VALUE[\\s]*LAST UPDATED[\\s]*myDep[\\s]*admin[\\s]*text[\\s]*myVal[\\s]*2023-01-18 23:04:25.764 \\+0000 UTC[\\s]*testdb[\\s]*admin[\\s]*dropdown[\\s]*db[\\s]*2024-03-12 20:39:12.149 \\+0000 UTC[\\s]*testweb[\\s]*admin[\\s]*dropdown[\\s]*aaa[\\s]*2024-03-15 17:42:56.752 \\+0000 UTC[\\s]*$",
		},
		{
			name:           "get deployment parameters json output",
			statusCode:     http.StatusOK,
			args:           []string{"deploymentparameters", "--json"},
			body:           response,
			wantOutputJson: response,
		},
		{
			name:            "get single parameter",
			statusCode:      http.StatusOK,
			body:            responseSingleText,
			args:            []string{"deploymentparameters", "--name", "myDep"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*TYPE[\\s]*VALUE[\\s]*LAST UPDATED[\\s]*myDep[\\s]*admin[\\s]*text[\\s]*myVal[\\s]*2023-01-18 23:04:25.764 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get deployment parameters custom columns V4",
			statusCode:      http.StatusOK,
			args:            []string{"repositories", "--output=custom-columns=NAME:.name"},
			body:            response,
			wantOutputRegex: "^[\\s]*NAME[\\s]*myDep[\\s]*testdb[\\s]*testweb[\\s]*$",
		},
		{
			name:            "get deployment parameters custom columns no headers V4",
			statusCode:      http.StatusOK,
			args:            []string{"repositories", "--output=custom-columns=NAME:.name", "--no-headers"},
			body:            response,
			wantOutputRegex: "^[\\s]*myDep[\\s]*testdb[\\s]*testweb[\\s]*$",
		},
	}

	runTests(cases, t)

}
