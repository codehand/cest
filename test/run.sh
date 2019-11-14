rm -rf ../cmd/tests
go install github.com/codehand/cest && cest -all handler.go