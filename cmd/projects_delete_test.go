package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectsDelete(t *testing.T) {
	testProjectJson := `{
		"items": [
		  {
			"id": "a64297e7-a119-4e10-ac37-5d0bba12194b",
			"name": "test",
			"hubUid": "",
			"hubPublisherUid": "",
			"description": "test1",
			"readme": "",
			"version": "1.0.0",
			"lastUpdated": "2024-03-26T18:44:30.713Z",
			"owner": "admin",
			"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"shareable": true,
			"lastUpdateUser": "admin",
			"lastUpdateUserID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"hasIcon": false
		  },
		  {
			"id": "5a70afe2-db56-4a15-8ea8-559a10c326ec",
			"name": "test2",
			"hubUid": null,
			"hubPublisherUid": null,
			"description": "test2",
			"readme": null,
			"version": "1.0.0",
			"lastUpdated": "2024-03-26T18:45:32.991Z",
			"owner": "admin",
			"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"shareable": true,
			"lastUpdateUser": "admin",
			"lastUpdateUserID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"hasIcon": false
		  },
		  {
			"id": "0adc89c8-5eda-4cbb-a599-eefd66b1c8b9",
			"name": "testAutomation",
			"hubUid": null,
			"hubPublisherUid": null,
			"description": "test",
			"readme": null,
			"version": "1.0.0",
			"lastUpdated": "2024-03-26T22:52:38.092Z",
			"owner": "admin",
			"ownerID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"shareable": true,
			"lastUpdateUser": "admin",
			"lastUpdateUserID": "fb2dd313-e5cf-432e-a24a-814e46929ab7",
			"hasIcon": false
		  }
		],
		"totalCount": 3,
		"limit": 100,
		"offset": 0
	  }`

	projectDoesNotExistBody := `{
		"message": "Unauthorized request by user admin due to lack of proper permissions or the object does not exist."
	  }`

	customHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(testProjectJson))
			require.NoError(t, err)

		}
		if r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
			_, err := w.Write([]byte(""))
			require.NoError(t, err)

		}

	}

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"projects", "delete", "--name", "test", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"projects", "delete", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--no-prompt"},
		},
		{
			name:        "missing flag",
			wantErrText: "required flag(s) \"id\" or \"name\" not set",
			args:        []string{"projects", "delete"},
		},
		{
			name:            "delete project by name",
			statusCode:      http.StatusNoContent,
			args:            []string{"projects", "delete", "--name", "test", "--no-prompt"},
			wantOutputRegex: "^Project successfully deleted.[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:        "delete project not found",
			statusCode:  http.StatusNotFound,
			args:        []string{"projects", "delete", "--name", "myDep", "--no-prompt"},
			wantErrText: "404 Not Found: check that the specified project exists",
		},
		{
			name:            "delete project by id",
			statusCode:      http.StatusNoContent,
			args:            []string{"projects", "delete", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--no-prompt"},
			wantOutputRegex: "^Project successfully deleted.[\\s]*$",
			wantURLContains: "/fmeapiv4/projects/a64297e7-a119-4e10-ac37-5d0bba12194b",
		},
		{
			name:        "delete project by id not found",
			statusCode:  http.StatusNotFound,
			body:        projectDoesNotExistBody,
			args:        []string{"projects", "delete", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--no-prompt"},
			wantErrText: "Unauthorized request by user admin due to lack of proper permissions or the object does not exist.",
		},
		{
			name:            "delete project by id all content",
			statusCode:      http.StatusNoContent,
			args:            []string{"projects", "delete", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--no-prompt", "--all"},
			wantURLContains: "/fmeapiv4/projects/a64297e7-a119-4e10-ac37-5d0bba12194b/delete-all",
		},
		{
			name:            "delete project by id all content including dependencies",
			statusCode:      http.StatusNoContent,
			args:            []string{"projects", "delete", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--no-prompt", "--all", "--dependencies"},
			wantURLContains: "/fmeapiv4/projects/a64297e7-a119-4e10-ac37-5d0bba12194b/delete-all",
			wantFormParams:  map[string]string{"deleteDependencies": "true"},
		},
	}

	runTests(cases, t)

}
