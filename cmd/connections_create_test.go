package cmd

import (
	"net/http"
	"testing"
)

func TestConnectionsCreate(t *testing.T) {
	responseOracleMissing := `{
		"message": "The connection parameter(s) are required but are missing or empty: CONNECTION_MODE, DATASET, WALLET_PATH, SELECTED_SERVICE."
	  }`

	responseDatabaseTypeMissing := `{
		"message": "Specified database type does not exist."
	  }`

	responseParameterValidationFailed := `{
		"message": "Parameter Validation Failed",
		"details": {
		  "authenticationMethod": "Connection authentication method must be supplied",
		  "type": "must not be blank"
		}
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"connections", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"connections"},
		},
		{
			name:            "create connection",
			statusCode:      http.StatusCreated,
			args:            []string{"connections", "create", "--name", "test123aa", "--category", "database", "--type", "PostgreSQL", "--username", "test", "--password", "test", "--parameter", "HOST=a", "--parameter", "PORT=5432", "--parameter", "DATASET=dbname", "--parameter", "USER_NAME=a", "--parameter", "SSL_OPTIONS=a", "--parameter", "SSLMODE=prefer"},
			wantOutputRegex: "^[\\s]*Connection successfully created.[\\s]*$",
		},
		{
			name:           "create connection json output",
			statusCode:     http.StatusCreated,
			args:           []string{"connections", "create", "--name", "test123aa", "--category", "database", "--type", "PostgreSQL", "--username", "test", "--password", "test", "--parameter", "HOST=a", "--parameter", "PORT=5432", "--parameter", "DATASET=dbname", "--parameter", "USER_NAME=a", "--parameter", "SSL_OPTIONS=a", "--parameter", "SSLMODE=prefer", "--json"},
			wantOutputJson: "{}",
		},
		{
			name:        "create connection missing oracle parameters",
			statusCode:  http.StatusBadRequest,
			args:        []string{"connections", "create", "--name", "test123aa", "--category", "database", "--type", "Oracle"},
			body:        responseOracleMissing,
			wantErrText: "The connection parameter(s) are required but are missing or empty: CONNECTION_MODE, DATASET, WALLET_PATH, SELECTED_SERVICE.",
		},
		{
			name:        "create connection missing database type",
			statusCode:  http.StatusBadRequest,
			args:        []string{"connections", "create", "--name", "test123aa", "--category", "database", "--type", "test"},
			body:        responseDatabaseTypeMissing,
			wantErrText: "Specified database type does not exist.",
		},
		{
			name:        "create connection missing parameters",
			statusCode:  http.StatusBadRequest,
			args:        []string{"connections", "create", "--name", "test123aa", "--category", "basic"},
			body:        responseParameterValidationFailed,
			wantErrText: "Parameter Validation Failed\nauthenticationMethod: Connection authentication method must be supplied\ntype: must not be blank",
		},
	}

	runTests(cases, t)

}
