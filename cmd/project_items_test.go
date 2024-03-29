package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectsItems(t *testing.T) {
	itemsResponse := `{
		"items": [
		  {
			"id": "6a3ebaf9-e537-4aff-9be0-b8d88542069e",
			"name": "author",
			"type": "user",
			"owner": null,
			"lastUpdated": null,
			"dependencies": null
		  },
		  {
			"id": "TestRepo1/TestWorkspace1.fmw",
			"name": "TestWorkspace1.fmw",
			"type": "workspace",
			"owner": "admin",
			"lastUpdated": "2024-03-26T18:41:27.471Z",
			"dependencies": null
		  }
		],
		"totalCount": 2,
		"limit": 100,
		"offset": 0
	  }`

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
		  }
		],
		"totalCount": 1,
		"limit": 100,
		"offset": 0
	}`

	customHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {

		// send the file if we are downloading
		if strings.Contains(r.URL.Path, "items") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(itemsResponse))
			require.NoError(t, err)
		} else {
			// otherwise we are getting the project by name
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(testProjectJson))
			require.NoError(t, err)
		}

	}

	cases := []testCase{
		{
			name:            "Get all items for a project via id",
			args:            []string{"projects", "items", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b"},
			statusCode:      http.StatusOK,
			body:            itemsResponse,
			wantOutputRegex: "ID[\\s]*NAME[\\s]*TYPE[\\s]*OWNER[\\s]*LAST[\\s]*UPDATED[\\s]*[\\s]*6a3ebaf9-e537-4aff-9be0-b8d88542069e[\\s]*author[\\s]*user[\\s]*[\\s]*0001-01-01[\\s]*00:00:00[\\s]*\\+[\\s]*0000[\\s]*UTC[\\s]*[\\s]*TestRepo1/TestWorkspace1.fmw[\\s]*TestWorkspace1.fmw[\\s]*workspace[\\s]*admin[\\s]*2024-03-26[\\s]*18:41:27.471[\\s]*\\+[\\s]*0000[\\s]*UTC",
		},
		{
			name:            "Get all items for a project via name",
			args:            []string{"projects", "items", "--name", "test"},
			wantOutputRegex: "ID[\\s]*NAME[\\s]*TYPE[\\s]*OWNER[\\s]*LAST[\\s]*UPDATED[\\s]*[\\s]*6a3ebaf9-e537-4aff-9be0-b8d88542069e[\\s]*author[\\s]*user[\\s]*[\\s]*0001-01-01[\\s]*00:00:00[\\s]*\\+[\\s]*0000[\\s]*UTC[\\s]*[\\s]*TestRepo1/TestWorkspace1.fmw[\\s]*TestWorkspace1.fmw[\\s]*workspace[\\s]*admin[\\s]*2024-03-26[\\s]*18:41:27.471[\\s]*\\+[\\s]*0000[\\s]*UTC",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:               "unknown flag",
			args:               []string{"projects", "items", "--name", "test", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"projects", "items", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b"},
		},
		{
			name:        "missing flag",
			wantErrText: "required flag(s) \"id\" or \"name\" not set",
			args:        []string{"projects", "items"},
		},
		{
			name:           "Get items including dependencies",
			args:           []string{"projects", "items", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--include-dependencies"},
			statusCode:     http.StatusOK,
			body:           itemsResponse,
			wantFormParams: map[string]string{"includeDependencies": "true"},
		},
		{
			name:           "Get items of type workspace",
			args:           []string{"projects", "items", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--type", "workspace"},
			statusCode:     http.StatusOK,
			body:           itemsResponse,
			wantFormParams: map[string]string{"type": "workspace"},
		},
		{
			name:           "Get items with filter string",
			args:           []string{"projects", "items", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--filter-string", "TestWorkspace1.fmw"},
			statusCode:     http.StatusOK,
			body:           itemsResponse,
			wantFormParams: map[string]string{"filterString": "TestWorkspace1.fmw"},
		},
		{
			name:           "Get items with filter property",
			args:           []string{"projects", "items", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--filter-property", "name", "--filter-string", "TestWorkspace1.fmw"},
			statusCode:     http.StatusOK,
			body:           itemsResponse,
			wantFormParams: map[string]string{"filterProperties": "name", "filterString": "TestWorkspace1.fmw"},
		},
		{
			name:        "Not allowed filter property without filter string",
			args:        []string{"projects", "items", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--filter-property", "name"},
			wantErrText: "flag \"filter-property\" specified without flag \"filter-string\"",
		},
		{
			name:            "Get items custom columns output",
			args:            []string{"projects", "items", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--output=custom-columns=ID:.id,NAME:.name"},
			statusCode:      http.StatusOK,
			body:            itemsResponse,
			wantOutputRegex: "ID[\\s]*NAME[\\s]*[\\s]*6a3ebaf9-e537-4aff-9be0-b8d88542069e[\\s]*author[\\s]*[\\s]*TestRepo1/TestWorkspace1.fmw[\\s]*TestWorkspace1.fmw[\\s]*",
		},
		{
			name:           "Get items json output",
			args:           []string{"projects", "items", "--id", "a64297e7-a119-4e10-ac37-5d0bba12194b", "--output=json"},
			statusCode:     http.StatusOK,
			body:           itemsResponse,
			wantOutputJson: itemsResponse,
		},
	}

	runTests(cases, t)

}
