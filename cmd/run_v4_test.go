package cmd

import (
	"net/http"
	"testing"
)

func TestRunV4(t *testing.T) {
	responseV4ASync := `{
		"id": 1
	}`

	responseV4Sync := `{
		"id": 1,
		"featureOutputCount": 1539,
		"requesterHost": "10.1.113.39",
		"requesterResultPort": 37805,
		"status": "SUCCESS",
		"statusMessage": "Translation Successful",
		"timeFinished": "2023-02-04T00:16:30Z",
		"timeQueued": "2023-02-04T00:16:28Z",
		"timeStarted": "2023-02-04T00:16:28Z"
	}`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:         "500 bad status code",
			statusCode:   http.StatusInternalServerError,
			wantErrText:  "500 Internal Server Error",
			args:         []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw"},
			fmeflowBuild: 26018,
		},
		{
			name:        "repository flag required",
			wantErrText: "required flag(s) \"repository\" not set",
			args:        []string{"run", "--workspace", "austinApartments.fmw"},
		},
		{
			name:        "workspace flag required",
			wantErrText: "required flag(s) \"workspace\" not set",
			args:        []string{"run", "--repository", "Samples"},
		},
		{
			name:            "run sync job table output",
			statusCode:      http.StatusOK,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--wait"},
			body:            responseV4Sync,
			wantOutputRegex: "^[\\s]*ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539[\\s]*$",
			fmeflowBuild:    26018,
		},
		{
			name:            "run async job regular output",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			fmeflowBuild:    26018,
		},
		{
			name:           "run async job json",
			statusCode:     http.StatusOK,
			body:           responseV4ASync,
			args:           []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--json"},
			wantOutputJson: responseV4ASync,
			fmeflowBuild:   26018,
		},
		{
			name:           "run sync job json output",
			statusCode:     http.StatusOK,
			body:           responseV4Sync,
			args:           []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--json", "--wait"},
			wantOutputJson: responseV4Sync,
			fmeflowBuild:   26018,
		},
		{
			name:            "failure topic flag async",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--failure-topic", "FAILURE_TOPIC"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"failureTopics\":\\[\"FAILURE_TOPIC\"\\].*",
			fmeflowBuild:    26018,
		},
		{
			name:            "success topic flag async",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--success-topic", "SUCCESS_TOPIC"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"successTopics\":\\[\"SUCCESS_TOPIC\"\\].*",
			fmeflowBuild:    26018,
		},
		{
			name:            "directive flag async",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--directive", "directive1=value1", "--directive", "directive2=value2"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"directives\":{.*\"directive1\":\"value1\".*\"directive2\":\"value2\".*}.*",
			fmeflowBuild:    26018,
		},
		{
			name:            "published parameter async",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--published-parameter", "COORDSYS=TX83-CF"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"publishedParameters\":{.*\"COORDSYS\":\"TX83-CF\".*}.*",
			fmeflowBuild:    26018,
		},
		{
			name:            "published parameter list async",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--published-parameter-list", "THEMES=railroad,airports"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"publishedParameters\":{.*\"THEMES\":\\[\"railroad\",\"airports\"\\].*}.*",
			fmeflowBuild:    26018,
		},
		{
			name:            "max job runtime flag async",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--max-job-runtime", "10"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"maxJobRuntime\":10.*",
			fmeflowBuild:    26018,
		},
		{
			name:            "time until canceled flag async deprecated",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--time-until-canceled", "10"},
			wantOutputRegex: "Flag --time-until-canceled has been deprecated, please use --max-job-runtime instead[\\s\\S]*Job submitted with id: 1",
			wantBodyRegEx:   ".*\"maxJobRuntime\":10.*",
			fmeflowBuild:    26018,
		},
		{
			name:            "max job runtime invalid value ignored",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--max-job-runtime", "-5"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			fmeflowBuild:    26018,
		},
		{
			name:            "max time in queue flag async",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--max-time-in-queue", "60"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"maxTimeInQueue\":60.*",
			fmeflowBuild:    26018,
		},
		{
			name:            "time to live flag async deprecated",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--time-to-live", "60"},
			wantOutputRegex: "Flag --time-to-live has been deprecated, please use --max-time-in-queue instead[\\s\\S]*Job submitted with id: 1",
			wantBodyRegEx:   ".*\"maxTimeInQueue\":60.*",
			fmeflowBuild:    26018,
		},
		{
			name:            "max time in queue invalid value ignored",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--max-time-in-queue", "-1"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			fmeflowBuild:    26018,
		},
		{
			name:            "max total life time flag sync",
			statusCode:      http.StatusOK,
			body:            responseV4Sync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--wait", "--max-total-life-time", "300"},
			wantOutputRegex: "^[\\s]*ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539[\\s]*$",
			wantBodyRegEx:   ".*\"maxTotalLifeTime\":300.*",
			fmeflowBuild:    26018,
		},
		{
			name:            "max total life time invalid value ignored",
			statusCode:      http.StatusOK,
			body:            responseV4Sync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--wait", "--max-total-life-time", "100000"},
			wantOutputRegex: "^[\\s]*ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539[\\s]*$",
			fmeflowBuild:    26018,
		},
		{
			name:            "queue flag async",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--queue", "MyQueue"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"queue\":\"MyQueue\".*",
			fmeflowBuild:    26018,
		},
		{
			name:            "tag flag async deprecated",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--tag", "MyQueue"},
			wantOutputRegex: "Flag --tag has been deprecated, please use --queue instead[\\s\\S]*Job submitted with id: 1",
			wantBodyRegEx:   ".*\"queue\":\"MyQueue\".*",
			fmeflowBuild:    26018,
		},
		{
			name:            "published parameter and list combined",
			statusCode:      http.StatusOK,
			body:            responseV4ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--published-parameter", "COORDSYS=TX83-CF", "--published-parameter-list", "THEMES=railroad,airports"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"publishedParameters\":{.*\"COORDSYS\":\"TX83-CF\".*\"THEMES\":\\[\"railroad\",\"airports\"\\].*}.*",
			fmeflowBuild:    26018,
		},
	}
	runTests(cases, t)
}
