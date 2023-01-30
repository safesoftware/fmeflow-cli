package cmd

import (
	"net/http"
	"testing"
)

func TestDeploymentParametersUpdate(t *testing.T) {
	parameterDoesNotExistBody := `{
		"message": "Unauthorized request by user admin due to lack of proper permissions or the object does not exist."
	  }`
	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"deploymentparameters", "update", "--name", "myDep", "--value", "myValue", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"deploymentparameters", "update", "--name", "myDep", "--value", "myValue"},
		},
		{
			name:        "missing name flag",
			wantErrText: "required flag(s) \"name\", \"value\" not set",
			args:        []string{"deploymentparameters", "update"},
		},
		{
			name:        "missing value flag",
			wantErrText: "required flag(s) \"value\" not set",
			args:        []string{"deploymentparameters", "update", "--name", "myDep"},
		},
		{
			name:            "update parameter",
			statusCode:      http.StatusNoContent,
			args:            []string{"deploymentparameters", "update", "--name", "myDep", "--value", "myValue"},
			wantOutputRegex: "^Deployment Parameter successfully updated.[\\s]*$",
		},
		{
			name:        "parameter does not exist",
			statusCode:  http.StatusConflict,
			body:        parameterDoesNotExistBody,
			args:        []string{"deploymentparameters", "update", "--name", "myDep", "--value", "myValue"},
			wantErrText: "Unauthorized request by user admin due to lack of proper permissions or the object does not exist.",
		},
	}

	runTests(cases, t)

}
