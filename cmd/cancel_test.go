package cmd

import (
	"net/http"
	"testing"
)

func TestCancel(t *testing.T) {
	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"cancel", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"cancel", "--id", "1"},
		},
		{
			name:        "invalid job id v3",
			statusCode:  http.StatusNotFound,
			wantErrText: "the specified job ID was not found",
			args:        []string{"cancel", "--id", "1", "--api-version", "v3"},
		},
		{
			name:            "cancel valid job v3",
			statusCode:      http.StatusNoContent,
			args:            []string{"cancel", "--id", "1234", "--api-version", "v3"},
			wantOutputRegex: "Success. The job with id 1234 was cancelled.",
		},
		{
			name:            "cancel valid job json v3",
			statusCode:      http.StatusNoContent,
			args:            []string{"cancel", "--id", "1234", "--json", "--api-version", "v3"},
			wantOutputRegex: "{}",
		},
		{
			name:        "job already complete",
			statusCode:  http.StatusUnprocessableEntity,
			body:        `{"message": "Job \"1234\" is already complete and cannot be cancelled."}`,
			args:        []string{"cancel", "--id", "1234", "--api-version", "v4"},
			wantErrText: "Job \"1234\" is already complete and cannot be cancelled.",
		},
		{
			name:        "job id does not exist",
			statusCode:  http.StatusUnprocessableEntity,
			body:        `{"message": "The job for ID \"55\" does not exist."}`,
			args:        []string{"cancel", "--id", "1234", "--api-version", "v4"},
			wantErrText: "The job for ID \"55\" does not exist.",
		},
		{
			name:           "job already complete json",
			statusCode:     http.StatusUnprocessableEntity,
			body:           `{"message": "Job \"1234\" is already complete and cannot be cancelled."}`,
			args:           []string{"cancel", "--id", "1234", "--json", "--api-version", "v4"},
			wantErrText:    "Job \"1234\" is already complete and cannot be cancelled.",
			wantOutputJson: `{"message": "Job \"1234\" is already complete and cannot be cancelled."}`,
		},
		{
			name:           "job id does not exist json",
			statusCode:     http.StatusUnprocessableEntity,
			body:           `{"message": "The job for ID \"55\" does not exist."}`,
			args:           []string{"cancel", "--id", "1234", "--json", "--api-version", "v4"},
			wantErrText:    "The job for ID \"55\" does not exist.",
			wantOutputJson: `{"message": "The job for ID \"55\" does not exist."}`,
		},
		{
			name:            "cancel valid job",
			statusCode:      http.StatusNoContent,
			args:            []string{"cancel", "--id", "1234", "--api-version", "v4"},
			wantOutputRegex: "Success. The job with id 1234 was cancelled.",
		},
		{
			name:            "cancel valid job json",
			statusCode:      http.StatusNoContent,
			args:            []string{"cancel", "--id", "1234", "--json", "--api-version", "v4"},
			wantOutputRegex: "{}",
		},
	}

	runTests(cases, t)

}
