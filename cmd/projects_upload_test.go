package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectUpload(t *testing.T) {
	// standard responses for v3 and v4
	response := `{
		"id": 1
	  }`
	projectContents := "Pretend project file"

	// generate random file to restore from
	f, err := os.CreateTemp("", "fmeserver-project")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up
	err = os.WriteFile(f.Name(), []byte(projectContents), 0644)
	require.NoError(t, err)

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"projects", "upload", "--file", f.Name(), "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"projects", "upload", "--file", f.Name()},
		},
		{
			name:        "422 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"projects", "upload", "--file", f.Name()},
		},
		{
			name:        "missing required flags",
			wantErrText: "required flag(s) \"file\" not set",
			args:        []string{"projects", "upload"},
		},
		{
			name:            "upload project",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "upload", "--file", f.Name()},
			body:            response,
			wantOutputRegex: "Project Upload task submitted with id: 1",
		},
		{
			name:            "import mode",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "upload", "--file", f.Name(), "--import-mode", "UPDATE"},
			body:            response,
			wantOutputRegex: "Project Upload task submitted with id: 1",
			wantFormParams:  map[string]string{"importMode": "UPDATE"},
		},
		{
			name:            "projects import mode",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "upload", "--file", f.Name(), "--projects-import-mode", "UPDATE"},
			body:            response,
			wantOutputRegex: "Project Upload task submitted with id: 1",
			wantFormParams:  map[string]string{"projectsImportMode": "UPDATE"},
			wantBodyRegEx:   projectContents,
		},
		{
			name:            "pause-notifications",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "upload", "--file", f.Name(), "--pause-notifications"},
			body:            response,
			wantOutputRegex: "Project Upload task submitted with id: 1",
			wantFormParams:  map[string]string{"pauseNotifications": "true"},
		},
		{
			name:            "disable project items",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "upload", "--file", f.Name(), "--disable-project-items"},
			body:            response,
			wantOutputRegex: "Project Upload task submitted with id: 1",
			wantFormParams:  map[string]string{"disableProjectItems": "true"},
		},
	}

	runTests(cases, t)

}
