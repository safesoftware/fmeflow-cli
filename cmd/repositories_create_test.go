package cmd

import (
	"net/http"
	"testing"
)

func TestRepositoriesCreate(t *testing.T) {

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"repositories", "create", "--name", "MyRepo", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"repositories", "create", "--name", "MyRepo"},
		},
		{
			name:        "missing flag",
			wantErrText: "required flag(s) \"name\" not set",
			args:        []string{"repositories", "create"},
		},
		{
			name:            "create repository",
			statusCode:      http.StatusCreated,
			args:            []string{"repositories", "create", "--name", "MyRepo"},
			wantOutputRegex: "^Repository successfully created.[\\s]*$",
		},
		{
			name:            "create repository with description",
			statusCode:      http.StatusCreated,
			args:            []string{"repositories", "create", "--name", "MyRepo", "--description", "My description"},
			wantOutputRegex: "^Repository successfully created.[\\s]*$",
		},
		{
			name:        "repository already exists",
			statusCode:  http.StatusConflict,
			args:        []string{"repositories", "create", "--name", "MyRepo"},
			wantErrText: "409 Conflict: The repository already exists",
		},
	}

	runTests(cases, t)

}
