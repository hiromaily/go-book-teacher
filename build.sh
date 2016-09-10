#!/bin/sh

###########################################################
# Variable
###########################################################
#export GOTRACEBACK=single
GOTRACEBACK=all
CURRENTDIR=`pwd`

TEST_MODE=0  #0:off, 1:All, 2...5, 9:All and coverage.
AUTO_EXEC=0
GODEP_MODE=1
AUTO_GITCOMMIT=0
HEROKU_MODE=0
DOCKER_MODE=0

GO_GET=0
GO_LINT=0

# when using go 1.7 for the first time, delete all inside pkg directory and run go install.
#go install -v ./...
#go get -u -v ./...

###########################################################
# Update all package
###########################################################
if [ $GO_GET -eq 1 ]; then
    go get -u -v ./...
    #go get -d -v ./...
    go get -u github.com/tools/godep
fi

###########################################################
# go fmt and go vet
###########################################################
echo '============== go fmt; go vet; =============='
#go fmt ./...
go fmt `go list ./... | grep -v '/vendor/'`
#go vet ./...
go vet `go list ./... | grep -v '/vendor/'`
EXIT_STATUS=$?

if [ $EXIT_STATUS -gt 0 ]; then
    exit $EXIT_STATUS
fi

###########################################################
# go lint
###########################################################
# it's too strict
#go get -u github.com/golang/lint/golint
if [ $GO_LINT -eq 1 ]; then
    #golint ./...
    #golint `go list ./... | grep -v '/vendor/'`
    golint ./... | grep -v '^vendor\/' || true
fi


###########################################################
# go list for check import package
###########################################################
#go list -f '{{.ImportPath}} -> {{join .Imports "\n"}}' ./cmd/book/main.go


###########################################################
# Adjust version dependency of projects
###########################################################
#cd ${GOPATH}/src/github.com/aws/aws-sdk-go
#git checkout v0.9.17
#git checkout master


###########################################################
# go build and install
###########################################################
echo '============== go build -i -v -o; =============='
if [ $GODEP_MODE -eq 1 ]; then
    rm -rf Godeps
    rm -rf ./vendor
fi

#-n show just command for build
#go build -i -n ./cmd/book/

#rebuild dependent packages (rebuild all package)
#go build -a -v -o ${GOPATH}/bin/book ./cmd/book/

#build and install
go build -i -v -o ${GOPATH}/bin/book ./cmd/book/
EXIT_STATUS=$?

if [ $EXIT_STATUS -gt 0 ]; then
    exit $EXIT_STATUS
fi

###########################################################
# go test
###########################################################
if [ $TEST_MODE -ne 0 ]; then
    echo '============== test =============='

    #call another shell
    #sh ./mail.sh

    if [ $TEST_MODE -eq 1 ]; then
        echo '1.All'
        #1.All
        go test -v -covermode=count -coverprofile=profile.cov cmd/book/*.go \
        -t ${PWD}/config/settings.toml
    fi

    if [ $TEST_MODE -eq 2 ]; then
        echo '2.json'
        #2.json
        go test -v -covermode=count -coverprofile=profile.cov cmd/book/*.go \
        -run TestIntegrationOnLocalUsingTxtAndBrowserAndJson
    fi

    if [ $TEST_MODE -eq 3 ]; then
        echo '3.saved file test'
        #3.saved file test
        go test -covermode=count -coverprofile=profile.cov -v cmd/book/*.go \
        -run TestIntegrationOnLocalUsingTxtAndBrowser
    fi

    if [ $TEST_MODE -eq 4 ]; then
        echo '4.saved file test ver.2'
        go test -covermode=count -coverprofile=profile.cov -v cmd/book/*.go \
        -run TestIntegrationOnLocalUsingRedisAndMail
    fi

    if [ $TEST_MODE -eq 5 ]; then
        echo '5.Godep test check'
        #Godep test check
        godep go test -v -covermode=count -coverprofile=profile.cov cmd/book/*.go \
        -run TestIntegrationOnLocalUsingTxtAndBrowser
    fi
fi


###########################################################
# exec
###########################################################
if [ $AUTO_EXEC -eq 1 ]; then
    echo '============== exec =============='
    #when using json file
    #book -f ${PWD}/json/teachers/settings.json

    go run ./cmd/book/main.go
    #go run ./cmd/book/main.go -i 90
    #GOTRACEBACK=all go run ./cmd/book/main.go
    #GOTRACEBACK=all go run ./cmd/book/main.go -t settings.toml -i 90
fi

###########################################################
# godep
###########################################################
if [ $GODEP_MODE -eq 1 ]; then
    echo '============== godeps =============='

    #go get -u github.com/tools/godep

    #Save
    rm -rf Godeps
    rm -rf vendor

    godep save ./...
    EXIT_STATUS=$?

    if [ $EXIT_STATUS -gt 0 ]; then
        exit $EXIT_STATUS
    fi
fi

#Build
#godep go build -o book ./cmd/book/

#Restore
#godep restore


###########################################################
# git add, commit, push
###########################################################
if [ $AUTO_GITCOMMIT -eq 1 ]; then
    echo '============== git recm, pufom =============='
    git recm
    git pufom
    git st
fi


###########################################################
# heroku
###########################################################
if [ $HEROKU_MODE -eq 1 ]; then
    echo '============== heroku: git push =============='
    git push -f heroku master
fi

#heroku config:add HEROKU_FLG=1
#heroku addons:create scheduler:standard

#heroku run book -t /app/config/settings.toml
#heroku run bash
#heroku logs -t
#heroku ps -a book
#heroku ps
#heroku config


###########################################################
# Docker
###########################################################
if [ $DOCKER_MODE -eq 1 ]; then
    echo '============== docker =============='
    ./docker-create.sh
fi


###########################################################
# cross-compile for linux
###########################################################
#GOOS=linux go install -v ./...


###########################################################
# godoc
###########################################################
#godoc -http :8000
#http://localhost:8000/pkg/
