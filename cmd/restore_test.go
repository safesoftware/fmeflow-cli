package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRestore(t *testing.T) {
	// standard responses for v3 and v4
	response := `{
		"id": 1
	  }`
	backupContents := "Pretend backup file"

	// generate random file to restore from
	f, err := os.CreateTemp("", "fmeflow-backup")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up
	err = os.WriteFile(f.Name(), []byte(backupContents), 0644)
	require.NoError(t, err)

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"restore", "--file", f.Name(), "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error: check that the file specified is a valid backup file",
			args:        []string{"restore", "--file", f.Name()},
		},
		{
			name:        "422 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"restore", "--file", f.Name()},
		},
		{
			name:        "missing required flags",
			wantErrText: "required flag \"file\" or \"resource\" not set",
			args:        []string{"restore"},
		},
		{
			name:            "restore from file",
			statusCode:      http.StatusOK,
			args:            []string{"restore", "--file", f.Name()},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
		},
		{
			name:            "restore from resource",
			statusCode:      http.StatusAccepted,
			args:            []string{"restore", "--resource"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
		},
		{
			name:            "restore from resource specific file",
			statusCode:      http.StatusAccepted,
			args:            []string{"restore", "--resource", "--file", "ServerConfigPackage.fsconfig"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			wantFormParams:  map[string]string{"importPackage": "ServerConfigPackage.fsconfig"},
		},
		{
			name:            "restore from resource specific file failure and success topics",
			statusCode:      http.StatusAccepted,
			args:            []string{"restore", "--resource", "--file", "ServerConfigPackage.fsconfig", "--success-topic", "SUCCESS", "--failure-topic", "FAILURE"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			wantFormParams:  map[string]string{"importPackage": "ServerConfigPackage.fsconfig", "successTopic": "SUCCESS", "failureTopic": "FAILURE"},
		},
		{
			name:            "restore from resource specific file and specific shared resource",
			statusCode:      http.StatusAccepted,
			args:            []string{"restore", "--resource", "--file", "ServerConfigPackage.fsconfig", "--resource-name", "OTHER_RESOURCE"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			wantFormParams:  map[string]string{"importPackage": "ServerConfigPackage.fsconfig", "resourceName": "OTHER_RESOURCE"},
		},
		{
			name:            "import mode",
			statusCode:      http.StatusOK,
			args:            []string{"restore", "--file", f.Name(), "--import-mode", "UPDATE"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			wantFormParams:  map[string]string{"importMode": "UPDATE"},
		},
		{
			name:            "projects import mode",
			statusCode:      http.StatusOK,
			args:            []string{"restore", "--file", f.Name(), "--projects-import-mode", "UPDATE"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			wantFormParams:  map[string]string{"projectsImportMode": "UPDATE"},
			wantBodyRegEx:   backupContents,
		},
		{
			name:            "pause-notifications",
			statusCode:      http.StatusOK,
			args:            []string{"restore", "--file", f.Name(), "--pause-notifications"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			wantFormParams:  map[string]string{"pauseNotifications": "true"},
		},
	}

	runTests(cases, t)

}
