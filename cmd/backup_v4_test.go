package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBackupV4(t *testing.T) {
	// standard responses for v3 and v4
	okResponseV4 := `Random file contents`

	// generate random file to back up to
	f, err := os.CreateTemp("", "*fmeflow-backup.fsconfig")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"backup", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
			fmeflowBuild:       26000,
		},
		{
			name:         "500 bad status code",
			statusCode:   http.StatusInternalServerError,
			wantErrText:  "500 Internal Server Error",
			args:         []string{"backup"},
			fmeflowBuild: 26000,
		},
		{
			name:         "422 bad status code",
			statusCode:   http.StatusNotFound,
			wantErrText:  "404 Not Found",
			args:         []string{"backup"},
			fmeflowBuild: 26000,
		},
		{
			name:             "backup to file",
			statusCode:       http.StatusOK,
			args:             []string{"backup", "--file", f.Name()},
			body:             okResponseV4,
			wantOutputRegex:  "FME Server backed up to",
			wantFileContents: fileContents{file: f.Name(), contents: okResponseV4},
			fmeflowBuild:     26000,
		},
		{
			name:            "backup to shared resource",
			statusCode:      http.StatusAccepted,
			args:            []string{"backup", "--resource"},
			body:            `{"id":4}`,
			wantOutputRegex: "Backup task submitted with id: 4",
			wantBodyRegEx:   `"resourceName":"FME_SHAREDRESOURCE_BACKUP"`,
			fmeflowBuild:    26000,
		},
		{
			name:            "export to specific file name",
			statusCode:      http.StatusAccepted,
			body:            `{"id":4}`,
			wantOutputRegex: "Backup task submitted with id: 4",
			args:            []string{"backup", "--resource", "--export-package", "TestPackageName.fsconfig"},
			wantBodyRegEx:   `"resourceName":"FME_SHAREDRESOURCE_BACKUP".*"packagePath":"TestPackageName.fsconfig"`,
			fmeflowBuild:    26000,
		},
		{
			name:            "specify failure topic",
			statusCode:      http.StatusAccepted,
			body:            `{"id":4}`,
			wantOutputRegex: "Backup task submitted with id: 4",
			args:            []string{"backup", "--resource", "--failure-topic", "SAMPLE_TOPIC"},
			wantBodyRegEx:   `"resourceName":"FME_SHAREDRESOURCE_BACKUP".*"failureTopic":"SAMPLE_TOPIC"`,
			fmeflowBuild:    26000,
		},
		{
			name:            "specify success topic",
			statusCode:      http.StatusAccepted,
			body:            `{"id":4}`,
			wantOutputRegex: "Backup task submitted with id: 4",
			args:            []string{"backup", "--resource", "--success-topic", "SAMPLE_TOPIC"},
			wantBodyRegEx:   `"resourceName":"FME_SHAREDRESOURCE_BACKUP".*"successTopic":"SAMPLE_TOPIC"`,
			fmeflowBuild:    26000,
		},
		{
			name:         "don't allow file and resource flags",
			args:         []string{"backup", "--file", f.Name(), "--resource"},
			wantErrText:  "if any flags in the group [file resource] are set none of the others can be; [file resource] were all set",
			fmeflowBuild: 26000,
		},
		{
			name:         "don't allow file and resource-name flags",
			args:         []string{"backup", "--file", f.Name(), "--resource-name", "test.fsconfig"},
			wantErrText:  "if any flags in the group [file resource-name] are set none of the others can be; [file resource-name] were all set",
			fmeflowBuild: 26000,
		},
		{
			name:         "don't allow file and export-package flags",
			args:         []string{"backup", "--file", f.Name(), "--export-package", "FME_SHAREDRESOURCE_BACKUP"},
			wantErrText:  "if any flags in the group [file export-package] are set none of the others can be; [export-package file] were all set",
			fmeflowBuild: 26000,
		},
		{
			name:         "don't allow file and failure-topic flags",
			args:         []string{"backup", "--file", f.Name(), "--failure-topic", "FAILURE_TOPIC"},
			wantErrText:  "if any flags in the group [file failure-topic] are set none of the others can be; [failure-topic file] were all set",
			fmeflowBuild: 26000,
		},
		{
			name:         "don't allow file and success-topic flags",
			args:         []string{"backup", "--file", f.Name(), "--success-topic", "SUCCESS_TOPIC"},
			wantErrText:  "if any flags in the group [file success-topic] are set none of the others can be; [file success-topic] were all set",
			fmeflowBuild: 26000,
		},
		{
			name:               "missing value for resource name",
			args:               []string{"backup", "--file", f.Name(), "--resource-name"},
			wantErrOutputRegex: "flag needs an argument: --resource-name",
			fmeflowBuild:       26000,
		},
		{
			name:               "missing value for export-package",
			args:               []string{"backup", "--file", f.Name(), "--export-package"},
			wantErrOutputRegex: "flag needs an argument: --export-package",
			fmeflowBuild:       26000,
		},
		{
			name:               "missing value for success topic",
			args:               []string{"backup", "--file", f.Name(), "--success-topic"},
			wantErrOutputRegex: "flag needs an argument: --success-topic",
			fmeflowBuild:       26000,
		},
		{
			name:               "missing value for failure-topic",
			args:               []string{"backup", "--file", f.Name(), "--failure-topic"},
			wantErrOutputRegex: "flag needs an argument: --failure-topic",
			fmeflowBuild:       26000,
		},
	}

	runTests(cases, t)

}
