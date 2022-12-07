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
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "the specified job ID was not found",
			args:        []string{"cancel", "--id", "1"},
		},
		{
			name:            "cancel valid job",
			statusCode:      http.StatusNoContent,
			args:            []string{"cancel", "--id", "1234"},
			wantOutputRegex: "",
		},
		{
			name:            "cancel valid job json",
			statusCode:      http.StatusNoContent,
			args:            []string{"cancel", "--id", "1234", "--json"},
			wantOutputRegex: "{}",
		},
	}

	runTests(cases, t)

}
