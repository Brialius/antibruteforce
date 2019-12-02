# antibruteforce
[![Go Report Card](https://goreportcard.com/badge/github.com/Brialius/antibruteforce)](https://goreportcard.com/report/github.com/Brialius/antibruteforce)

## Build
### make goals
|Goal|Description|
|----|-----------|
|build (default)|build binaries|
|build-server|build server binary|
|build-client|build client binary|
|deploy|run docker-compose environment|
|undeploy|destroy docker-compose environment|
|deploy-tests|run test docker-compose environment and execute integration tests|
|undeploy-tests|destroy test docker-compose environment|
|build (default)|build binaries|
|setup|download and install required dependencies|
|test|run tests|
|integration-test|run integration tests|
|install|install binary to `$GOPATH/bin`|
|lint|run linters|
|clean|run `go clean`|
|mod-refresh|run `go mod tidy` and `go mod vendor`|
|ci|run all steps needed for CI|
|version|show current git tag if any matched to `v*` exists|
|release|set git tag and push to repo `make release ver=v1.2.3`|
