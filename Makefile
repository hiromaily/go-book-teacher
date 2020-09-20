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
.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: imports
imports:
	./scripts/imports.sh

###########################################################
# go list for check import package
###########################################################
.PHONY: golist
golist:
	go list -f '{{.ImportPath}} -> {{join .Imports "\n"}}' ./cmd/book/main.go


###############################################################################
# Build
###############################################################################
.PHONY: build
build:
	go build -i -v -o ${GOPATH}/bin/book ./cmd/book/

# run by save:text, notify:command using defined teacher data
.PHONY: exec1
exec1:
	book -t ./config/toml/text-command.toml

# run by save:text, notify:command using defined teacher jsondata
.PHONY: exec2
exec2:
	book -t ./config/toml/text-command.toml -j ./testdata/json/teachers.json

# run by save:text, notify:command using defined teacher data with loop
.PHONY: exec3
exec3:
	book -t ./config/toml/text-command.toml -i 10

# run by save:text, notify:slack using defined teacher data
.PHONY: exec4
exec4:
	book -t ./config/toml/text-slack.toml -crypto

# run by save:redis, notify:command using defined teacher data
.PHONY: exec5
exec5:
	book -t ./config/toml/redis-command.toml -crypto

.PHONY: exec-heroku
exec-heroku:
	book -t ./config/toml/heroku.toml -crypto

# run by save:text, notify:mail using defined teacher data
# for now, mail is not available because of security issue
# exec10:
# 	book -t ./config/toml/text-mail.toml


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
	go run ./tools/encryption/ -m e important-password

.PHONY: tool-decode
tool-decode:
	go run ./tools/encryption/ -m d o5PDC2aLqoYxhY9+mL0W/IdG+rTTH0FWPUT4u1XBzko=

###############################################################################
# Utility
###############################################################################

.PHONY: clean
clean:
	rm -rf status.log

###############################################################################
# Docker
###############################################################################
.PHONY: dclogin
dclogin:
	docker-compose exec book bash

.PHONY: dcexec
dcexec:
	docker-compose exec book /bin/sh -c "book -t ./config/toml/text-command.toml"

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
	heroku run book -t /app/config/toml/heroku.toml -crypto

.PHONY: heroku-info
heroku-info:
	#heroku config | grep REDIS
	heroku config
	heroku ps
