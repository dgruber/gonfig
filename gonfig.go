package gonfig

import (
	"crypto/tls"
	"fmt"
	"github.com/cloudfoundry-community/go-cfenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"os"
)

// URL represents the particles for building the URL to access the configuration
// from the Pivotal Cloud Foundry configuration server.
type URL struct {
	URI     string
	App     string
	Profile string
	Label   string
}

// Credentials are used in order to access the PCF Config Server with oauth2.
type Credentials struct {
	AccessTokenURI string
	ClientID       string
	ClientSecret   string
	URL            URL
}

// GetConfigurationFromServer requests the current configuration from the Configuration Server
// using the given Credentials.
func (c *Credentials) GetConfigurationFromServer() (map[string]interface{}, error) {
	var client *http.Client
	var url string

	if os.Getenv("gonfig_testing") == "1" {
		client = &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
		url = c.URL.URI
	} else {
		conf := &clientcredentials.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			TokenURL:     c.AccessTokenURI,
		}
		client = oauth2.NewClient(oauth2.NoContext, conf.TokenSource(oauth2.NoContext))
		url = fmt.Sprintf("%s/%s/%s/%s", c.URL.URI, c.URL.App, c.URL.Profile, c.URL.Label)
	}
	return makeRequest(client, url)
}

// GetConfigServerCredentialsFromEnv returns oauth2 and other parameters (app name,
// space name) out of the PCF enviornment variables. They are required for accessing
// the PCF Configuration Server.
func GetConfigServerCredentialsFromEnv() (*Credentials, error) {
	env, err := cfenv.Current()
	if err != nil {
		return nil, fmt.Errorf("Error during fetching current configuration: %s", err.Error())
	}

	services, err := env.Services.WithLabel("p-config-server")
	if err != nil {
		return nil, err
	}

	var credentials Credentials

	credentials.AccessTokenURI, _ = services[0].CredentialString("access_token_uri")
	if credentials.AccessTokenURI == "" {
		return nil, fmt.Errorf("access_token_uri not found in credentials")
	}

	credentials.ClientID, _ = services[0].CredentialString("client_id")
	if credentials.ClientID == "" {
		return nil, fmt.Errorf("client_id not found in credentials")
	}

	credentials.ClientSecret, _ = services[0].CredentialString("client_secret")
	if credentials.ClientSecret == "" {
		return nil, fmt.Errorf("client_secret not found in credentials")
	}

	uri, _ := services[0].CredentialString("uri")
	if uri == "" {
		return nil, fmt.Errorf("uri not found in credentials")
	}

	credentials.URL = URL{
		URI:     uri,
		App:     env.Name,
		Profile: env.SpaceName,
		Label:   "master",
	}

	return &credentials, nil
}

// FetchConfig returns the configuration given by the PCF Config Server which is bound
// as service to the app.
func FetchConfig() (map[string]interface{}, error) {
	return FetchConfigByLabel("master")
}

// FetchConfigByLabel returns the configuration from the PCF Config Server for a specifc
// label. The default label is "master" which is used by FetchConfig(). The label represents
// for a git configuration typically a branch name.
func FetchConfigByLabel(label string) (map[string]interface{}, error) {
	credentials, err := GetConfigServerCredentialsFromEnv()
	if err != nil {
		return nil, err
	}
	credentials.URL.Label = label
	return credentials.GetConfigurationFromServer()
}
