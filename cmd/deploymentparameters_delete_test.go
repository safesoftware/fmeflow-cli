package cmd

import (
	"net/http"
	"testing"
)

func TestDeploymentParametersDelete(t *testing.T) {
	paramMissingBody := `{
		"message": "Unauthorized request by user admin due to lack of proper permissions or the object does not exist."
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"deploymentparameters", "delete", "--name", "myDep", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"deploymentparameters", "delete", "--name", "myDep", "--no-prompt"},
		},
		{
			name:        "missing flag",
			wantErrText: "required flag(s) \"name\" not set",
			args:        []string{"deploymentparameters", "delete"},
		},
		{
			name:            "delete repository V4",
			statusCode:      http.StatusNoContent,
			args:            []string{"deploymentparameters", "delete", "--name", "myDep", "--no-prompt"},
			wantOutputRegex: "^Deployment Parameter successfully deleted.[\\s]*$",
		},
		{
			name:        "delete repository not found V4",
			statusCode:  http.StatusNotFound,
			body:        paramMissingBody,
			args:        []string{"deploymentparameters", "delete", "--name", "myDep", "--no-prompt"},
			wantErrText: "Unauthorized request by user admin due to lack of proper permissions or the object does not exist.",
		},
	}

	runTests(cases, t)

}
