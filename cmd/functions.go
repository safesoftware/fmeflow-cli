package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/viper"
)

// define the default style for all tables that are output
var defaultStyle = table.Style{
	Name: "Borderless",
	Box:  table.StyleBoxLight,
	Format: table.FormatOptions{
		Footer: text.FormatUpper,
		Header: text.FormatUpper,
		Row:    text.FormatDefault,
	},
	Options: table.Options{
		DrawBorder:      false,
		SeparateColumns: false,
		SeparateFooter:  false,
		SeparateHeader:  false,
		SeparateRows:    false,
	},
}

func buildFmeServerRequest(endpoint string, method string, body io.Reader) (http.Request, error) {
	// retrieve url and token
	fmeserverUrl := viper.GetString("url")
	fmeserverToken := viper.GetString("token")

	req, err := http.NewRequest(method, fmeserverUrl+endpoint, body)
	req.Header.Set("Authorization", "fmetoken token="+fmeserverToken)
	return *req, err
}

// since the JSON for published parameters has subtypes, we need to implement this ourselves
func (f *Job) UnmarshalJSON(b []byte) error {
	type job Job
	err := json.Unmarshal(b, (*job)(f))
	if err != nil {
		return err
	}

	for _, raw := range f.RawPublishedParameters {
		data := make(map[string]json.RawMessage)
		err = json.Unmarshal(raw, &data)
		if err != nil {
			return err
		}

		var i interface{}
		for k, v := range data {
			if k == "value" {
				if strings.HasPrefix(string(v), "[") {
					i = &ListParameter{}
				} else {
					i = &SimpleParameter{}
				}

			}
		}

		if i != nil {
			err = json.Unmarshal(raw, i)
			if err != nil {
				return err
			}
			f.PublishedParameters = append(f.PublishedParameters, i)
		}
	}
	return nil
}

func (f *Job) MarshalJSON() ([]byte, error) {

	type job Job
	if f.PublishedParameters != nil {
		for _, v := range f.PublishedParameters {
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			f.RawPublishedParameters = append(f.RawPublishedParameters, b)
		}
	}
	return json.Marshal((*job)(f))
}

func prettyPrintJSON(s []byte) (string, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, s, "", "  ")
	if err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
