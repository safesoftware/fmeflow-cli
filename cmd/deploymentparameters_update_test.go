package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeploymentParametersUpdate(t *testing.T) {
	parameterDoesNotExistBody := `{
		"message": "Unauthorized request by user admin due to lack of proper permissions or the object does not exist."
}`

	responseSingleText := `{
		"type": "text",
		"name": "testText",
		"value": "textValue",
		"type": "text",
		"owner": "admin",
		"updated": "2024-03-15T17:12:21.407Z",
		"resourceMissing": false
	  }`

	responseSingleDb := `{
		"type": "dropdown",
		"name": "testDb",
		"value": "db",
		"type": "dropdown",
		"owner": "admin",
		"updated": "2024-03-12T20:39:12.149Z",
		"resourceMissing": false,
		"choiceSettings": {
		  "choiceSet": "dbConnections",
		  "family": "PostgreSQL"
		}
	  }`
	responseSingleWeb := `{
		"type": "dropdown",
		"name": "testWeb",
		"value": "slackConnection",
		"type": "dropdown",
		"owner": "admin",
		"updated": "2024-03-15T17:42:56.752Z",
		"resourceMissing": true,
		"choiceSettings": {
		  "choiceSet": "webConnections",
		  "services": [
			"Slack"
		  ]
		}
	  }`

	customHttpServerHandlerText := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseSingleText))
			require.NoError(t, err)
		}

		if r.Method == "PUT" {
			// check we are updating the right value
			require.Contains(t, r.URL.Path, "testText")
			// check the json body is correct
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result UpdateDeploymentParameter
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, "text", result.Type)
			require.Equal(t, "myValue", result.Value)
			w.WriteHeader(http.StatusNoContent)
			_, err = w.Write([]byte(""))
			require.NoError(t, err)
		}

	}
	customHttpServerHandlerDb := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseSingleDb))
			require.NoError(t, err)
		}

		if r.Method == "PUT" {
			// check we are updating the right value
			require.Contains(t, r.URL.Path, "testDb")
			// check the json body is correct
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result UpdateDeploymentParameter
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, "dropdown", result.Type)
			require.Equal(t, "myValue", result.Value)
			require.Equal(t, "dbConnections", result.ChoiceSettings.ChoiceSet)
			w.WriteHeader(http.StatusNoContent)
			_, err = w.Write([]byte(""))
			require.NoError(t, err)
		}

	}
	customHttpServerHandlerWeb := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseSingleWeb))
			require.NoError(t, err)
		}

		if r.Method == "PUT" {
			// check we are updating the right value
			require.Contains(t, r.URL.Path, "testWeb")
			// check the json body is correct
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result UpdateDeploymentParameter
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, "dropdown", result.Type)
			require.Equal(t, "myValue", result.Value)
			require.Equal(t, "webConnections", result.ChoiceSettings.ChoiceSet)
			w.WriteHeader(http.StatusNoContent)
			_, err = w.Write([]byte(""))
			require.NoError(t, err)
		}

	}

	customHttpServerHandlerWebExcludeInclude := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseSingleWeb))
			require.NoError(t, err)
		}

		if r.Method == "PUT" {
			// check we are updating the right value
			require.Contains(t, r.URL.Path, "testWeb")
			// check the json body is correct
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var result UpdateDeploymentParameter
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)
			require.Equal(t, "dropdown", result.Type)
			require.Equal(t, "myValue", result.Value)
			require.Equal(t, "webConnections", result.ChoiceSettings.ChoiceSet)
			require.Equal(t, []string{"Slack"}, result.ChoiceSettings.Services)
			require.Equal(t, []string{"Teams"}, result.ChoiceSettings.ExcludedServices)
			w.WriteHeader(http.StatusNoContent)
			_, err = w.Write([]byte(""))
			require.NoError(t, err)
		}

	}

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"deploymentparameters", "update", "--name", "myDep", "--value", "myValue", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"deploymentparameters", "update", "--name", "myDep", "--value", "myValue"},
		},
		{
			name:        "missing name flag",
			wantErrText: "required flag(s) \"name\", \"value\" not set",
			args:        []string{"deploymentparameters", "update"},
		},
		{
			name:        "missing value flag",
			wantErrText: "required flag(s) \"value\" not set",
			args:        []string{"deploymentparameters", "update", "--name", "myDep"},
		},
		{
			name:            "update parameter text",
			args:            []string{"deploymentparameters", "update", "--name", "testText", "--value", "myValue"},
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandlerText)),
			wantOutputRegex: "^Deployment Parameter successfully updated.[\\s]*$",
		},
		{
			name:            "update parameter db",
			args:            []string{"deploymentparameters", "update", "--name", "testDb", "--value", "myValue"},
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandlerDb)),
			wantOutputRegex: "^Deployment Parameter successfully updated.[\\s]*$",
		},
		{
			name:            "update parameter web",
			args:            []string{"deploymentparameters", "update", "--name", "testWeb", "--value", "myValue"},
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandlerWeb)),
			wantOutputRegex: "^Deployment Parameter successfully updated.[\\s]*$",
		},
		{
			name:            "update parameter web include exclude",
			args:            []string{"deploymentparameters", "update", "--name", "testWeb", "--value", "myValue", "--type", "web", "--included-service", "Slack", "--excluded-service", "Teams"},
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandlerWebExcludeInclude)),
			wantOutputRegex: "^Deployment Parameter successfully updated.[\\s]*$",
		},
		{
			name:        "parameter does not exist",
			statusCode:  http.StatusConflict,
			body:        parameterDoesNotExistBody,
			args:        []string{"deploymentparameters", "update", "--name", "myDep", "--value", "myValue"},
			wantErrText: "Unauthorized request by user admin due to lack of proper permissions or the object does not exist.",
		},
	}

	runTests(cases, t)

}
