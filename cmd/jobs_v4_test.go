package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJobsV4(t *testing.T) {
	responseV4Completed := `{
		"offset": 0,
		"limit": 0,
		"totalCount": 3,
		"items": [
			{
				"id": 1,
				"description": "",
				"engineHost": "387f74cd4e1f",
				"engineName": "387f74cd4e1f",
				"repository": "Samples",
				"queue": "Default",
				"queueType": "Standard",
				"resultDatasetDownloadUrl": "",
				"status": "success",
				"timeFinished": "2023-11-08T21:52:50Z",
				"timeQueued": "2023-11-08T21:52:04Z",
				"timeStarted": "2023-11-08T21:52:49Z",
				"runtimeUsername": "admin",
				"runtimeUserID": "5d288fb0-273a-4fba-a19c-a3259ac1487a",
				"workspace": "austinApartments.fmw",
				"elapsedTime": 1118,
				"cpuTime": 994,
				"cpuPercent": 88.9,
				"peakMemoryUsage": 10988808,
				"lineCount": 117,
				"warningCount": 0,
				"errorCount": 0
			},
			{
				"id": 2,
				"description": "",
				"engineHost": "10f259e906e5",
				"engineName": "10f259e906e5",
				"repository": "Test",
				"queue": "Default",
				"queueType": "Standard",
				"resultDatasetDownloadUrl": "",
				"status": "failure",
				"timeFinished": "2023-11-15T00:42:31Z",
				"timeQueued": "2023-11-15T00:42:30Z",
				"timeStarted": "2023-11-15T00:42:30Z",
				"runtimeUsername": "admin",
				"runtimeUserID": "5d288fb0-273a-4fba-a19c-a3259ac1487a",
				"workspace": "none2none.fmw",
				"elapsedTime": 850,
				"cpuTime": 697,
				"cpuPercent": 82.0,
				"peakMemoryUsage": 5190544,
				"lineCount": 70,
				"warningCount": 0,
				"errorCount": 2
			},
			{
				"id": 3,
				"description": "",
				"engineHost": "145929514b24",
				"engineName": "145929514b24",
				"repository": "Test",
				"queue": "Priority",
				"queueType": "Priority",
				"resultDatasetDownloadUrl": "",
				"status": "cancelled",
				"timeFinished": "2023-12-07T21:23:12Z",
				"timeQueued": "2023-12-07T21:22:48Z",
				"timeStarted": "2023-12-07T21:22:48Z",
				"runtimeUsername": "admin",
				"runtimeUserID": "5d288fb0-273a-4fba-a19c-a3259ac1487a",
				"workspace": "cancelled.fmw",
				"elapsedTime": 0,
				"cpuTime": 0,
				"cpuPercent": 0,
				"peakMemoryUsage": 0,
				"lineCount": 0,
				"warningCount": 0,
				"errorCount": 0
			}
		]
	}`

	responseV4Active := `{
		"offset": 0,
		"limit": 100,
		"totalCount": 2,
		"items": [
			{
				"id": 4,
				"description": "",
				"engineHost": "10f259e906e5",
				"engineName": "10f259e906e5",
				"repository": "Test",
				"queue": "Default",
				"queueType": "Standard",
				"resultDatasetDownloadUrl": "",
				"status": "running",
				"timeFinished": "0001-01-01T00:00:00Z",
				"timeQueued": "2023-11-15T01:22:14Z",
				"timeStarted": "2023-11-15T01:22:14Z",
				"runtimeUsername": "admin",
				"runtimeUserID": "5d288fb0-273a-4fba-a19c-a3259ac1487a",
				"workspace": "running.fmw",
				"elapsedTime": 0,
				"cpuTime": 0,
				"cpuPercent": 0,
				"peakMemoryUsage": 0,
				"lineCount": 0,
				"warningCount": 0,
				"errorCount": 0
			},
			{
				"id": 5,
				"description": "",
				"engineHost": "",
				"engineName": "",
				"repository": "Samples",
				"queue": "Default",
				"queueType": "Standard",
				"resultDatasetDownloadUrl": "",
				"status": "queued",
				"timeFinished": "0001-01-01T00:00:00Z",
				"timeQueued": "2023-11-08T21:52:04Z",
				"timeStarted": "0001-01-01T00:00:00Z",
				"runtimeUsername": "admin",
				"runtimeUserID": "5d288fb0-273a-4fba-a19c-a3259ac1487a",
				"workspace": "austinApartments.fmw",
				"elapsedTime": 0,
				"cpuTime": 0,
				"cpuPercent": 0,
				"peakMemoryUsage": 0,
				"lineCount": 0,
				"warningCount": 0,
				"errorCount": 0
			}
		]
	}`

	responseV4Queued := `{
		"offset": 0,
		"limit": 0,
		"totalCount": 1,
		"items": [
			{
				"id": 5,
				"description": "",
				"engineHost": "",
				"engineName": "",
				"repository": "Samples",
				"queue": "Default",
				"queueType": "Standard",
				"resultDatasetDownloadUrl": "",
				"status": "queued",
				"timeFinished": "0001-01-01T00:00:00Z",
				"timeQueued": "2023-11-08T21:52:04Z",
				"timeStarted": "0001-01-01T00:00:00Z",
				"runtimeUsername": "admin",
				"runtimeUserID": "5d288fb0-273a-4fba-a19c-a3259ac1487a",
				"workspace": "austinApartments.fmw",
				"elapsedTime": 0,
				"cpuTime": 0,
				"cpuPercent": 0,
				"peakMemoryUsage": 0,
				"lineCount": 0,
				"warningCount": 0,
				"errorCount": 0
			}
		]
	}`

	responseV4SuccessAndFailure := `{
		"offset": 0,
		"limit": 0,
		"totalCount": 2,
		"items": [
			{
				"id": 1,
				"description": "",
				"engineHost": "387f74cd4e1f",
				"engineName": "387f74cd4e1f",
				"repository": "Samples",
				"queue": "Default",
				"queueType": "Standard",
				"resultDatasetDownloadUrl": "",
				"status": "success",
				"timeFinished": "2023-11-08T21:52:50Z",
				"timeQueued": "2023-11-08T21:52:04Z",
				"timeStarted": "2023-11-08T21:52:49Z",
				"runtimeUsername": "admin",
				"runtimeUserID": "5d288fb0-273a-4fba-a19c-a3259ac1487a",
				"workspace": "austinApartments.fmw",
				"elapsedTime": 1118,
				"cpuTime": 994,
				"cpuPercent": 88.9,
				"peakMemoryUsage": 10988808,
				"lineCount": 117,
				"warningCount": 0,
				"errorCount": 0
			},
			{
				"id": 2,
				"description": "",
				"engineHost": "10f259e906e5",
				"engineName": "10f259e906e5",
				"repository": "Test",
				"queue": "Default",
				"queueType": "Standard",
				"resultDatasetDownloadUrl": "",
				"status": "failure",
				"timeFinished": "2023-11-15T00:42:31Z",
				"timeQueued": "2023-11-15T00:42:30Z",
				"timeStarted": "2023-11-15T00:42:30Z",
				"runtimeUsername": "admin",
				"runtimeUserID": "5d288fb0-273a-4fba-a19c-a3259ac1487a",
				"workspace": "none2none.fmw",
				"elapsedTime": 850,
				"cpuTime": 697,
				"cpuPercent": 82.0,
				"peakMemoryUsage": 5190544,
				"lineCount": 70,
				"warningCount": 0,
				"errorCount": 2
			}
		]
	}`

	responseV4SingleJob := `{
		"id": 999,
		"description": "",
		"engineHost": "10f259e906e5",
		"engineName": "10f259e906e5",
		"repository": "Test",
		"queue": "Default",
		"queueType": "Standard",
		"resultDatasetDownloadUrl": "",
		"status": "failure",
		"timeFinished": "2023-11-15T00:42:31Z",
		"timeQueued": "2023-11-15T00:42:30Z",
		"timeStarted": "2023-11-15T00:42:30Z",
		"runtimeUsername": "admin",
		"runtimeUserID": "5d288fb0-273a-4fba-a19c-a3259ac1487a",
		"workspace": "none2none.fmw",
		"elapsedTime": 850,
		"cpuTime": 697,
		"cpuPercent": 82.0,
		"peakMemoryUsage": 5190544,
		"lineCount": 70,
		"warningCount": 0,
		"errorCount": 2
	}`

	responseV4SingleJobOutput := `{
		"offset": 0,
		"limit": 0,
		"totalCount": 1,
		"items": [
			{
				"id": 999,
				"description": "",
				"engineHost": "10f259e906e5",
				"engineName": "10f259e906e5",
				"repository": "Test",
				"queue": "Default",
				"queueType": "Standard",
				"resultDatasetDownloadUrl": "",
				"status": "failure",
				"timeFinished": "2023-11-15T00:42:31Z",
				"timeQueued": "2023-11-15T00:42:30Z",
				"timeStarted": "2023-11-15T00:42:30Z",
				"runtimeUsername": "admin",
				"runtimeUserID": "5d288fb0-273a-4fba-a19c-a3259ac1487a",
				"workspace": "none2none.fmw",
				"elapsedTime": 850,
				"cpuTime": 697,
				"cpuPercent": 82.0,
				"peakMemoryUsage": 5190544,
				"lineCount": 70,
				"warningCount": 0,
				"errorCount": 2
			}
		]
	}`

	customV4HttpServerHandler := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.RawQuery
		w.WriteHeader(http.StatusOK)

		if strings.Contains(query, "status=queued") && strings.Contains(query, "status=running") {
			_, err := w.Write([]byte(responseV4Active))
			require.NoError(t, err)
		}
		if strings.Contains(query, "status=success") && strings.Contains(query, "status=failure") && strings.Contains(query, "status=cancelled") {
			_, err := w.Write([]byte(responseV4Completed))
			require.NoError(t, err)
		}
		if strings.Contains(query, "status=success") && strings.Contains(query, "status=failure") &&
			!strings.Contains(query, "status=cancelled") &&
			!strings.Contains(query, "status=queued") &&
			!strings.Contains(query, "status=running") {
			_, err := w.Write([]byte(responseV4SuccessAndFailure))
			require.NoError(t, err)
		}
	}

	cases := []testCase{
		{
			name:               "unknown flag v4",
			statusCode:         http.StatusOK,
			args:               []string{"jobs", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
			fmeflowBuild:       25300,
		},
		{
			name:         "500 bad status code v4",
			statusCode:   http.StatusInternalServerError,
			wantErrText:  "500 Internal Server Error: ",
			args:         []string{"jobs"},
			fmeflowBuild: 25300,
		},
		{
			name:         "404 bad status code v4",
			statusCode:   http.StatusNotFound,
			wantErrText:  "404 Not Found: ",
			args:         []string{"jobs"},
			fmeflowBuild: 25300,
		},
		{
			name:            "get jobs v4 all jobs",
			statusCode:      http.StatusOK,
			args:            []string{"jobs"},
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*4[\\s]*10f259e906e5[\\s]*running.fmw[\\s]*running[\\s]*5[\\s]*austinApartments.fmw[\\s]*queued[\\s]*1[\\s]*387f74cd4e1f[\\s]*austinApartments.fmw[\\s]*success[\\s]*2[\\s]*10f259e906e5[\\s]*none2none.fmw[\\s]*failure[\\s]*3[\\s]*145929514b24[\\s]*cancelled.fmw[\\s]*cancelled[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customV4HttpServerHandler)),
			fmeflowBuild:    25300,
		},
		{
			name:            "get jobs v4 all jobs explicitly",
			statusCode:      http.StatusOK,
			args:            []string{"jobs", "--all"},
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*4[\\s]*10f259e906e5[\\s]*running.fmw[\\s]*running[\\s]*5[\\s]*austinApartments.fmw[\\s]*queued[\\s]*1[\\s]*387f74cd4e1f[\\s]*austinApartments.fmw[\\s]*success[\\s]*2[\\s]*10f259e906e5[\\s]*none2none.fmw[\\s]*failure[\\s]*3[\\s]*145929514b24[\\s]*cancelled.fmw[\\s]*cancelled[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customV4HttpServerHandler)),
			fmeflowBuild:    25300,
		},
		{
			name:            "get jobs v4 queued jobs",
			statusCode:      http.StatusOK,
			args:            []string{"jobs", "--queued"},
			body:            responseV4Queued,
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*5[\\s]*austinApartments.fmw[\\s]*queued[\\s]*$",
			fmeflowBuild:    25300,
		},
		{
			name:            "get jobs v4 success and failure",
			statusCode:      http.StatusOK,
			args:            []string{"jobs", "--success", "--failure"},
			body:            responseV4SuccessAndFailure,
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*1[\\s]*387f74cd4e1f[\\s]*austinApartments.fmw[\\s]*success[\\s]*2[\\s]*10f259e906e5[\\s]*none2none.fmw[\\s]*failure[\\s]*$",
			fmeflowBuild:    25300,
		},
		{
			name:           "get jobs v4 by repository",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--repository", "Samples"},
			wantFormParams: map[string]string{"repository": "Samples"},
			body:           responseV4Active,
			fmeflowBuild:   25300,
		},
		{
			name:           "get jobs v4 by workspace",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--repository", "Samples", "--workspace", "austinApartments.fmw"},
			wantFormParams: map[string]string{"workspace": "austinApartments.fmw", "repository": "Samples"},
			body:           responseV4Active,
			fmeflowBuild:   25300,
		},
		{
			name:           "get jobs v4 by source-type",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--source-type", "source"},
			wantFormParams: map[string]string{"sourceType": "source"},
			body:           responseV4Active,
			fmeflowBuild:   25300,
		},
		{
			name:           "get jobs v4 by engine-name",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--engine-name", "10f259e906e5"},
			wantFormParams: map[string]string{"engineName": "10f259e906e5"},
			body:           responseV4Active,
			fmeflowBuild:   25300,
		},
		{
			name:           "get jobs v4 by source-id",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--source-id", "63f2489a-f3fc-4fa0-8df8-198de602b922", "--source-type", "automations"},
			wantFormParams: map[string]string{"sourceID": "63f2489a-f3fc-4fa0-8df8-198de602b922"},
			body:           responseV4Active,
			fmeflowBuild:   25300,
		},
		{
			name:         "get jobs v4 by source-id no source-type",
			statusCode:   http.StatusOK,
			args:         []string{"jobs", "--source-id", "63f2489a-f3fc-4fa0-8df8-198de602b922"},
			wantErrText:  "required flag(s) \"source-type\" not set",
			body:         responseV4Active,
			fmeflowBuild: 25300,
		},
		{
			name:           "get jobs v4 success and failure json",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--success", "--failure", "--json"},
			body:           responseV4SuccessAndFailure,
			wantOutputJson: responseV4SuccessAndFailure,
			fmeflowBuild:   25300,
		},
		{
			name:            "get jobs v4 success and failure no headers",
			statusCode:      http.StatusOK,
			args:            []string{"jobs", "--success", "--failure", "--no-headers"},
			body:            responseV4SuccessAndFailure,
			wantOutputRegex: "^[\\s]*1[\\s]*387f74cd4e1f[\\s]*austinApartments.fmw[\\s]*success[\\s]*2[\\s]*10f259e906e5[\\s]*none2none.fmw[\\s]*failure[\\s]*$",
			fmeflowBuild:    25300,
		},
		{
			name:           "get jobs v4 success and failure json output type",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--success", "--failure", "--output=json"},
			body:           responseV4SuccessAndFailure,
			wantOutputJson: responseV4SuccessAndFailure,
			fmeflowBuild:   25300,
		},
		{
			name:         "queued and engine-name can't both be specified v4",
			statusCode:   http.StatusOK,
			args:         []string{"jobs", "--queued", "--engine-name", "10f259e906e5"},
			wantErrText:  "if any flags in the group [queued engine-name] are set none of the others can be; [engine-name queued] were all set",
			fmeflowBuild: 25300,
		},
		{
			name:            "get single job v4",
			statusCode:      http.StatusOK,
			args:            []string{"jobs", "--id", "999"},
			body:            responseV4SingleJob,
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*999[\\s]*10f259e906e5[\\s]*none2none.fmw[\\s]*failure[\\s]*$",
			fmeflowBuild:    25300,
		},
		{
			name:           "get single job v4 json",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--id", "999", "--json"},
			body:           responseV4SingleJob,
			wantOutputJson: responseV4SingleJobOutput,
			fmeflowBuild:   25300,
		},
		{
			name:         "get single job v4 does not exist",
			statusCode:   http.StatusNotFound,
			args:         []string{"jobs", "--id", "243"},
			wantErrText:  "404 Not Found: ",
			fmeflowBuild: 25300,
		},
		{
			name:            "get jobs v4 success and failure custom columns",
			statusCode:      http.StatusOK,
			body:            responseV4SuccessAndFailure,
			args:            []string{"jobs", "--success", "--failure", "--output", "custom-columns=CPU:.cpuTime,MEMORY:.peakMemoryUsage"},
			wantOutputRegex: "^[\\s]*CPU[\\s]*MEMORY[\\s]*994[\\s]*(10988808|1\\.0988808e\\+07)[\\s]*697[\\s]*(5190544|5\\.190544e\\+06)[\\s]*$",
			fmeflowBuild:    25300,
		},
		{
			name:            "get jobs v4 completed",
			statusCode:      http.StatusOK,
			body:            responseV4Completed,
			args:            []string{"jobs", "--completed"},
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*1[\\s]*387f74cd4e1f[\\s]*austinApartments.fmw[\\s]*success[\\s]*2[\\s]*10f259e906e5[\\s]*none2none.fmw[\\s]*failure[\\s]*3[\\s]*145929514b24[\\s]*cancelled.fmw[\\s]*cancelled[\\s]*$",
			fmeflowBuild:    25300,
		},
	}

	runTests(cases, t)
}
