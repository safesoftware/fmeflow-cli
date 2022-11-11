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
			name:        "unknown flag",
			statusCode:  http.StatusOK,
			args:        []string{"--badflag"},
			wantErrText: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
		},
		{
			name:       "v3 health check ok",
			statusCode: http.StatusOK,
			args:       []string{"--api-version", "v3"},
			body:       okResponseV3,
			wantOutput: "ok",
		},
		{
			name:       "v3 health check ready ok",
			statusCode: http.StatusOK,
			args:       []string{"--ready", "--api-version", "v3"},
			body:       okResponseV3,
			wantOutput: "ok",
		},
		{
			name:       "v4 health check ok",
			statusCode: http.StatusOK,
			body:       okResponseV4,
			wantOutput: "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
		},
		{
			name:       "v4 health check ready ok",
			statusCode: http.StatusOK,
			body:       okResponseV4,
			args:       []string{"--ready"},
			wantOutput: "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
		},
		{
			name:           "v3 health check used for 2022.2 build",
			statusCode:     http.StatusOK,
			body:           okResponseV3,
			wantOutput:     "^ok\n$",
			fmeserverBuild: 22765,
		},
		{
			name:           "v4 health check used for 2023.0 build",
			statusCode:     http.StatusOK,
			body:           okResponseV4,
			wantOutput:     "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
			fmeserverBuild: 23200,
		},
		{
			name:       "extra json fields",
			statusCode: http.StatusOK,
			body: `{
				"status": "ok",
				"message": "FME Server is healthy.",
				"extra": "Extra field"
				}`,
			wantOutput:     "STATUS[\\s]*MESSAGE[\\s]*[\\s]*ok[\\s]*FME Server is healthy",
			fmeserverBuild: 23200,
		},
	}
	runTests(cases, newHealthcheckCmd, t)
}
