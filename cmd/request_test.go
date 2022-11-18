package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLicenseRequest(t *testing.T) {
	// standard responses for v3
	responseV3Status := `{
		"message": "Success! Your FME Server has now been licensed.",
		"status": "SUCCESS"
	  }`

	customHttpServerHandler := func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/fmerest/v3/licensing/request") {
			w.WriteHeader(http.StatusAccepted)
			_, err := w.Write([]byte(""))
			require.NoError(t, err)

		}
		if strings.Contains(r.URL.Path, "/fmerest/v3/licensing/request/status") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseV3Status))
			require.NoError(t, err)
		}

	}

	cases := []testCase{
		{
			name:               "unknown flag",
			statusCode:         http.StatusAccepted,
			args:               []string{"license", "request", "--badflag"},
			wantErrOutputRegex: "unknown flag: --badflag",
		},
		{
			name:        "500 bad status code",
			statusCode:  http.StatusInternalServerError,
			wantErrText: "500 Internal Server Error",
			args:        []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com"},
		},
		{
			name:        "404 bad status code",
			statusCode:  http.StatusNotFound,
			wantErrText: "404 Not Found",
			args:        []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com"},
		},
		{
			name:        "request license missing email flag",
			statusCode:  http.StatusAccepted,
			args:        []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob"},
			wantErrText: "required flag(s) \"email\" not set",
		},
		{
			name:        "request license missing last name flag",
			statusCode:  http.StatusAccepted,
			args:        []string{"license", "request", "--first-name", "Billy", "--email", "billy.bob@example.com"},
			wantErrText: "required flag(s) \"last-name\" not set",
		},
		{
			name:        "request license missing first name flag",
			statusCode:  http.StatusAccepted,
			args:        []string{"license", "request", "--last-name", "Bob", "--email", "billy.bob@example.com"},
			wantErrText: "required flag(s) \"first-name\" not set",
		},
		{
			name:            "request license",
			args:            []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com"},
			wantOutputRegex: "^License Request Successfully sent\\.[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:            "request license and wait",
			statusCode:      http.StatusOK,
			args:            []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--wait"},
			wantOutputRegex: "^License Request Successfully sent\\.[\\s]*Success! Your FME Server has now been licensed\\.[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:           "request license check form params",
			statusCode:     http.StatusAccepted,
			args:           []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--serial-number", "AAAA-AAAA-AAAA", "--company", "Example Inc.", "--industry", "Industry", "--sales-source", "source", "--subscribe-to-updates", "--category", "Category"},
			wantFormParams: map[string]string{"firstName": "Billy", "lastName": "Bob", "email": "billy.bob@example.com", "serialNumber": "AAAA-AAAA-AAAA", "company": "Example Inc.", "category": "Category", "industry": "Industry", "salesSource": "source", "subscribeToUpdates": "true"},
		},
	}

	runTests(cases, t)

}
