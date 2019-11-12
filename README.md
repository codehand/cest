# cest
[![Build Status](https://travis-ci.com/codehand/cest.svg?token=xSfYAJ5sB8Z6maxH16Mj&branch=beta)](https://travis-ci.com/codehand/cest)
[![codecov](https://codecov.io/gh/codehand/cest/branch/beta/graph/badge.svg?token=22X76FVtsG)](https://codecov.io/gh/codehand/cest)
## Contents

* [Contents](#contents)
* [How to build and install](#how-to-build-and-install)
* [How to use](#how-to-use)
* [License](#license)

## How to build and install

    go install github.com/codehand/cest
    
## Supported Versions

`cest` is tool genarate test with echo labstack 


## Changes
* Great support across OS's
* Easy cross compilation with GOOS and GOARCH
* Wrapper server echo v3 ([echo](https://echo.labstack.com/)) 
* Multi package, version
* Sub package handlers, skip handlers
* Restful API/CRUD
* Support model mongodb
* Eye catching user interface
* Progress bars
* Fix static echo version
* Handler exit process
* Echo swagger document

## Features
* Support gen service GRPC protocol
* Cache redis and memory internal
* Worker pool
* Genarate model, handler func
* Genarate config kube deployment
* Auto cron job
* Http client circuit breaker
* Plus redis shake
* Plus mongo shake
* Validation http request

## Examples

    cest -all .
    go get -u github.com/mjibson/esc