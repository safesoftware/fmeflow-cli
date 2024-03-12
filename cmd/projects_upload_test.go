package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProjectUpload(t *testing.T) {
	// standard responses for v3 and v4
	responsev3 := `{
		"id": 1
	  }`
	projectContents := "Pretend project file"

	ProjectNotificationStruct := ProjectNotification{
		Type:         "TOPIC",
		SuccessTopic: "MIGRATION_ASYNC_JOB_SUCCESS",
		FailureTopic: "MIGRATION_ASYNC_JOB_FAILURE",
	}

	ProjectGetStruct := ProjectUploadV4{
		JobID:     1,
		Status:    "importing",
		Owner:     "admin",
		OwnerID:   "fb2dd313-e5cf-432e-a24a-814e46929ab7",
		Requested: time.Date(2024, 3, 8, 19, 53, 25, 518000000, time.UTC),
		Generated: time.Date(2024, 3, 8, 19, 53, 25, 518000000, time.UTC),
		FileName:  "test",
		Request: ProjectImportRun{
			FallbackOwnerID:    "",
			Overwrite:          true,
			PauseNotifications: true,
			DisableItems:       false,
			Notification:       &ProjectNotificationStruct,
			SelectedItems:      nil,
		},
	}

	ProjectItemsJson := `{
		"items": [
		  {
			"id": "test",
			"jobId": 68,
			"name": "test",
			"type": "deploymentParameter",
			"ownerId": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"ownerName": "admin",
			"ownerStatus": "id_match",
			"originalOwner": "admin",
			"selected": true,
			"existing": true,
			"previewAction": "skipped",
			"action": "unknown",
			"source": "project"
		  },
		  {
			"id": "6a3ebaf9-e537-4aff-9be0-b8d88542069e",
			"jobId": 68,
			"name": "author",
			"type": "user",
			"ownerId": null,
			"ownerName": null,
			"ownerStatus": "none",
			"originalOwner": null,
			"selected": true,
			"existing": true,
			"previewAction": "overwritten",
			"action": "unknown",
			"source": "project"
		  }
		],
		"totalCount": 2,
		"limit": 100,
		"offset": 0
	  }`

	// generate random file to restore from
	f, err := os.CreateTemp("", "fmeflow-project")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up
	err = os.WriteFile(f.Name(), []byte(projectContents), 0644)
	require.NoError(t, err)

	// state variable for our custom http handler
	getCount := 0

	// this is the generic mock up for this test which will respond similar to how FME Flow should respond
	customHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "upload") {
			// set a location header
			w.Header().Set("Location", "http://localhost:8080/fmeapiv4/migrations/imports/1")
			w.WriteHeader(http.StatusCreated)

			// check if there is a URL argument
			r.ParseForm()
			urlParams := r.Form
			if urlParams.Get("skipPreview") != "" {
				// set status to generating preview
				ProjectGetStruct.Status = "generating_preview"
			} else {
				// set status to ready
				ProjectGetStruct.Status = "ready"
			}

		} else if strings.Contains(r.URL.Path, "run") && r.Method == "POST" {
			w.WriteHeader(http.StatusAccepted)
			require.Contains(t, r.URL.Path, "migrations/imports/1/run")
			// set status to importing
			ProjectGetStruct.Status = "importing"
		} else if strings.Contains(r.URL.Path, "migrations/imports/1/items") && r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			// return the test list of selectable items
			_, err := w.Write([]byte(ProjectItemsJson))
			require.NoError(t, err)
			// set status to importing
			ProjectGetStruct.Status = "importing"
		} else if strings.Contains(r.URL.Path, "migrations/imports/1") && r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			if getCount < 1 {
				getCount++
			} else {
				if ProjectGetStruct.Status == "importing" {
					// set status to imported
					ProjectGetStruct.Status = "imported"
				} else if ProjectGetStruct.Status == "generating_preview" {
					//set status to ready
					ProjectGetStruct.Status = "ready"
				}
				getCount = 0
			}
			// marshal the struct to json
			projectGetJson, err := json.Marshal(ProjectGetStruct)
			require.NoError(t, err)

			// write the json to the response
			_, err = w.Write(projectGetJson)
			require.NoError(t, err)
		} else if strings.Contains(r.URL.Path, "migrations/imports/1") && r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	// this will check for the correct settings for the quick test, then call the custom handler
	quickHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "upload") {
			// check if the quick flag is set on upload
			r.ParseForm()
			urlParams := r.Form
			require.True(t, urlParams.Has("skipPreview"))
		}
		customHttpServerHandler(w, r)
	}

	// this will check for the correct settings for the overwrite test, then call the custom handler
	overwriteHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "run") && r.Method == "POST" {
			// check if the overwite is in the json body
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, true, result["overwrite"])
		}
		customHttpServerHandler(w, r)
	}

	// this will check for the correct settings for the overwrite test, then call the custom handler
	pauseNotificationsHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "run") && r.Method == "POST" {
			// check if the overwite is in the json body
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, true, result["pauseNotifications"])
		}
		customHttpServerHandler(w, r)
	}

	// this will check for the correct settings for the disableItems test, then call the custom handler
	disableItemsHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "run") && r.Method == "POST" {
			// check if the overwite is in the json body
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, true, result["disableItems"])
		}
		customHttpServerHandler(w, r)
	}

	// this will check for the correct settings for the topics test, then call the custom handler
	topicsHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "run") && r.Method == "POST" {
			// check if the overwite is in the json body
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result ProjectImportRun
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, "SUCCESS_TOPIC", result.Notification.SuccessTopic)
			require.Equal(t, "FAILURE_TOPIC", result.Notification.FailureTopic)
		}
		customHttpServerHandler(w, r)
	}

	// this will check that the selected items are set correctly
	selectedItemsAllHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "run") && r.Method == "POST" {
			// check if the overwite is in the json body
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result ProjectImportRun
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, 2, len(result.SelectedItems))
			require.Equal(t, "test", result.SelectedItems[0].ID)
			require.Equal(t, "deploymentParameter", result.SelectedItems[0].Type)
			require.Equal(t, "6a3ebaf9-e537-4aff-9be0-b8d88542069e", result.SelectedItems[1].ID)
			require.Equal(t, "user", result.SelectedItems[1].Type)
		}
		customHttpServerHandler(w, r)
	}

	// this will check that only the test item is set
	selectedItemsListHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "run") && r.Method == "POST" {
			// check if the overwite is in the json body
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result ProjectImportRun
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, 1, len(result.SelectedItems))
			require.Equal(t, "test", result.SelectedItems[0].ID)
			require.Equal(t, "deploymentParameter", result.SelectedItems[0].Type)
		}
		customHttpServerHandler(w, r)
	}

	// check that no selected items are set
	selectedItemsNoneHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "run") && r.Method == "POST" {
			// check if the overwite is in the json body
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result ProjectImportRun
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, 0, len(result.SelectedItems))
		}
		customHttpServerHandler(w, r)
	}

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
			args:        []string{"projects", "upload", "--file", f.Name(), "--overwrite", "--projects-import-mode", "UPDATE", "--api-version", "v3"},
		},
		{
			name:        "duplicate flags overwrite import-mode",
			wantErrText: "if any flags in the group [overwrite import-mode] are set none of the others can be; [import-mode overwrite] were all set",
			args:        []string{"projects", "upload", "--file", f.Name(), "--overwrite", "--import-mode", "UPDATE", "--api-version", "v3"},
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
			name:            "upload project V4 quick",
			args:            []string{"projects", "upload", "--file", f.Name(), "--quick"},
			httpServer:      httptest.NewServer(http.HandlerFunc(quickHttpServerHandler)),
			wantOutputRegex: "Project import complete.",
		},
		{
			name:            "upload project V4 overwrite",
			args:            []string{"projects", "upload", "--file", f.Name(), "--overwrite"},
			httpServer:      httptest.NewServer(http.HandlerFunc(overwriteHttpServerHandler)),
			wantOutputRegex: "Project import complete.",
		},
		{
			name:            "upload project V4 pause-notifications",
			args:            []string{"projects", "upload", "--file", f.Name(), "--pause-notifications"},
			httpServer:      httptest.NewServer(http.HandlerFunc(pauseNotificationsHttpServerHandler)),
			wantOutputRegex: "Project import complete.",
		},
		{
			name:            "upload project V4 disable-project-items",
			args:            []string{"projects", "upload", "--file", f.Name(), "--disable-project-items"},
			httpServer:      httptest.NewServer(http.HandlerFunc(disableItemsHttpServerHandler)),
			wantOutputRegex: "Project import complete.",
		},
		{
			name:            "upload project V4 get-selectable",
			args:            []string{"projects", "upload", "--file", f.Name(), "--get-selectable"},
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
			wantOutputRegex: "^[\\s]*ID[\\s]*TYPE[\\s]*test[\\s]*deploymentParameter[\\s]*6a3ebaf9-e537-4aff-9be0-b8d88542069e[\\s]*user[\\s]*$",
		},
		{
			name:           "upload project V4 get-selectable json",
			args:           []string{"projects", "upload", "--file", f.Name(), "--get-selectable", "--json"},
			httpServer:     httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
			wantOutputJson: ProjectItemsJson,
		},
		{
			name:       "upload project V4 set success and failure topics",
			args:       []string{"projects", "upload", "--file", f.Name(), "--success-topic", "SUCCESS_TOPIC", "--failure-topic", "FAILURE_TOPIC"},
			httpServer: httptest.NewServer(http.HandlerFunc(topicsHttpServerHandler)),
		},
		{
			name:       "upload project V4 selected items all",
			args:       []string{"projects", "upload", "--file", f.Name(), "--selected-items", "all"},
			httpServer: httptest.NewServer(http.HandlerFunc(selectedItemsAllHttpServerHandler)),
		},
		{
			name:       "upload project V4 selected items none",
			args:       []string{"projects", "upload", "--file", f.Name(), "--selected-items", "none"},
			httpServer: httptest.NewServer(http.HandlerFunc(selectedItemsNoneHttpServerHandler)),
		},
		{
			name:       "upload project V4 selected items list",
			args:       []string{"projects", "upload", "--file", f.Name(), "--selected-items", "test:deploymentParameter"},
			httpServer: httptest.NewServer(http.HandlerFunc(selectedItemsListHttpServerHandler)),
		},
		{
			name:       "upload project V4 selected items list 2",
			args:       []string{"projects", "upload", "--file", f.Name(), "--selected-items", "test:deploymentParameter,6a3ebaf9-e537-4aff-9be0-b8d88542069e:user"},
			httpServer: httptest.NewServer(http.HandlerFunc(selectedItemsAllHttpServerHandler)),
		},
		{
			name:        "upload project V4 invalid selected items syntax",
			args:        []string{"projects", "upload", "--file", f.Name(), "--selected-items", "test:deploymentParameter,6a3ebaf9-e537-4aff-9be0-b8d88542069e"},
			wantErrText: "invalid selected items. Must be a comma separated list of item ids and types. e.g. item1:itemtype1,item2:itemtype2",
			httpServer:  httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:        "upload project V4 invalid selected items in package",
			args:        []string{"projects", "upload", "--file", f.Name(), "--selected-items", "test:deploymentParameter,author:user"},
			wantErrText: "selected item author (user) is not in the list of selectable items",
			httpServer:  httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
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
