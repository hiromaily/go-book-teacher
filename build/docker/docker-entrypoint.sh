#!/bin/sh
###
# initialize for docker environment
###

#go get -u -d -v ./...
go build -v -o /go/bin/book ./cmd/book/