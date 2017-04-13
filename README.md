# gonfig

Accessing the Pivotal Cloud Foundry Config Server from Go apps

Simple library which can be used within Go (#golang) microservice applications in order to get the
configuration from the PCF Config Server.

All what it does is fetching the access details of the bound PCF Configuration Server from the
VCAP environment variables and requests the configuration from the Config Server by using the given
credentials.

## CI Status

Pivotal's Concourse CI/CF tool also provides badges! For an example pipeline compiling and testing the project, see below:

| Job | Status |
|---------|--------|
| Building | [![Build Badge](http://ci.route.today/api/v1/teams/main/pipelines/gonfig/jobs/building/badge)](http://ci.route.today/teams/main/pipelines/gonfig/jobs/building) |
| Testing | [![Test Badge](http://ci.route.today/api/v1/teams/main/pipelines/gonfig/jobs/testing/badge)](http://ci.route.today/teams/main/pipelines/gonfig/jobs/testing) |
| Push of Example | [![Push Badge](http://ci.route.today/api/v1/teams/main/pipelines/gonfig/jobs/example/badge)](http://ci.route.today/teams/main/pipelines/gonfig/jobs/example) |

See the [complete pipeline](http://ci.route.today/teams/main/pipelines/gonfig) for more details
