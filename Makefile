# Note: tabs by space can't not used for Makefile!

CURRENTDIR=`pwd`
modVer=$(shell cat go.mod | head -n 3 | tail -n 1 | awk '{print $2}' | cut -d'.' -f2)
currentVer=$(shell go version | awk '{print $3}' | sed -e "s/go//" | cut -d'.' -f2)

###############################################################################
# Managing Dependencies
###############################################################################
.PHONY: check-ver
check-ver:
	#echo $(modVer)
	#echo $(currentVer)
	@if [ ${currentVer} -lt ${modVer} ]; then\
		echo go version ${modVer}++ is required but your go version is ${currentVer};\
	fi

.PHONY: update
update:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u -d -v ./...


###############################################################################
# Golang formatter and detection
###############################################################################
.PHONY: imports
imports:
	./scripts/imports.sh

.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: lintall
lintall: imports lint


###############################################################################
# Build
###############################################################################
.PHONY: build
build:
	go build -i -v -o ${GOPATH}/bin/book ./cmd/book/

.PHONY: run
run: build
	book


###############################################################################
# Test
###############################################################################
.PHONY: test
test:
	go test -v -race ./cmd/book
	#go test -v -covermode=count -coverprofile=profile.cov cmd/book/*.go -t ${PWD}/data/toml/mailon.toml


###############################################################################
# Tools
# Note: environment variable `ENC_KEY`, `ENC_IV` should be set in advance
###############################################################################
.PHONY: tool-encode
tool-encode:
	go run ./tools/encryption/ -encode important-password

.PHONY: tool-decode
tool-decode:
	go run ./tools/encryption/ -decode o5PDC2aLqoYxhY9+mL0W/IdG+rTTH0FWPUT4u1XBzko=


###############################################################################
# Docker
###############################################################################
.PHONY: dclogin
dclogin:
	docker-compose exec book bash

.PHONY: dcexec
dcexec:
	docker-compose exec book /bin/sh -c "book -t ./config/text-command.toml"

###############################################################################
# Build Heroku
#
#heroku config:add HEROKU_FLG=1
#heroku addons:create scheduler:standard

#heroku run book -t /app/config/toml/heroku.toml -crypto
#heroku run bash
#heroku logs -t
#heroku ps -a book
#heroku ps
#heroku config
#
# $ git push -f heroku master
#  The go.mod file for this project does not specify a Go version
#  https://devcenter.heroku.com/articles/go-apps-with-modules#build-configuration
#
###############################################################################
.PHONY: heroku-deploy
heroku-deploy:
	git push -f heroku master

.PHONY: heroku-run
heroku-run:
	heroku run book -t /app/config/heroku.toml

.PHONY: heroku-info
heroku-info:
	#heroku config | grep REDIS
	heroku config
	heroku ps
