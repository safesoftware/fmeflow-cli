package cmd

import (
	"net/http"
	"testing"
)

func TestDeploymentParametersCreate(t *testing.T) {
	parameterExistsBody := `{
		"message": "A deployment parameter with name \"myDep\" already exists."
	  }`
	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue"},
		},
		{
			name:        "missing name flag",
			wantErrText: "required flag(s) \"name\" not set",
			args:        []string{"deploymentparameters", "create"},
		},
		{
			name:            "create parameter",
			statusCode:      http.StatusCreated,
			args:            []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue"},
			wantOutputRegex: "^Deployment Parameter successfully created.[\\s]*$",
		},
		{
			name:        "parameter already exists",
			statusCode:  http.StatusConflict,
			body:        parameterExistsBody,
			args:        []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue"},
			wantErrText: "A deployment parameter with name \"myDep\" already exists.",
		},
		{
			name:               "create parameter with  invalid type",
			statusCode:         http.StatusCreated,
			args:               []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue", "--type", "invalid"},
			wantErrOutputRegex: "invalid argument \"invalid\" for \"--type\" flag: must be one of \"text\" or \"database\" or \"web\"",
		},
		{
			name:            "create parameter with text type",
			statusCode:      http.StatusCreated,
			args:            []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue", "--type", "text"},
			wantOutputRegex: "^Deployment Parameter successfully created.[\\s]*$",
		},
		{
			name:            "create parameter with web type",
			statusCode:      http.StatusCreated,
			args:            []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue", "--type", "web"},
			wantBodyRegEx:   `{"name":"myDep","type":"dropdown","value":"myValue","choiceSettings":{"choiceSet":"webConnections"}}`,
			wantOutputRegex: "^Deployment Parameter successfully created.[\\s]*$",
		},
		{
			name:            "create parameter with web type and included services",
			statusCode:      http.StatusCreated,
			args:            []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue", "--type", "web", "--included-service", "service1", "--included-service", "service2"},
			wantBodyRegEx:   `{"name":"myDep","type":"dropdown","value":"myValue","choiceSettings":{"choiceSet":"webConnections","services":\["service1","service2"\]}}`,
			wantOutputRegex: "^Deployment Parameter successfully created.[\\s]*$",
		},
		{
			name:            "create parameter with web type and excluded services",
			statusCode:      http.StatusCreated,
			args:            []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue", "--type", "web", "--excluded-service", "service1", "--excluded-service", "service2"},
			wantBodyRegEx:   `{"name":"myDep","type":"dropdown","value":"myValue","choiceSettings":{"choiceSet":"webConnections","excludedServices":\["service1","service2"\]}}`,
			wantOutputRegex: "^Deployment Parameter successfully created.[\\s]*$",
		},
		{
			name:            "create parameter with db type",
			statusCode:      http.StatusCreated,
			args:            []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue", "--type", "database"},
			wantBodyRegEx:   `{"name":"myDep","type":"dropdown","value":"myValue","choiceSettings":{"choiceSet":"dbConnections"}}`,
			wantOutputRegex: "^Deployment Parameter successfully created.[\\s]*$",
		},
		{
			name:            "create parameter with db type and dbType flag",
			statusCode:      http.StatusCreated,
			args:            []string{"deploymentparameters", "create", "--name", "myDep", "--value", "myValue", "--type", "database", "--database-type", "dbType"},
			wantBodyRegEx:   `{"name":"myDep","type":"dropdown","value":"myValue","choiceSettings":{"choiceSet":"dbConnections","family":"dbType"}}`,
			wantOutputRegex: "^Deployment Parameter successfully created.[\\s]*$",
		},
	}

	runTests(cases, t)

}
