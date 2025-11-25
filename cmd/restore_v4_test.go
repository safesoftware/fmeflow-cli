package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRestoreV4(t *testing.T) {
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
			fmeflowBuild:       26000,
		},
		{
			name:         "500 bad status code",
			statusCode:   http.StatusInternalServerError,
			wantErrText:  "500 Internal Server Error",
			args:         []string{"restore", "--file", f.Name()},
			fmeflowBuild: 26000,
		},
		{
			name:         "422 bad status code",
			statusCode:   http.StatusNotFound,
			wantErrText:  "404 Not Found",
			args:         []string{"restore", "--file", f.Name()},
			fmeflowBuild: 26000,
		},
		{
			name:         "missing required flags",
			wantErrText:  "required flag \"file\" or \"resource\" not set",
			args:         []string{"restore"},
			fmeflowBuild: 26000,
		},
		{
			name:            "resource without file",
			statusCode:      http.StatusAccepted,
			args:            []string{"restore", "--resource"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			fmeflowBuild:    26000,
		},
		{
			name:            "restore from file",
			statusCode:      http.StatusOK,
			args:            []string{"restore", "--file", f.Name()},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			fmeflowBuild:    26000,
		},
		{
			name:            "restore from resource",
			statusCode:      http.StatusAccepted,
			args:            []string{"restore", "--resource", "--file", "ServerConfigPackage.fsconfig"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			fmeflowBuild:    26000,
		},
		{
			name:            "restore from resource specific file",
			statusCode:      http.StatusAccepted,
			args:            []string{"restore", "--resource", "--file", "ServerConfigPackage.fsconfig"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			fmeflowBuild:    26000,
		},
		{
			name:            "restore from resource specific file failure and success topics",
			statusCode:      http.StatusAccepted,
			args:            []string{"restore", "--resource", "--file", "ServerConfigPackage.fsconfig", "--success-topic", "SUCCESS", "--failure-topic", "FAILURE"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			wantBodyRegEx:   `.*"successTopic":"SUCCESS".*"failureTopic":"FAILURE".*`,
			fmeflowBuild:    26000,
		},
		{
			name:            "restore from resource specific file and specific shared resource",
			statusCode:      http.StatusAccepted,
			args:            []string{"restore", "--resource", "--file", "ServerConfigPackage.fsconfig", "--resource-name", "OTHER_RESOURCE"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			wantBodyRegEx:   `.*"resourceName":"OTHER_RESOURCE".*"packagePath":"ServerConfigPackage.fsconfig".*`,
			fmeflowBuild:    26000,
		},
		{
			name:            "pause-notifications false",
			statusCode:      http.StatusOK,
			args:            []string{"restore", "--file", f.Name(), "--pause-notifications=false"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			wantBodyRegEx:   `.*"pauseNotifications":false.*`,
			fmeflowBuild:    26000,
		},
		{
			name:            "overwrite true",
			statusCode:      http.StatusOK,
			args:            []string{"restore", "--file", f.Name(), "--overwrite"},
			body:            response,
			wantOutputRegex: "Restore task submitted with id: 1",
			wantBodyRegEx:   `.*"overwrite":true.*`,
			fmeflowBuild:    26000,
		},
		{
			name:            "json output",
			statusCode:      http.StatusOK,
			args:            []string{"restore", "--file", f.Name(), "--json"},
			body:            response,
			wantOutputRegex: `"id": 1`,
			fmeflowBuild:    26000,
		},
	}

	runTests(cases, t)
}
