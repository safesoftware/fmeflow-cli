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
			"updated": "2023-01-18T23:04:25.764Z"
		  },
		  {
			"name": "test",
			"value": "test value",
			"type": "text",
			"owner": "admin",
			"updated": "2023-01-19T17:56:54.472Z"
		  }
		],
		"totalCount": 2,
		"limit": 100,
		"offset": 0
	  }`

	responseSingle := `{
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
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*TYPE[\\s]*VALUE[\\s]*LAST UPDATED[\\s]*myDep[\\s]*admin[\\s]*text[\\s]*myVal[\\s]*2023-01-18 23:04:25.764 \\+0000 UTC[\\s]*test[\\s]*admin[\\s]*text[\\s]*test[\\s]*value[\\s]*2023-01-19 17:56:54.472 \\+0000 UTC[\\s]*$",
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
			body:            responseSingle,
			args:            []string{"deploymentparameters", "--name", "myDep"},
			wantOutputRegex: "^[\\s]*NAME[\\s]*OWNER[\\s]*TYPE[\\s]*VALUE[\\s]*LAST UPDATED[\\s]*myDep[\\s]*admin[\\s]*text[\\s]*myVal[\\s]*2023-01-18 23:04:25.764 \\+0000 UTC[\\s]*$",
		},
		{
			name:            "get deployment parameters custom columns V4",
			statusCode:      http.StatusOK,
			args:            []string{"repositories", "--output=custom-columns=NAME:.name"},
			body:            response,
			wantOutputRegex: "^[\\s]*NAME[\\s]*myDep[\\s]*test[\\s]*$",
		},
	}

	runTests(cases, t)

}
