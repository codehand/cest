# cest
[![Build Status](https://travis-ci.com/codehand/cest.svg?token=xSfYAJ5sB8Z6maxH16Mj&branch=master)](https://travis-ci.com/codehand/cest)
[![codecov](https://codecov.io/gh/codehand/cest/branch/master/graph/badge.svg)](https://codecov.io/gh/codehand/cest)
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
* 

## Features
* 

## Examples

    cest -all .
    cest -all ...
    cest -all *
    cest -all handler/abc/abc.go
    cest -all handler/abc
    cest -all -output=tests ...
    cest -only=EchoContext -output=tests handler/abc/abc.go
    