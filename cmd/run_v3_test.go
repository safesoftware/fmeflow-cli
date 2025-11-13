package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	responseV3ASync := `{
		"id": 1
	  }`

	responseV3Sync := `{
		"timeRequested": "2023-02-04T00:16:28Z",
		"requesterResultPort": 37805,
		"numFeaturesOutput": 1539,
		"requesterHost": "10.1.113.39",
		"timeStarted": "2023-02-04T00:16:28Z",
		"id": 1,
		"timeFinished": "2023-02-04T00:16:30Z",
		"priority": -1,
		"statusMessage": "Translation Successful",
		"status": "SUCCESS"
	  }`

	dataFileContents := "Pretend backup file"

	// generate random file to restore from
	f, err := os.CreateTemp("", "datafile")
	require.NoError(t, err)
	defer os.Remove(f.Name()) // clean up
	err = os.WriteFile(f.Name(), []byte(dataFileContents), 0644)
	require.NoError(t, err)

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw"},
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
			body:            responseV3Sync,
			wantOutputRegex: "^[\\s]*ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539[\\s]*$",
		},
		{
			name:            "run async job regular output",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
		},
		{
			name:           "run async job json",
			statusCode:     http.StatusOK,
			body:           responseV3ASync,
			args:           []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--json"},
			wantOutputJson: responseV3ASync,
		},
		{
			name:           "run sync job json output",
			statusCode:     http.StatusOK,
			body:           responseV3Sync,
			args:           []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--json", "--wait"},
			wantOutputJson: responseV3Sync,
		},
		{
			name:            "description flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--description", "My Description"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"TMDirectives\".*:[\\s]*{.*\"description\":\"My Description\".*",
		},
		{
			name:            "failure topic flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--failure-topic", "FAILURE_TOPIC"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"NMDirectives\".*:[\\s]*{.*\"failureTopics\":\\[\"FAILURE_TOPIC\"\\].*",
		},
		{
			name:            "success topic flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--success-topic", "SUCCESS_TOPIC"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"NMDirectives\".*:[\\s]*{.*\"successTopics\":\\[\"SUCCESS_TOPIC\"\\].*",
		},
		{
			name:            "node manager directive flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--node-manager-directive", "directive1=value1", "--node-manager-directive", "directive2=value2"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"NMDirectives\".*:[\\s]*{.*\"directives\":\\[{\"name\":\"directive1\",\"value\":\"value1\"},{\"name\":\"directive2\",\"value\":\"value2\".*",
		},
		{
			name:            "run until canceled flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--run-until-canceled"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"TMDirectives\":{\"rtc\":true}.*",
		},
		{
			name:            "tag flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--tag", "myqueue"},
			wantOutputRegex: "Job submitted with id: 1",
			wantBodyRegEx:   ".*\"TMDirectives\":{.*\"tag\":\"myqueue\".*}.*",
		},
		{
			name:            "queue flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--queue", "myqueue"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"TMDirectives\":{.*\"tag\":\"myqueue\".*}.*",
		},

		{
			name:            "time to live flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--time-to-live", "60"},
			wantOutputRegex: "Job submitted with id: 1",
			wantBodyRegEx:   ".*\"TMDirectives\":{.*\"ttl\":60.*}.*",
		},
		{
			name:            "max time in queue flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--max-time-in-queue", "60"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"TMDirectives\":{.*\"ttl\":60.*}.*",
		},

		{
			name:            "timeuntil canceled flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--time-until-canceled", "60"},
			wantOutputRegex: "Job submitted with id: 1",
			wantBodyRegEx:   ".*\"TMDirectives\":{.*\"ttc\":60.*}.*",
		},
		{
			name:            "max job runtime flag async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--max-job-runtime", "60"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"TMDirectives\":{.*\"ttc\":60.*}.*",
		},

		{
			name:            "published parameter async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--published-parameter", "COORDSYS=TX83-CF"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"publishedParameters\":\\[{\"value\":\"TX83-CF\",\"name\":\"COORDSYS\".*}.*",
		},
		{
			name:            "published parameter list async",
			statusCode:      http.StatusOK,
			body:            responseV3ASync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--published-parameter-list", "THEMES=railroad,airports"},
			wantOutputRegex: "^[\\s]*Job submitted with id: 1[\\s]*$",
			wantBodyRegEx:   ".*\"publishedParameters\":\\[{\"value\":\\[\"railroad\",\"airports\"],\"name\":\"THEMES\".*}.*",
		},

		{
			name:            "description flag transact data",
			statusCode:      http.StatusOK,
			body:            responseV3Sync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--description", "My Description", "--file", f.Name()},
			wantOutputRegex: "^[\\s]*ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539[\\s]*$",
			wantFormParams:  map[string]string{"opt_description": "My Description"},
		},
		{
			name:            "failure topic flag transact data",
			statusCode:      http.StatusOK,
			body:            responseV3Sync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--failure-topic", "FAILURE_TOPIC", "--file", f.Name()},
			wantOutputRegex: "^[\\s]*ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539[\\s]*$",
			wantFormParams:  map[string]string{"opt_failuretopics": "FAILURE_TOPIC"},
		},
		{
			name:            "success topic flag transact data",
			statusCode:      http.StatusOK,
			body:            responseV3Sync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--success-topic", "SUCCESS_TOPIC", "--file", f.Name()},
			wantOutputRegex: "^[\\s]*ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539[\\s]*$",
			wantFormParams:  map[string]string{"opt_successtopics": "SUCCESS_TOPIC"},
		},
		{
			name:            "tag flag transact data",
			statusCode:      http.StatusOK,
			body:            responseV3Sync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--tag", "myqueue", "--file", f.Name()},
			wantOutputRegex: "ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539",
			wantFormParams:  map[string]string{"opt_tag": "myqueue"},
		},
		{
			name:            "time to live flag transact data",
			statusCode:      http.StatusOK,
			body:            responseV3Sync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--time-to-live", "60", "--file", f.Name()},
			wantOutputRegex: "ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539",
			wantFormParams:  map[string]string{"opt_ttl": "60"},
		},
		{
			name:            "timeuntil canceled flag transact data",
			statusCode:      http.StatusOK,
			body:            responseV3Sync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--time-until-canceled", "60", "--file", f.Name()},
			wantOutputRegex: "ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539",
			wantFormParams:  map[string]string{"opt_ttc": "60"},
		},
		{
			name:            "published parameter transact data",
			statusCode:      http.StatusOK,
			body:            responseV3Sync,
			args:            []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--published-parameter", "COORDSYS=TX83-CF", "--file", f.Name()},
			wantOutputRegex: "^[\\s]*ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539[\\s]*$",
			wantFormParams:  map[string]string{"COORDSYS": "TX83-CF"},
		},
		{
			name:               "published parameter list transact data",
			statusCode:         http.StatusOK,
			body:               responseV3Sync,
			args:               []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--published-parameter-list", "THEMES=railroad,airports", "--file", f.Name()},
			wantOutputRegex:    "^[\\s]*ID[\\s]*STATUS[\\s]*STATUS MESSAGE[\\s]*FEATURES OUTPUT[\\s]*1[\\s]*SUCCESS[\\s]*Translation Successful[\\s]*1539[\\s]*$",
			wantFormParamsList: map[string][]string{"THEMES": {"railroad", "airports"}},
		},
		{
			name:        "transact data node manager mutually exclusive",
			args:        []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--node-manager-directive", "directive1=value1", "--file", f.Name()},
			wantErrText: "if any flags in the group [file node-manager-directive] are set none of the others can be; [file node-manager-directive] were all set",
		},
		{
			name:        "transact data run until canceled mutually exclusive",
			args:        []string{"run", "--repository", "Samples", "--workspace", "austinApartments.fmw", "--run-until-canceled", "--file", f.Name()},
			wantErrText: "if any flags in the group [file run-until-canceled] are set none of the others can be; [file run-until-canceled] were all set",
		},
	}
	runTests(cases, t)

}
