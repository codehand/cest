# Go super project Makefile
# Authored by mac

# You want to edit these
REPOSITORY	:= gitlab.com/travelr/crm
PACKAGES 	:= common common/sub

# Maybe even these
CMD_DIR		:= cmd
GOFMT		:= @gofmt -s -tabwidth=4 -w -l
COV_EXT     := cov
TEST_COV	:= @go test -cover -coverprofile=
CLI_COV		:= go tool cover -func=
WEB_COV 	:= go tool cover -html=
GOTEST		:= @go test
GORUN		:= go run
GOVET		:= @go vet
OPEN_CMD	:= xdg-open
ZIP_FILE	:= godoc.zip
BIN_DIR		:= bin

# But dont edit these
GOOS		:= $(shell go env GOOS)
GOARCH		:= $(shell go env GOARCH)
GOBUILD		:= GOOS=$(GOOS) GOARCH=$(GOARCH) go build
GOINSTALL	:= GOOS=$(GOOS) GOARCH=$(GOARCH) go install -v
RUNNABLES	:= $(wildcard $(CMD_DIR)/*.go)
PKG_ROOT	:= $(GOPATH)/pkg/$(GOOS)_$(GOARCH)/$(REPOSITORY)
A_FILES		:= $(foreach pkg, $(PACKAGES), $(PKG_ROOT)/$(pkg).a)
TESTABLE	:= $(foreach pkg, $(PACKAGES), $(wildcard $(pkg)/*.go)) 
GOFILES		:= $(TESTABLE) $(RUNNABLE)
COV_FILES	:= $(foreach pkg, $(PACKAGES), $(pkg).$(COV_EXT))

default:	test1 test2 test3 test4
test1:
		$(RM) -f **/*_test.go
		cest -all .
		go test ./... -v
test2:
		$(RM) -f **/*_test.go
		$(RM) -rf tests
		cest -all -output=tests ...
		go test ./... -v
test3:
		$(RM) -f **/*_test.go
		cest -all *
		go test ./... -v
test4:
		$(RM) -f **/*_test.go
		cest -all ...
		go test ./... -v
test5:
		$(RM) -f **/*_test.go
		$(RM) -rf tests
		- cest -only=EchoContext -output=tests handler/abc/abc.go
		go test ./... -v