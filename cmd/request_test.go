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
	responseStatus := `{
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
			_, err := w.Write([]byte(responseStatus))
			require.NoError(t, err)
		}
		if strings.Contains(r.URL.Path, "/fmeapiv4/license/request") {
			w.WriteHeader(http.StatusAccepted)
			_, err := w.Write([]byte(""))
			require.NoError(t, err)

		}
		if strings.Contains(r.URL.Path, "/fmeapiv4/license/request/status") {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(responseStatus))
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
			name:            "request license v3",
			args:            []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--api-version", "v3"},
			wantOutputRegex: "^License Request Successfully sent\\.[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:            "request license v4",
			args:            []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--api-version", "v4"},
			wantOutputRegex: "^License Request Successfully sent\\.[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:            "request license and wait v3",
			statusCode:      http.StatusOK,
			args:            []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--wait", "--api-version", "v3"},
			wantOutputRegex: "^License Request Successfully sent\\.[\\s]*Success! Your FME Server has now been licensed\\.[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:            "request license and wait v4",
			statusCode:      http.StatusOK,
			args:            []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--wait", "--api-version", "v4"},
			wantOutputRegex: "^License Request Successfully sent\\.[\\s]*Success! Your FME Server has now been licensed\\.[\\s]*$",
			httpServer:      httptest.NewServer(http.HandlerFunc(customHttpServerHandler)),
		},
		{
			name:           "request license check form params v3",
			statusCode:     http.StatusAccepted,
			args:           []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--serial-number", "AAAA-AAAA-AAAA", "--company", "Example Inc.", "--industry", "Industry", "--sales-source", "source", "--subscribe-to-updates", "--category", "Category", "--api-version", "v3"},
			wantFormParams: map[string]string{"firstName": "Billy", "lastName": "Bob", "email": "billy.bob@example.com", "serialNumber": "AAAA-AAAA-AAAA", "company": "Example Inc.", "category": "Category", "industry": "Industry", "salesSource": "source", "subscribeToUpdates": "true"},
		},
		{
			name:         "request license check body params v4",
			statusCode:   http.StatusAccepted,
			args:         []string{"license", "request", "--first-name", "Billy", "--last-name", "Bob", "--email", "billy.bob@example.com", "--serial-number", "AAAA-AAAA-AAAA", "--company", "Example Inc.", "--industry", "Industry", "--subscribe-to-updates", "--api-version", "v4"},
			wantBodyJson: `{"firstName":"Billy","lastName":"Bob","email":"billy.bob@example.com","serialNumber":"AAAA-AAAA-AAAA","company":"Example Inc.","industry":"Industry","subscribeToUpdates":true}`,
		},
	}

	runTests(cases, t)

}
