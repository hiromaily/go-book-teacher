#!/bin/sh
###
# initialize for docker environment
###

go get -d -v ./...
go build -v -o /go/bin/book ./cmd/book/