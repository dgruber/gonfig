package gonfig

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

var testVCAPApp = `
{"instance_id":"451f045fd16427bb99c895a2649b7b2a",
"instance_index":0,
"host":"0.0.0.0",
"port":61857,
"started_at":"2013-08-12 00:05:29 +0000",
"started_at_timestamp":1376265929,
"start":"2013-08-12 00:05:29 +0000",
"state_timestamp":1376265929,
"limits":{"mem":512,"disk":1024,"fds":16384},
"application_version":"c1063c1c-40b9-434e-a797-db240b587d32",
"application_name":"styx-james",
"application_uris":["styx-james.a1-app.cf-app.com"],
"version":"c1063c1c-40b9-434e-a797-db240b587d32",
"name":"styx-james",
"space_id":"3e0c28c5-6d9c-436b-b9ee-1f4326e54d05",
"space_name":"jdk",
"uris":["styx-james.a1-app.cf-app.com"],"users":null}
`

var testVCAPwithoutConfig = `
{
  "elephantsql": [
    {
      "name": "elephantsql-c6c60",
      "label": "elephantsql",
      "tags": [
        "postgres",
        "postgresql",
        "relational"
      ],
      "plan": "turtle",
      "credentials": {
        "uri": "postgres://exampleuser:examplepass@babar.elephantsql.com:5432/exampleuser"
      }
    }
  ],
  "sendgrid": [
    {
      "name": "mysendgrid",
      "label": "sendgrid",
      "tags": [
        "smtp"
      ],
      "plan": "free",
      "credentials": {
        "hostname": "smtp.sendgrid.net",
        "username": "QvsXMbJ3rK",
        "password": "HCHMOYluTv"
      }
    }
  ]
}
`

var testVCAPwithConfig = `
{
  "p-config-server": [
   {
    "credentials": {
     "access_token_uri": "https://p-spring-cloud-services.uaa.cf.wise.com/oauth/token",
     "client_id": "p-config-server-c4a56a3d-9507-4c2f-9cd1-f858dbf9e11c",
     "client_secret": "9aGx9K5Vx0cM",
     "uri": "https://config-51711835-4626-4823-b5a1-e5d91012f3f2.apps.wise.com"
    },
    "label": "p-config-server",
    "name": "config-server",
    "plan": "standard",
    "tags": [
     "configuration",
     "spring-cloud"
    ]
   }
  ]
 }
 `

var testConfigServerResponse = `
{
  "name":"gonfig",
  "profiles":["dgruber-dev"],
  "label":"master",
  "version":"77091415ec22c45e8f76407e2b30c226b33271",
  "state":null,
  "propertySources":[
    {
      "name":"https://github.com/dgruber/sample-config/gonfig.yml",
      "source":{"resolutionX":640,"resolutionY":480}
    }
  ]
}
`

var _ = Describe("Gonfig", func() {

	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", testVCAPwithConfig)
		os.Setenv("VCAP_APPLICATION", testVCAPApp)
		os.Setenv("gonfig_testing", "0")
	})

	Describe("GetConfigServerCredentialsFromEnv()", func() {
		Context("Correct VCAP set for one config server", func() {
			It("must parse the credentials correctly", func() {
				c, err := GetConfigServerCredentialsFromEnv()
				Expect(err).To(BeNil())
				Expect(c.URL.URI).To(BeEquivalentTo("https://config-51711835-4626-4823-b5a1-e5d91012f3f2.apps.wise.com"))
				Expect(c.URL.App).To(BeEquivalentTo("styx-james"))
				Expect(c.URL.Label).To(BeEquivalentTo("master"))
				Expect(c.URL.Profile).To(BeEquivalentTo("jdk"))
				Expect(c.ClientSecret).To(BeEquivalentTo("9aGx9K5Vx0cM"))
				Expect(c.ClientID).To(BeEquivalentTo("p-config-server-c4a56a3d-9507-4c2f-9cd1-f858dbf9e11c"))
				Expect(c.AccessTokenURI).To(BeEquivalentTo("https://p-spring-cloud-services.uaa.cf.wise.com/oauth/token"))
			})
		})
		Context("No config server bound", func() {
			It("must return an error", func() {
				os.Setenv("VCAP_SERVICES", testVCAPwithoutConfig)
				c, err := GetConfigServerCredentialsFromEnv()
				Expect(err).NotTo(BeNil())
				Expect(c).To(BeNil())
			})
		})
	})

	Describe("Cloud Config Server Response", func() {
		It("must be parsable as simple key int pairs", func() {
			var conf ConfigServerResponse
			err := json.Unmarshal([]byte(testConfigServerResponse), &conf)
			Expect(err).To(BeNil())
			Expect(conf.Name).To(BeEquivalentTo("gonfig"))
			Expect(conf.Label).To(BeEquivalentTo("master"))
			Expect(conf.Version).To(BeEquivalentTo("77091415ec22c45e8f76407e2b30c226b33271"))
			Expect(len(conf.PropertySources)).To(BeEquivalentTo(1))
			Expect(conf.PropertySources[0].Name).To(BeEquivalentTo("https://github.com/dgruber/sample-config/gonfig.yml"))
			Expect(conf.PropertySources[0].Source["resolutionX"].(float64)).To(BeEquivalentTo(640))
			Expect(conf.PropertySources[0].Source["resolutionY"].(float64)).To(BeEquivalentTo(480))
		})
	})

	Describe("Fetch Config", func() {
		It("must return the configured result", func() {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, testConfigServerResponse)
			}))
			defer ts.Close()

			newConfig := strings.Replace(testVCAPwithConfig,
				"https://config-51711835-4626-4823-b5a1-e5d91012f3f2.apps.wise.com", ts.URL, 1)

			os.Setenv("VCAP_SERVICES", newConfig)
			os.Setenv("gonfig_testing", "1")

			config, err := FetchConfig()

			Expect(err).To(BeNil())
			Expect(config["resolutionX"].(float64)).To(BeEquivalentTo(640))
			Expect(config["resolutionY"].(float64)).To(BeEquivalentTo(480))
		})
	})

})
