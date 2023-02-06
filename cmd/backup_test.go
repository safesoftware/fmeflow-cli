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
	f, err := os.CreateTemp("", "*fmeserver-backup.fsconfig")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"backup", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"backup"},
		},
		{
			name:        "422 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"backup"},
		},
		{
			name:             "backup to file",
			statusCode:       http.StatusOK,
			args:             []string{"backup", "--file", f.Name()},
			body:             okResponseV3,
			wantOutputRegex:  "FME Server backed up to",
			wantFileContents: fileContents{file: f.Name(), contents: okResponseV3},
		},
		{
			name:            "backup to shared resource",
			statusCode:      http.StatusAccepted,
			args:            []string{"backup", "--resource"},
			body:            `{"id":4}`,
			wantOutputRegex: "Backup task submitted with id: 4",
			wantFormParams:  map[string]string{"resourceName": "FME_SHAREDRESOURCE_BACKUP"},
		},
		{
			name:            "export to specific file name",
			statusCode:      http.StatusAccepted,
			body:            `{"id":4}`,
			wantOutputRegex: "Backup task submitted with id: 4",
			args:            []string{"backup", "--resource", "--export-package", "TestPackageName.fsconfig"},
			wantFormParams:  map[string]string{"exportPackage": "TestPackageName.fsconfig", "resourceName": "FME_SHAREDRESOURCE_BACKUP"},
		},
		{
			name:            "specify failure topic",
			statusCode:      http.StatusAccepted,
			body:            `{"id":4}`,
			wantOutputRegex: "Backup task submitted with id: 4",
			args:            []string{"backup", "--resource", "--failure-topic", "SAMPLE_TOPIC"},
			wantFormParams:  map[string]string{"failureTopic": "SAMPLE_TOPIC", "resourceName": "FME_SHAREDRESOURCE_BACKUP"},
		},
		{
			name:            "specify success topic",
			statusCode:      http.StatusAccepted,
			body:            `{"id":4}`,
			wantOutputRegex: "Backup task submitted with id: 4",
			args:            []string{"backup", "--resource", "--success-topic", "SAMPLE_TOPIC"},
			wantFormParams:  map[string]string{"successTopic": "SAMPLE_TOPIC", "resourceName": "FME_SHAREDRESOURCE_BACKUP"},
		},
		{
			name:        "don't allow file and resource flags",
			args:        []string{"backup", "--file", f.Name(), "--resource"},
			wantErrText: "if any flags in the group [file resource] are set none of the others can be; [file resource] were all set",
		},
		{
			name:        "don't allow file and resource-name flags",
			args:        []string{"backup", "--file", f.Name(), "--resource-name", "test.fsconfig"},
			wantErrText: "if any flags in the group [file resource-name] are set none of the others can be; [file resource-name] were all set",
		},
		{
			name:        "don't allow file and export-package flags",
			args:        []string{"backup", "--file", f.Name(), "--export-package", "FME_SHAREDRESOURCE_BACKUP"},
			wantErrText: "if any flags in the group [file export-package] are set none of the others can be; [export-package file] were all set",
		},
		{
			name:        "don't allow file and failure-topic flags",
			args:        []string{"backup", "--file", f.Name(), "--failure-topic", "FAILURE_TOPIC"},
			wantErrText: "if any flags in the group [file failure-topic] are set none of the others can be; [failure-topic file] were all set",
		},
		{
			name:        "don't allow file and success-topic flags",
			args:        []string{"backup", "--file", f.Name(), "--success-topic", "SUCCESS_TOPIC"},
			wantErrText: "if any flags in the group [file success-topic] are set none of the others can be; [file success-topic] were all set",
		},
		{
			name:               "missing value for resource name",
			args:               []string{"backup", "--file", f.Name(), "--resource-name"},
			wantErrOutputRegex: "flag needs an argument: --resource-name",
		},
		{
			name:               "missing value for export-package",
			args:               []string{"backup", "--file", f.Name(), "--export-package"},
			wantErrOutputRegex: "flag needs an argument: --export-package",
		},
		{
			name:               "missing value for success topic",
			args:               []string{"backup", "--file", f.Name(), "--success-topic"},
			wantErrOutputRegex: "flag needs an argument: --success-topic",
		},
		{
			name:               "missing value for failure-topic",
			args:               []string{"backup", "--file", f.Name(), "--failure-topic"},
			wantErrOutputRegex: "flag needs an argument: --failure-topic",
		},
	}

	runTests(cases, t)

}
