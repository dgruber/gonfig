# gonfig
Accessing the Pivotal Cloud Foundry Config Server from Go apps

Simple library which can be used within Go (#golang) microservice applications in order to get the
configuration from the PCF Config Server.

All what it does is fetching the access details of the bound PCF Configuration Server from the
VCAP environment variables and requests the configuration from the Config Server by using the given
credentials.
