package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

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

// since the JSON for published parameters has subtypes, we need to implement this ourselves
func (f *Job) UnmarshalJSON(b []byte) error {
	type job Job
	err := json.Unmarshal(b, (*job)(f))
	if err != nil {
		return err
	}

	for _, raw := range f.RawPublishedParameters {
		//var v map[string]interface{}
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
			//fmt.Println(string(k))
			//fmt.Println(string(v))
			//fmt.Println("----")
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
