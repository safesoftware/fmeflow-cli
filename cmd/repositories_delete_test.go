package cmd

import (
	"net/http"
	"testing"
)

func TestRepositoriesDelete(t *testing.T) {
	repoMissingBodyV4 := `{
		"message": "Unauthorized request by user admin due to lack of proper permissions or the object does not exist."
	  }`

	repoMissingBodyV3 := `{
		"message": "Repository MyRepo does not exist."
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"repositories", "delete", "--name", "MyRepo", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"repositories", "delete", "--name", "MyRepo", "--no-prompt"},
		},
		{
			name:        "missing flag",
			wantErrText: "required flag(s) \"name\" not set",
			args:        []string{"repositories", "delete"},
		},
		{
			name:            "delete repository V4",
			statusCode:      http.StatusNoContent,
			args:            []string{"repositories", "delete", "--name", "MyRepo", "--no-prompt", "--api-version", "v4"},
			wantOutputRegex: "^Repository successfully deleted.[\\s]*$",
		},
		{
			name:        "delete repository not found V4",
			statusCode:  http.StatusNotFound,
			body:        repoMissingBodyV4,
			args:        []string{"repositories", "delete", "--name", "MyRepo", "--no-prompt", "--api-version", "v4"},
			wantErrText: "Unauthorized request by user admin due to lack of proper permissions or the object does not exist.",
		},
		{
			name:            "delete repository V3",
			statusCode:      http.StatusNoContent,
			args:            []string{"repositories", "delete", "--name", "MyRepo", "--no-prompt", "--api-version", "v3"},
			wantOutputRegex: "^Repository successfully deleted.[\\s]*$",
		},
		{
			name:        "delete repository not found V3",
			statusCode:  http.StatusNotFound,
			body:        repoMissingBodyV3,
			args:        []string{"repositories", "delete", "--name", "MyRepo", "--no-prompt", "--api-version", "v3"},
			wantErrText: "Repository MyRepo does not exist.",
		},
	}

	runTests(cases, t)

}
