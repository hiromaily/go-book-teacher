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

# lint:
# 	golint ./... | grep -v '^vendor\/' || true
# 	misspell `find . -name "*.go" | grep -v '/vendor/'`
# 	ineffassign .
#
# fmt:
# 	go fmt `go list ./... | grep -v '/vendor/'`
#
# vet:
# 	go vet `go list ./... | grep -v '/vendor/'`
#
# fix:
# 	go fix `go list ./... | grep -v '/vendor/'`
#
# chk:
# 	go fmt `go list ./... | grep -v '/vendor/'`
# 	go vet `go list ./... | grep -v '/vendor/'`
# 	go fix `go list ./... | grep -v '/vendor/'`
# 	golint ./... | grep -v '^vendor\/' || true
# 	misspell `find . -name "*.go" | grep -v '/vendor/'`
# 	ineffassign .


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
	book -t ./data/toml/text-command.toml

# run by save:text, notify:command using defined teacher jsondata
exec2:
	book -t ./data/toml/text-command.toml -j ./data/json/teachers.json

# run by save:text, notify:command using defined teacher data with loop
exec3:
	book -t ./data/toml/text-command.toml  -i 10

# run by save:text, notify:slack using defined teacher data
exec4:
	book -t ./data/toml/text-slack.toml

# run by save:redis, notify:command using defined teacher data
exec5:
	book -t ./data/toml/redis-command.toml

# run by save:text, notify:mail using defined teacher data
# for now, mail is not available because of security issue
exec6:
	book -t ./data/toml/text-mail.toml

run:
	rm -rf ./status.log
	go run ./cmd/book/main.go -t data/toml/local.toml
	#go run ./cmd/book/main.go -i 90
	#go run ./cmd/book/main.go -t data/toml/settings.toml -i 90


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
test1:
	# mail mode
	go test -v -covermode=count -coverprofile=profile.cov cmd/book/*.go -t ${PWD}/data/toml/mailon.toml

test2:
	# slack mode
	go test -v -covermode=count -coverprofile=profile.cov cmd/book/*.go -t ${PWD}/data/toml/slackon.toml

test3:
	go test -v -covermode=count -coverprofile=profile.cov cmd/book/*.go -run TestIntegrationOnLocalUsingTxtAndBrowserAndJson

test4:
	go test -covermode=count -coverprofile=profile.cov -v cmd/book/*.go -run TestIntegrationOnLocalUsingTxtAndBrowser

test5:
	go test -covermode=count -coverprofile=profile.cov -v cmd/book/*.go -run TestIntegrationOnLocalUsingRedisAndMail

test6:
	go test -v -covermode=count -coverprofile=profile.cov cmd/book/*.go -run TestIntegrationOnLocalUsingTxtAndBrowser


###############################################################################
# Build Heroku
#
#heroku config:add HEROKU_FLG=1
#heroku addons:create scheduler:standard

#heroku run book -t /app/data/toml/settings.toml
#heroku run bash
#heroku logs -t
#heroku ps -a book
#heroku ps
#heroku config
#
###############################################################################
heroku:
	git push -f heroku master

herokuinfo:
	heroku config
	heroku ps
