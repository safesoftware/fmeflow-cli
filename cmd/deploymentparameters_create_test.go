package cmd

import (
	"net/http"
	"testing"
)

func TestDeploymentParametersCreate(t *testing.T) {
	parameterExistsBody := `{
		"message": "A deployment parameter with name \"myDep\" already exists."
	  }`
	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue"},
		},
		{
			name:        "missing name flag",
			wantErrText: "required flag(s) \"name\", \"value\" not set",
			args:        []string{"deploymentparameters", "create"},
		},
		{
			name:        "missing value flag",
			wantErrText: "required flag(s) \"value\" not set",
			args:        []string{"deploymentparameters", "create", "--name", "myDep"},
		},
		{
			name:            "create parameter",
			statusCode:      http.StatusCreated,
			args:            []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue"},
			wantOutputRegex: "^Deployment Parameter successfully created.[\\s]*$",
		},
		{
			name:        "parameter already exists",
			statusCode:  http.StatusConflict,
			body:        parameterExistsBody,
			args:        []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue"},
			wantErrText: "A deployment parameter with name \"myDep\" already exists.",
		},
	}

	runTests(cases, t)

}
