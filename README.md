# gonfig
Accessing PCF Config Server from Go Apps

Simple library which can be used within Go (#golang) microservice applications in order to get the
configuration from the PCF Config Server.

All what it does is fetching the access details of the bound PCF Configuration Server from the
VCAP enviornment variables and requests the configuration from the Config Server by using the given
credentials.
