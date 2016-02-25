#!/usr/bin/env bash
#go get -u -v ./

go fmt ./...
go vet ./...

#go build -o ./book ./book.go
go build -o ./book
