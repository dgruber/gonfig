package gonfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ConfigServerResponse encodes is the expected JSON response of the
// Configuration Server.
type ConfigServerResponse struct {
	Name     string   `json:"name"`
	Profiles []string `json:"profiles"`
	Label    string   `json:"label"`
	Version  string   `json:"version"`
	// State           interface{} `json:"state"`
	PropertySources []struct {
		// Name is the path to the github repository for example
		Name string `json:"name"`
		// Source contains the actual configuration (float64 for ints)
		Source map[string]interface{} `json:"source"` // Rest of the fields should go here.
	} `json:"propertySources"`
}

func makeRequest(client *http.Client, url string) (map[string]interface{}, error) {
	resp, errGet := client.Get(url)
	if errGet != nil {
		return nil, fmt.Errorf("Error during getting configuration from URL %s: %s", url, errGet.Error())
	}
	defer resp.Body.Close()

	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, fmt.Errorf("Error during reading response body from config server: %s", errRead.Error())
	}

	var response ConfigServerResponse
	errDecode := json.Unmarshal(body, &response)
	if errDecode != nil {
		return nil, fmt.Errorf("Error during deconding response from config server: %s", errDecode.Error())
	} else if len(response.PropertySources) == 0 {
		return nil, fmt.Errorf("Response from config server has zero length: %s", body)
	}
	return response.PropertySources[0].Source, nil
}
