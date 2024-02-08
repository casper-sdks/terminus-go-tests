## Terminus Go

This repo holds a set of tests to be run against the Casper Go SDK.

Points to note are:

- The tests can be run manually via the Terminus project [here](https://github.com/casper-sdks/terminus) 
- The tests are built using Cucumber features


### How to run locally

- Install Go as per the instructions [here](https://go.dev/doc/install)

- Clone repo and start NCTL (please note the NCTL Casper node version in the script 'docker-run')

  ```bash
  git clone git@github.com:casper-sdks/terminus-go-tests.git
  cd terminus-js-tests/script
  chmod +x docker-run && ./docker-run
  chmod +x docker-copy-assets && /docker-copy-assets 
  cd ..
  ```

- Go get the required SDK branch and run the tests

  ```bash
  go get github.com/make-software/casper-go-sdk@$[required-branch]
  go install gotest.tools/gotestsum@latest
  mkdir reports
  gotestsum --format standard-verbose --junitfile reports/report.xml
  ```

- TODO script the above

- JUnit test results will be output to /reports

### How to run locally IDE

Alternatively the tests can be run using an IDE

They are developed using JetBrains GoLand
