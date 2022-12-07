package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJobs(t *testing.T) {
	// standard responses for v3
	responseV3Completed := `{
		"offset":0,
		"limit":0,
		"totalCount":2,
		"items":[
			{
				"request":{
					"TMDirectives":{
						"rtc":false,
						"ttc":-1,
						"tag":"Default",
						"ttl":-1
					},
					"NMDirectives":{
						"directives":[
							{
								"name":"logFullPath",
								"value":"{logHome}/{@logFileName(id)}"
							}
						]
					}
				},
				"timeDelivered":"2022-11-08T21:52:50Z",
				"workspace":"austinApartments.fmw",
				"numErrors":0,
				"numLines":117,
				"engineHost":"387f74cd4e1f",
				"timeQueued":"2022-11-08T21:52:04Z",
				"cpuPct":88.9087656529517,
				"description":"",
				"timeStarted":"2022-11-08T21:52:49Z",
				"repository":"Samples",
				"userName":"admin",
				"result":{
					"timeRequested":"2022-11-08T21:52:04Z",
					"requesterResultPort":-1,
					"numFeaturesOutput":49,
					"requesterHost":"192.168.16.5",
					"timeStarted":"2022-11-08T21:52:49Z",
					"id":3,
					"timeFinished":"2022-11-08T21:52:50Z",
					"priority":-1,
					"statusMessage":"Translation Successful",
					"status":"SUCCESS"
				},
				"cpuTime":994,
				"id":3,
				"timeFinished":"2022-11-08T21:52:50Z",
				"engineName":"387f74cd4e1f",
				"numWarnings":0,
				"timeSubmitted":"2022-11-08T21:52:04Z",
				"elapsedTime":1118,
				"peakMemUsage":10988808,
				"status":"SUCCESS"
			},
			{
				"request":{
				   "TMDirectives":{
					  "rtc":false,
					  "ttc":-1,
					  "tag":"Default",
					  "ttl":-1
				   },
				   "NMDirectives":{
					  "directives":[
						 {
							"name":"logFullPath",
							"value":"{logHome}/{@logFileName(id)}"
						 }
					  ]
				   }
				},
				"timeDelivered":"2022-11-15T00:42:31Z",
				"workspace":"none2none.fmw",
				"numErrors":2,
				"numLines":70,
				"engineHost":"10f259e906e5",
				"timeQueued":"2022-11-15T00:42:30Z",
				"cpuPct":82,
				"description":"",
				"timeStarted":"2022-11-15T00:42:30Z",
				"repository":"Test",
				"userName":"admin",
				"result":{
				   "timeRequested":"2022-11-15T00:42:30Z",
				   "requesterResultPort":-1,
				   "numFeaturesOutput":0,
				   "requesterHost":"192.168.32.4",
				   "timeStarted":"2022-11-15T00:42:30Z",
				   "id":2,
				   "timeFinished":"2022-11-15T00:42:31Z",
				   "priority":-1,
				   "statusMessage":"Terminator: Termination Message: 'Translation Terminated'",
				   "status":"FME_FAILURE"
				},
				"cpuTime":697,
				"id":2,
				"timeFinished":"2022-11-15T00:42:31Z",
				"engineName":"10f259e906e5",
				"numWarnings":0,
				"timeSubmitted":"2022-11-15T00:42:30Z",
				"elapsedTime":850,
				"peakMemUsage":5190544,
				"status":"FME_FAILURE"
			}
		]
	}`
	responseV3Active := `{
		"offset":0,
		"limit":0,
		"totalCount":2,
		"items":[
		   {
			  "request":{
				 "TMDirectives":{
					"rtc":false,
					"ttc":-1,
					"tag":"Default",
					"ttl":-1
				 },
				 "NMDirectives":{
					"directives":[
					   {
						  "name":"logFullPath",
						  "value":"{logHome}/{@logFileName(id)}"
					   }
					]
				 }
			  },
			  "timeDelivered":"0001-01-01T00:00:00Z",
			  "workspace":"running.fmw",
			  "numErrors":0,
			  "numLines":0,
			  "engineHost":"10f259e906e5",
			  "timeQueued":"2022-11-15T01:22:14Z",
			  "cpuPct":0,
			  "description":"",
			  "timeStarted":"2022-11-15T01:22:14Z",
			  "repository":"Test",
			  "userName":"admin",
			  "result":{
				 "timeRequested":"0001-01-01T00:00:00Z",
				 "requesterResultPort":0,
				 "numFeaturesOutput":0,
				 "requesterHost":"",
				 "timeStarted":"0001-01-01T00:00:00Z",
				 "id":0,
				 "timeFinished":"0001-01-01T00:00:00Z",
				 "priority":0,
				 "statusMessage":"",
				 "status":""
			  },
			  "cpuTime":0,
			  "id":4,
			  "timeFinished":"0001-01-01T00:00:00Z",
			  "engineName":"10f259e906e5",
			  "numWarnings":0,
			  "timeSubmitted":"2022-11-15T01:22:14Z",
			  "elapsedTime":0,
			  "peakMemUsage":0,
			  "status":"PULLED"
		   },
		   {
			  "request":{
				 "TMDirectives":{
					"rtc":false,
					"ttc":-1,
					"tag":"Default",
					"ttl":-1
				 },
				 "NMDirectives":{
					"directives":[
					   {
						  "name":"logFullPath",
						  "value":"{logHome}/{@logFileName(id)}"
					   }
					]
				 }
			  },
			  "timeDelivered":"0001-01-01T00:00:00Z",
			  "workspace":"austinApartments.fmw",
			  "numErrors":0,
			  "numLines":0,
			  "engineHost":"",
			  "timeQueued":"2022-11-08T21:52:04Z",
			  "cpuPct":0,
			  "description":"",
			  "timeStarted":"0001-01-01T00:00:00Z",
			  "repository":"Samples",
			  "userName":"admin",
			  "result":{
				 "timeRequested":"0001-01-01T00:00:00Z",
				 "requesterResultPort":0,
				 "numFeaturesOutput":0,
				 "requesterHost":"",
				 "timeStarted":"0001-01-01T00:00:00Z",
				 "id":0,
				 "timeFinished":"0001-01-01T00:00:00Z",
				 "priority":0,
				 "statusMessage":"",
				 "status":""
			  },
			  "cpuTime":0,
			  "id":1,
			  "timeFinished":"0001-01-01T00:00:00Z",
			  "engineName":"",
			  "numWarnings":0,
			  "timeSubmitted":"2022-11-08T21:52:04Z",
			  "elapsedTime":0,
			  "peakMemUsage":0,
			  "status":"QUEUED"
		   }
		]
	}`

	responseV3Running := `{
		"offset":0,
		"limit":0,
		"totalCount":1,
		"items":[
		   {
			  "request":{
				 "TMDirectives":{
					"rtc":false,
					"ttc":-1,
					"tag":"Default",
					"ttl":-1
				 },
				 "NMDirectives":{
					"directives":[
					   {
						  "name":"logFullPath",
						  "value":"{logHome}/{@logFileName(id)}"
					   }
					]
				 }
			  },
			  "timeDelivered":"0001-01-01T00:00:00Z",
			  "workspace":"running.fmw",
			  "numErrors":0,
			  "numLines":0,
			  "engineHost":"10f259e906e5",
			  "timeQueued":"2022-11-15T01:22:14Z",
			  "cpuPct":0,
			  "description":"",
			  "timeStarted":"2022-11-15T01:22:14Z",
			  "repository":"Test",
			  "userName":"admin",
			  "result":{
				 "timeRequested":"0001-01-01T00:00:00Z",
				 "requesterResultPort":0,
				 "numFeaturesOutput":0,
				 "requesterHost":"",
				 "timeStarted":"0001-01-01T00:00:00Z",
				 "id":0,
				 "timeFinished":"0001-01-01T00:00:00Z",
				 "priority":0,
				 "statusMessage":"",
				 "status":""
			  },
			  "cpuTime":0,
			  "id":4,
			  "timeFinished":"0001-01-01T00:00:00Z",
			  "engineName":"10f259e906e5",
			  "numWarnings":0,
			  "timeSubmitted":"2022-11-15T01:22:14Z",
			  "elapsedTime":0,
			  "peakMemUsage":0,
			  "status":"PULLED"
		   }
		]
	}`

	responseV3Queued := `{
		"offset":0,
		"limit":0,
		"totalCount":1,
		"items":[
		   {
			  "request":{
				 "TMDirectives":{
					"rtc":false,
					"ttc":-1,
					"tag":"Default",
					"ttl":-1
				 },
				 "NMDirectives":{
					"directives":[
					   {
						  "name":"logFullPath",
						  "value":"{logHome}/{@logFileName(id)}"
					   }
					]
				 }
			  },
			  "timeDelivered":"0001-01-01T00:00:00Z",
			  "workspace":"austinApartments.fmw",
			  "numErrors":0,
			  "numLines":0,
			  "engineHost":"",
			  "timeQueued":"2022-11-08T21:52:04Z",
			  "cpuPct":0,
			  "description":"",
			  "timeStarted":"0001-01-01T00:00:00Z",
			  "repository":"Samples",
			  "userName":"admin",
			  "result":{
				 "timeRequested":"0001-01-01T00:00:00Z",
				 "requesterResultPort":0,
				 "numFeaturesOutput":0,
				 "requesterHost":"",
				 "timeStarted":"0001-01-01T00:00:00Z",
				 "id":0,
				 "timeFinished":"0001-01-01T00:00:00Z",
				 "priority":0,
				 "statusMessage":"",
				 "status":""
			  },
			  "cpuTime":0,
			  "id":1,
			  "timeFinished":"0001-01-01T00:00:00Z",
			  "engineName":"",
			  "numWarnings":0,
			  "timeSubmitted":"2022-11-08T21:52:04Z",
			  "elapsedTime":0,
			  "peakMemUsage":0,
			  "status":"QUEUED"
		   }
		]
	 }`

	customHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if strings.Contains(r.URL.Path, "active") {
			_, err := w.Write([]byte(responseV3Active))
			require.NoError(t, err)
		}
		if strings.Contains(r.URL.Path, "completed") {
			_, err := w.Write([]byte(responseV3Completed))
			require.NoError(t, err)
		}
		if strings.Contains(r.URL.Path, "queued") {
			_, err := w.Write([]byte(responseV3Queued))
			require.NoError(t, err)
		}
		if strings.Contains(r.URL.Path, "running") {
			_, err := w.Write([]byte(responseV3Running))
			require.NoError(t, err)
		}

	}

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"jobs", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"jobs"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"jobs"},
		},
		{
			name:            "get jobs table output",
			statusCode:      http.StatusOK,
			args:            []string{"jobs"},
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*4[\\s]*10f259e906e5[\\s]*running.fmw[\\s]*PULLED[\\s]*1[\\s]*austinApartments.fmw[\\s]*QUEUED[\\s]*3[\\s]*387f74cd4e1f[\\s]*austinApartments.fmw[\\s]*SUCCESS[\\s]*2[\\s]*10f259e906e5[\\s]*none2none.fmw[\\s]*FME_FAILURE[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:            "get jobs all table output",
			statusCode:      http.StatusOK,
			args:            []string{"jobs", "--all"},
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*4[\\s]*10f259e906e5[\\s]*running.fmw[\\s]*PULLED[\\s]*1[\\s]*austinApartments.fmw[\\s]*QUEUED[\\s]*3[\\s]*387f74cd4e1f[\\s]*austinApartments.fmw[\\s]*SUCCESS[\\s]*2[\\s]*10f259e906e5[\\s]*none2none.fmw[\\s]*FME_FAILURE[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:            "get jobs running",
			statusCode:      http.StatusOK,
			body:            responseV3Running,
			args:            []string{"jobs", "--running"},
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*4[\\s]*10f259e906e5[\\s]*running.fmw[\\s]*PULLED[\\s]*$",
		},
		{
			name:            "get jobs active",
			statusCode:      http.StatusOK,
			body:            responseV3Active,
			args:            []string{"jobs", "--active"},
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*4[\\s]*10f259e906e5[\\s]*running.fmw[\\s]*PULLED[\\s]*1[\\s]*austinApartments.fmw[\\s]*QUEUED[\\s]*$",
		},
		{
			name:            "get jobs completed",
			statusCode:      http.StatusOK,
			body:            responseV3Completed,
			args:            []string{"jobs", "--completed"},
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*3[\\s]*387f74cd4e1f[\\s]*austinApartments.fmw[\\s]*SUCCESS[\\s]*2[\\s]*10f259e906e5[\\s]*none2none.fmw[\\s]*FME_FAILURE[\\s]*$",
		},
		{
			name:            "get jobs queued",
			statusCode:      http.StatusOK,
			body:            responseV3Queued,
			args:            []string{"jobs", "--queued"},
			wantOutputRegex: "^[\\s]*JOB ID[\\s]*ENGINE NAME[\\s]*WORKSPACE[\\s]*STATUS[\\s]*1[\\s]*austinApartments.fmw[\\s]*QUEUED[\\s]*$",
		},
		{
			name:            "get jobs queued no headers",
			statusCode:      http.StatusOK,
			body:            responseV3Queued,
			args:            []string{"jobs", "--queued", "--no-headers"},
			wantOutputRegex: "^[\\s]*1[\\s]*austinApartments.fmw[\\s]*QUEUED[\\s]*$",
		},
		{
			name:           "get jobs queued json",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--queued", "--json"},
			body:           responseV3Queued,
			wantOutputJson: responseV3Queued,
		},
		{
			name:           "get jobs queued json output type",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--queued", "--output=json"},
			body:           responseV3Queued,
			wantOutputJson: responseV3Queued,
		},
		{
			name:        "workspace flag requires repository",
			statusCode:  http.StatusOK,
			args:        []string{"jobs", "--workspace", "austinApartments.fmw"},
			wantErrText: "required flag(s) \"repository\" not set",
			body:        responseV3Completed,
		},
		{
			name:        "queued and active can't both be specified",
			statusCode:  http.StatusOK,
			args:        []string{"jobs", "--queued", "--active"},
			wantErrText: "if any flags in the group [queued active] are set none of the others can be; [active queued] were all set",
			body:        responseV3Completed,
		},
		{
			name:        "running and active can't both be specified",
			statusCode:  http.StatusOK,
			args:        []string{"jobs", "--running", "--active"},
			wantErrText: "if any flags in the group [running active] are set none of the others can be; [active running] were all set",
			body:        responseV3Completed,
		},
		{
			name:           "get jobs by repository",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--repository", "Samples"},
			wantFormParams: map[string]string{"repository": "Samples"},
			body:           responseV3Completed,
		},
		{
			name:           "get jobs by workspace",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--repository", "Samples", "--workspace", "austinApartments.fmw"},
			wantFormParams: map[string]string{"workspace": "austinApartments.fmw", "repository": "Samples"},
			body:           responseV3Completed,
		},
		{
			name:           "get jobs by source id",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--source-id", "some-source-id"},
			wantFormParams: map[string]string{"sourceID": "some-source-id"},
			body:           responseV3Completed,
		},
		{
			name:           "get jobs by user",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--user-name", "admin"},
			wantFormParams: map[string]string{"userName": "admin"},
			body:           responseV3Completed,
		},
		{
			name:           "get jobs by source-type",
			statusCode:     http.StatusOK,
			args:           []string{"jobs", "--source-type", "source"},
			wantFormParams: map[string]string{"sourceType": "source"},
			body:           responseV3Completed,
		},
		{
			name:            "get jobs completed custom columns",
			statusCode:      http.StatusOK,
			body:            responseV3Completed,
			args:            []string{"jobs", "--completed", "--output", "custom-columns=CPU:.cpuTime,FEATURES OUTPUT:.result.numFeaturesOutput"},
			wantOutputRegex: "^[\\s]*CPU[\\s]*FEATURES OUTPUT[\\s]*994[\\s]*49[\\s]*697[\\s]*0[\\s]*$",
		},
	}

	runTests(cases, t)

}
