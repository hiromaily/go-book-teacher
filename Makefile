# Note: tabs by space can't not used for Makefile!

CURRENTDIR=`pwd`

###############################################################################
# Managing Dependencies
###############################################################################
update:
	go get -u -d -v ./...


###############################################################################
# Golang formatter and detection
###############################################################################
fmt:
	go fmt `go list ./... | grep -v '/vendor/'`

vet:
	go vet `go list ./... | grep -v '/vendor/'`

fix:
	go fix `go list ./... | grep -v '/vendor/'`

lint:
	golint ./... | grep -v '^vendor\/' || true
	misspell `find . -name "*.go" | grep -v '/vendor/'`
	ineffassign .

chk:
	go fmt `go list ./... | grep -v '/vendor/'`
	go vet `go list ./... | grep -v '/vendor/'`
	go fix `go list ./... | grep -v '/vendor/'`
	golint ./... | grep -v '^vendor\/' || true
	misspell `find . -name "*.go" | grep -v '/vendor/'`
	ineffassign .


###########################################################
# go list for check import package
###########################################################
golist:
	go list -f '{{.ImportPath}} -> {{join .Imports "\n"}}' ./cmd/book/main.go


###############################################################################
# Build Local
###############################################################################
bld:
	go build -i -v -o ${GOPATH}/bin/book ./cmd/book/

run:
	rm -rf ./status.log
	go run ./cmd/book/main.go -t data/toml/local.toml
	#go run ./cmd/book/main.go -i 90
	#go run ./cmd/book/main.go -t data/toml/settings.toml -i 90


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
