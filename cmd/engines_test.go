package cmd

import (
	"net/http"
	"testing"
)

func TestEngines(t *testing.T) {
	// standard responses for v3 and v4
	responseV3 := `{
		"offset": -1,
		"limit": -1,
		"totalCount": 1,
		"items": [
		  {
			"hostName": "387f74cd4e1f",
			"assignedQueues": [
			  "Default"
			],
			"resultFailureCount": 0,
			"instanceName": "387f74cd4e1f",
			"registrationProperties": [
			  "Standard",
			  "387f74cd4e1f",
			  "387f74cd4e1f",
			  "23166",
			  "linux-x64"
			],
			"engineManagerNodeName": "fmeservercore",
			"maxTransactionResultFailure": 10,
			"type": "STANDARD",
			"buildNumber": 23166,
			"platform": "linux-x64",
			"resultSuccessCount": 0,
			"maxTransactionResultSuccess": 100,
			"assignedStreams": [],
			"transactionPort": 40059,
			"currentJobID": -1
		  }
		]
	  }`

	responseV3FourEngines := `{
		"offset": -1,
		"limit": -1,
		"totalCount": 4,
		"items": [
		  {
			"hostName": "eaf909ea8a98",
			"assignedQueues": [
			  "Default"
			],
			"resultFailureCount": 0,
			"instanceName": "eaf909ea8a98",
			"registrationProperties": [
			  "Standard",
			  "eaf909ea8a98",
			  "eaf909ea8a98",
			  "23166",
			  "linux-x64"
			],
			"engineManagerNodeName": "fmeservercore",
			"maxTransactionResultFailure": 10,
			"type": "STANDARD",
			"buildNumber": 23166,
			"platform": "linux-x64",
			"resultSuccessCount": 0,
			"maxTransactionResultSuccess": 100,
			"assignedStreams": [],
			"transactionPort": 40935,
			"currentJobID": -1
		  },
		  {
			"hostName": "10f259e906e5",
			"assignedQueues": [
			  "Default"
			],
			"resultFailureCount": 0,
			"instanceName": "10f259e906e5",
			"registrationProperties": [
			  "Standard",
			  "10f259e906e5",
			  "10f259e906e5",
			  "23166",
			  "linux-x64"
			],
			"engineManagerNodeName": "fmeservercore",
			"maxTransactionResultFailure": 10,
			"type": "STANDARD",
			"buildNumber": 23166,
			"platform": "linux-x64",
			"resultSuccessCount": 0,
			"maxTransactionResultSuccess": 100,
			"assignedStreams": [],
			"transactionPort": 36883,
			"currentJobID": -1
		  },
		  {
			"hostName": "fe1da0f5536d",
			"assignedQueues": [
			  "Default"
			],
			"resultFailureCount": 0,
			"instanceName": "fe1da0f5536d",
			"registrationProperties": [
			  "Standard",
			  "fe1da0f5536d",
			  "fe1da0f5536d",
			  "23166",
			  "linux-x64"
			],
			"engineManagerNodeName": "fmeservercore",
			"maxTransactionResultFailure": 10,
			"type": "STANDARD",
			"buildNumber": 23166,
			"platform": "linux-x64",
			"resultSuccessCount": 0,
			"maxTransactionResultSuccess": 100,
			"assignedStreams": [],
			"transactionPort": 44089,
			"currentJobID": -1
		  },
		  {
			"hostName": "005cafdec613",
			"assignedQueues": [
			  "Default"
			],
			"resultFailureCount": 0,
			"instanceName": "005cafdec613",
			"registrationProperties": [
			  "Standard",
			  "005cafdec613",
			  "005cafdec613",
			  "23166",
			  "linux-x64"
			],
			"engineManagerNodeName": "fmeservercore",
			"maxTransactionResultFailure": 10,
			"type": "STANDARD",
			"buildNumber": 23166,
			"platform": "linux-x64",
			"resultSuccessCount": 0,
			"maxTransactionResultSuccess": 100,
			"assignedStreams": [],
			"transactionPort": 44795,
			"currentJobID": -1
		  }
		]
	  }`

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusOK,
			args:               []string{"engines", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"engines"},
		},
		{
			name:        "422 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"engines"},
		},
		{
			name:            "get engines",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"engines"},
			wantOutputRegex: "NAME[\\s]*HOST[\\s]*BUILD[\\s]*PLATFORM[\\s]*TYPE[\\s]*CURRENT JOB ID[\\s]*REGISTRATION PROPERTIES[\\s]*QUEUES[\\s]*[\\s]*387f74cd4e1f[\\s]*387f74cd4e1f[\\s]*23166[\\s]*linux-x64[\\s]*STANDARD[\\s]*-1[\\s]*\\[Standard 387f74cd4e1f 387f74cd4e1f 23166 linux-x64\\][\\s]*\\[Default\\]",
		},
		{
			name:            "get engines no headers",
			statusCode:      http.StatusOK,
			body:            responseV3,
			args:            []string{"engines", "--no-headers"},
			wantOutputRegex: "[\\s]*387f74cd4e1f[\\s]*387f74cd4e1f[\\s]*23166[\\s]*linux-x64[\\s]*STANDARD[\\s]*-1[\\s]*\\[Standard 387f74cd4e1f 387f74cd4e1f 23166 linux-x64\\][\\s]*\\[Default\\]",
		},
		{
			name:           "get engines json",
			statusCode:     http.StatusOK,
			args:           []string{"engines", "--json"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
		{
			name:           "get engines json via output type",
			statusCode:     http.StatusOK,
			args:           []string{"engines", "--output=json"},
			body:           responseV3,
			wantOutputJson: responseV3,
		},
		{
			name:            "get engines count",
			statusCode:      http.StatusOK,
			body:            responseV3FourEngines,
			args:            []string{"engines", "--count"},
			wantOutputRegex: "4",
		},
		{
			name:            "get engines custom columns",
			statusCode:      http.StatusOK,
			body:            responseV3FourEngines,
			args:            []string{"engines", "--output=custom-columns=ENGINEMANAGER:$.engineManagerNodeName,TRANSACTIONPORT:$.transactionPort,CURRENTJOB:$.currentJobID"},
			wantOutputRegex: "[\\s]*ENGINEMANAGER[\\s]*TRANSACTIONPORT[\\s]*CURRENTJOB[\\s]*fmeservercore[\\s]*40935[\\s]*-1[\\s]*fmeservercore[\\s]*36883[\\s]*-1[\\s]*fmeservercore[\\s]*44089[\\s]*-1[\\s]*fmeservercore[\\s]*44795[\\s]*-1",
		},
	}

	runTests(cases, t)

}
