# Note: tabs by space can't not used for Makefile!

CURRENTDIR=`pwd`

###############################################################################
# Managing Dependencies
###############################################################################
update:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u -d -v ./...


###############################################################################
# Golang formatter and detection
###############################################################################
lint:
	golangci-lint run --fix

imports:
	./scripts/imports.sh

###########################################################
# go list for check import package
###########################################################
golist:
	go list -f '{{.ImportPath}} -> {{join .Imports "\n"}}' ./cmd/book/main.go


###############################################################################
# Build Local
###############################################################################
build:
	go build -i -v -o ${GOPATH}/bin/book ./cmd/book/

# run by save:text, notify:command using defined teacher data
exec1:
	book -t ./config/toml/text-command.toml -crypto

# run by save:text, notify:command using defined teacher jsondata
exec2:
	book -t ./config/toml/text-command.toml -j ./testdata/json/teachers.json -crypto

# run by save:text, notify:command using defined teacher data with loop
exec3:
	book -t ./config/toml/text-command.toml  -i 10 -crypto

# run by save:text, notify:slack using defined teacher data
exec4:
	book -t ./config/toml/text-slack.toml -crypto

# run by save:redis, notify:command using defined teacher data
exec5:
	book -t ./config/toml/redis-command.toml -crypto

exec-heroku:
	book -t ./config/toml/heroku.toml -crypto

# run by save:text, notify:mail using defined teacher data
# for now, mail is not available because of security issue
# exec10:
# 	book -t ./config/toml/text-mail.toml


.PHONY: clean
clean:
	rm -rf status.log

###############################################################################
# Docker
###############################################################################
dc_create:
	./docker-create.sh

dclogin:
	docker-compose exec book bash

dcbld:
	docker-compose exec book bash ./docker-entrypoint.sh
	#docker-compose exec book /bin/sh -c "go build -i -race -v -o /go/bin/book ./cmd/book/"

dcexec:
	docker-compose exec book /bin/sh -c "book -t ./data/toml/settings.toml"

###############################################################################
# Test
###############################################################################
test:
	go test -v -race ./cmd/book
	#go test -v -covermode=count -coverprofile=profile.cov cmd/book/*.go -t ${PWD}/data/toml/mailon.toml

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
heroku:
	git push -f heroku master

herokuinfo:
	#heroku config | grep REDIS
	heroku config
	heroku ps
