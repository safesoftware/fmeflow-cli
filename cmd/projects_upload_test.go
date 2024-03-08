package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectUpload(t *testing.T) {
	// standard responses for v3 and v4
	responsev3 := `{
		"id": 1
	  }`
	projectContents := "Pretend project file"

	// generate random file to restore from
	f, err := os.CreateTemp("", "fmeflow-project")
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
			wantErrText: "500 Internal Server Error: check that the file specified is a valid project file",
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
			name:        "duplicate flags overwrite projects-import-mode",
			wantErrText: "if any flags in the group [overwrite projects-import-mode] are set none of the others can be; [overwrite projects-import-mode] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--overwrite", "--projects-import-mode", "UPDATE"},
		},
		{
			name:        "duplicate flags overwrite import-mode",
			wantErrText: "if any flags in the group [overwrite import-mode] are set none of the others can be; [import-mode overwrite] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--overwrite", "--import-mode", "UPDATE"},
		},
		{
			name:        "duplicate flags quick get-selectable",
			wantErrText: "if any flags in the group [quick get-selectable] are set none of the others can be; [get-selectable quick] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--quick", "--get-selectable"},
		},
		{
			name:        "duplicate flags interactive get-selectable",
			wantErrText: "if any flags in the group [interactive get-selectable] are set none of the others can be; [get-selectable interactive] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--interactive", "--get-selectable"},
		},
		{
			name:        "duplicate flags overwrite get-selectable",
			wantErrText: "if any flags in the group [get-selectable overwrite] are set none of the others can be; [get-selectable overwrite] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--overwrite", "--get-selectable"},
		},
		{
			name:        "duplicate flags pause-notifications get-selectable",
			wantErrText: "if any flags in the group [get-selectable pause-notifications] are set none of the others can be; [get-selectable pause-notifications] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--pause-notifications", "--get-selectable"},
		},
		{
			name:        "duplicate flags selected-items interactive",
			wantErrText: "if any flags in the group [selected-items interactive] are set none of the others can be; [interactive selected-items] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--selected-items=[]", "--interactive"},
		},
		{
			name:        "duplicate flags selected-items get-selectable",
			wantErrText: "if any flags in the group [selected-items get-selectable] are set none of the others can be; [get-selectable selected-items] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--selected-items=[]", "--get-selectable"},
		},
		{
			name:        "duplicate flags selected-items quick",
			wantErrText: "if any flags in the group [selected-items quick] are set none of the others can be; [quick selected-items] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--selected-items=[]", "--quick"},
		},
		{
			name:        "duplicate flags interactive quick",
			wantErrText: "if any flags in the group [quick interactive] are set none of the others can be; [interactive quick] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--interactive", "--quick"},
		},
		{
			name:            "upload project V3",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "upload", "--file", f.Name(), "--api-version", "v3"},
			body:            responsev3,
			wantOutputRegex: "Project Upload task submitted with id: 1",
		},
		{
			name:            "import mode V3",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "upload", "--file", f.Name(), "--import-mode", "UPDATE", "--api-version", "v3"},
			body:            responsev3,
			wantOutputRegex: "Project Upload task submitted with id: 1",
			wantFormParams:  map[string]string{"importMode": "UPDATE"},
		},
		{
			name:            "projects import mode V3",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "upload", "--file", f.Name(), "--projects-import-mode", "UPDATE", "--api-version", "v3"},
			body:            responsev3,
			wantOutputRegex: "Project Upload task submitted with id: 1",
			wantFormParams:  map[string]string{"projectsImportMode": "UPDATE"},
			wantBodyRegEx:   projectContents,
		},
		{
			name:            "pause-notifications V3",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "upload", "--file", f.Name(), "--pause-notifications", "--api-version", "v3"},
			body:            responsev3,
			wantOutputRegex: "Project Upload task submitted with id: 1",
			wantFormParams:  map[string]string{"pauseNotifications": "true"},
		},
		{
			name:            "disable project items V3",
			statusCode:      http.StatusOK,
			args:            []string{"projects", "upload", "--file", f.Name(), "--disable-project-items", "--api-version", "v3"},
			body:            responsev3,
			wantOutputRegex: "Project Upload task submitted with id: 1",
			wantFormParams:  map[string]string{"disableProjectItems": "true"},
		},
	}

	runTests(cases, t)

}
