package cmd

import (
	"io"
	"net/http"

	"github.com/spf13/viper"
)

func buildFmeServerRequest(endpoint string, method string, body io.Reader) (http.Request, error) {
	// retrieve url and token
	fmeserverUrl := viper.GetString("url")
	fmeserverToken := viper.GetString("token")

	req, err := http.NewRequest(method, fmeserverUrl+endpoint, body)
	req.Header.Set("Authorization", "fmetoken token="+fmeserverToken)
	return *req, err
}
