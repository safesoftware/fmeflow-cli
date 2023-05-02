package cmd

import (
	"net/http"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	// standard responses for v3 and v4
	okResponseV3 := `{
		"status": "ok"
		}`
	okResponseV4 := `{
		"status": "ok",
		"message": "FME Server is healthy."
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"healthcheck", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"healthcheck"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"healthcheck"},
		},
		{
			name:            "v3 health check ok",
			statusCode:      http.StatusOK,
			args:            []string{"healthcheck", "--api-version", "v3"},
			body:            okResponseV3,
			wantOutputRegex: "ok",
		},
		{
			name:            "v3 health check ready ok",
			statusCode:      http.StatusOK,
			args:            []string{"healthcheck", "--ready", "--api-version", "v3"},
			body:            okResponseV3,
			wantOutputRegex: "ok",
		},
		{
			name:            "v4 health check ok",
			statusCode:      http.StatusOK,
			body:            okResponseV4,
			wantOutputRegex: "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
			args:            []string{"healthcheck"},
		},
		{
			name:            "v4 health check ready ok",
			statusCode:      http.StatusOK,
			body:            okResponseV4,
			args:            []string{"healthcheck", "--ready"},
			wantOutputRegex: "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
		},
		{
			name:            "v3 health check used for 2022.2 build",
			statusCode:      http.StatusOK,
			body:            okResponseV3,
			wantOutputRegex: "^ok\n$",
			fmeflowBuild:    22765,
			args:            []string{"healthcheck"},
		},
		{
			name:            "v4 health check used for 2023.0 build",
			statusCode:      http.StatusOK,
			body:            okResponseV4,
			wantOutputRegex: "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
			fmeflowBuild:    23200,
			args:            []string{"healthcheck"},
		},
		{
			name:       "extra json fields",
			statusCode: http.StatusOK,
			body: `{
				"status": "ok",
				"message": "FME Server is healthy.",
				"extra": "Extra field"
				}`,
			wantOutputRegex: "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
			fmeflowBuild:    23200,
			args:            []string{"healthcheck"},
		},
		{
			name:           "json output v4",
			statusCode:     http.StatusOK,
			body:           okResponseV4,
			wantOutputJson: okResponseV4,
			args:           []string{"healthcheck", "--json"},
		},
		{
			name:           "json output v3",
			statusCode:     http.StatusOK,
			body:           okResponseV3,
			wantOutputJson: okResponseV3,
			args:           []string{"healthcheck", "--json", "--api-version", "v3"},
		},
		{
			name:            "v4 health check with url flag",
			statusCode:      http.StatusOK,
			body:            okResponseV4,
			wantOutputRegex: "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
			args:            []string{"healthcheck", "--url", urlPlaceholder},
			omitConfig:      true,
		},
		{
			name:            "v4 health check with no token in config file",
			statusCode:      http.StatusOK,
			body:            okResponseV4,
			wantOutputRegex: "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
			args:            []string{"healthcheck"},
			omitConfigToken: true,
		},
	}
	runTests(cases, t)
}
