package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBackup(t *testing.T) {
	// standard responses for v3 and v4
	okResponseV3 := `Random file contents`

	// generate random file to back up to
	f, err := os.CreateTemp("", "fmeserver-backup")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up

	cases := []testCase{
		{
			name:        "unknown flag",
			statusCode:  http.StatusOK,
			args:        []string{"--badflag"},
			wantErrText: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
		},
		{
			name:        "422 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
		},
		{
			name:       "backup to file",
			statusCode: http.StatusOK,
			args:       []string{"--file", f.Name()},
			body:       okResponseV3,
			wantOutput: "FME Server backed up to",
		},
		{
			name:           "backup to shared resource",
			statusCode:     http.StatusAccepted,
			args:           []string{"--resource"},
			body:           `{"id":4}`,
			wantOutput:     "Backup task submitted with id: 4",
			wantFormParams: map[string]string{"resourceName": "FME_SHAREDRESOURCE_BACKUP"},
		},
		{
			name:           "export to specific file name",
			statusCode:     http.StatusAccepted,
			body:           `{"id":4}`,
			wantOutput:     "Backup task submitted with id: 4",
			args:           []string{"--resource", "--export-package", "TestPackageName.fsconfig"},
			wantFormParams: map[string]string{"exportPackage": "TestPackageName.fsconfig", "resourceName": "FME_SHAREDRESOURCE_BACKUP"},
		},
		{
			name:           "specify failure topic",
			statusCode:     http.StatusAccepted,
			body:           `{"id":4}`,
			wantOutput:     "Backup task submitted with id: 4",
			args:           []string{"--resource", "--failure-topic", "SAMPLE_TOPIC"},
			wantFormParams: map[string]string{"failureTopic": "SAMPLE_TOPIC", "resourceName": "FME_SHAREDRESOURCE_BACKUP"},
		},
		{
			name:           "specify success topic",
			statusCode:     http.StatusAccepted,
			body:           `{"id":4}`,
			wantOutput:     "Backup task submitted with id: 4",
			args:           []string{"--resource", "--success-topic", "SAMPLE_TOPIC"},
			wantFormParams: map[string]string{"successTopic": "SAMPLE_TOPIC", "resourceName": "FME_SHAREDRESOURCE_BACKUP"},
		},
		{
			name:        "don't allow file and resource flags",
			args:        []string{"--file", f.Name(), "--resource"},
			wantErrText: "if any flags in the group [file resource] are set none of the others can be; [file resource] were all set",
		},
		{
			name:        "don't allow file and resource-name flags",
			args:        []string{"--file", f.Name(), "--resource-name", "test.fsconfig"},
			wantErrText: "if any flags in the group [file resource-name] are set none of the others can be; [file resource-name] were all set",
		},
		{
			name:        "don't allow file and export-package flags",
			args:        []string{"--file", f.Name(), "--export-package", "FME_SHAREDRESOURCE_BACKUP"},
			wantErrText: "if any flags in the group [file export-package] are set none of the others can be; [export-package file] were all set",
		},
		{
			name:        "don't allow file and failure-topic flags",
			args:        []string{"--file", f.Name(), "--failure-topic", "FAILURE_TOPIC"},
			wantErrText: "if any flags in the group [file failure-topic] are set none of the others can be; [failure-topic file] were all set",
		},
		{
			name:        "don't allow file and success-topic flags",
			args:        []string{"--file", f.Name(), "--success-topic", "SUCCESS_TOPIC"},
			wantErrText: "if any flags in the group [file success-topic] are set none of the others can be; [file success-topic] were all set",
		},
		{
			name:        "missing value for resource name",
			args:        []string{"--file", f.Name(), "--resource-name"},
			wantErrText: "flag needs an argument: --resource-name",
		},
		{
			name:        "missing value for export-package",
			args:        []string{"--file", f.Name(), "--export-package"},
			wantErrText: "flag needs an argument: --export-package",
		},
		{
			name:        "missing value for success topic",
			args:        []string{"--file", f.Name(), "--success-topic"},
			wantErrText: "flag needs an argument: --success-topic",
		},
		{
			name:        "missing value for failure-topic",
			args:        []string{"--file", f.Name(), "--failure-topic"},
			wantErrText: "flag needs an argument: --failure-topic",
		},
	}

	runTests(cases, newBackupCmd, t)

}
