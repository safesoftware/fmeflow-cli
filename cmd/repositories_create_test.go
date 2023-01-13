package cmd

import (
	"net/http"
	"testing"
)

func TestRepositoriesCreate(t *testing.T) {
	repoExistsBodyV4 := `{
		"message": "MyRepo is not a valid name. A repository named Samples already exists."
	  }`
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
			name:            "create repository V4",
			statusCode:      http.StatusCreated,
			args:            []string{"repositories", "create", "--name", "MyRepo", "--api-version", "v4"},
			wantOutputRegex: "^Repository successfully created.[\\s]*$",
		},
		{
			name:            "create repository with description V4",
			statusCode:      http.StatusCreated,
			args:            []string{"repositories", "create", "--name", "MyRepo", "--description", "My description", "--api-version", "v4"},
			wantOutputRegex: "^Repository successfully created.[\\s]*$",
		},
		{
			name:        "repository already exists V4",
			statusCode:  http.StatusConflict,
			body:        repoExistsBodyV4,
			args:        []string{"repositories", "create", "--name", "MyRepo", "--api-version", "v4"},
			wantErrText: "MyRepo is not a valid name. A repository named Samples already exists.",
		},
		{
			name:            "create repository V3",
			statusCode:      http.StatusCreated,
			args:            []string{"repositories", "create", "--name", "MyRepo", "--api-version", "v3"},
			wantOutputRegex: "^Repository successfully created.[\\s]*$",
		},
		{
			name:            "create repository with description V3",
			statusCode:      http.StatusCreated,
			args:            []string{"repositories", "create", "--name", "MyRepo", "--description", "My description", "--api-version", "v3"},
			wantOutputRegex: "^Repository successfully created.[\\s]*$",
		},
		{
			name:        "repository already exists V3",
			statusCode:  http.StatusConflict,
			args:        []string{"repositories", "create", "--name", "MyRepo", "--api-version", "v3"},
			wantErrText: "409 Conflict: The repository already exists",
		},
	}

	runTests(cases, t)

}
