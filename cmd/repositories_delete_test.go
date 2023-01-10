package cmd

import (
	"net/http"
	"testing"
)

func TestRepositoriesDelete(t *testing.T) {

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
			name:            "delete repository",
			statusCode:      http.StatusNoContent,
			args:            []string{"repositories", "delete", "--name", "MyRepo", "--no-prompt"},
			wantOutputRegex: "^Repository successfully deleted.[\\s]*$",
		},
		{
			name:        "delete repository not found",
			statusCode:  http.StatusNotFound,
			args:        []string{"repositories", "delete", "--name", "MyRepo", "--no-prompt"},
			wantErrText: "404 Not Found: The repository does not exist",
		},
	}

	runTests(cases, t)

}
