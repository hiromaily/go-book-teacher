#!/bin/sh

###########################################################
# Variable
###########################################################
#export GOTRACEBACK=single
GOTRACEBACK=all
CURRENTDIR=`pwd`

export TEST_MODE=0
AUTO_EXEC=0
GODEP_MODE=1
AUTO_GITCOMMIT=0
HEROKU_MODE=0
DOCKER_MODE=0

###########################################################
# Update all package
###########################################################
#go get -u -v ./...
#go get -d -v ./...
#go get -u github.com/tools/godep


###########################################################
# go fmt and go vet
###########################################################
echo '============== go fmt; go vet; =============='
go fmt ./...
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
#golint ./...


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
rm -rf ./vendor

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
    sh ./mail.sh

    #MAIL_TO_ADDRESS=xxxxx@gmail.com
    #MAIL_FROM_ADDRESS=xxxx@xxxx.com
    #SMTP_ADDRESS=xxxx@xxxx.com
    #SMTP_PASS=xxxxx
    #SMTP_SERVER=smtp.xxxx.com
    #SMTP_PORT=587

    #break
    TEST_MODE=0

    #Don't Run below. it moved to mail.sh.
    if [ $TEST_MODE -eq 1 ]; then
        #1.All
        go test -v cmd/book/*.go \
        -toadd ${MAIL_TO_ADDRESS} -fradd ${MAIL_FROM_ADDRESS} \
        -smpass ${SMTP_PASS} -smsvr ${SMTP_SERVER} -smport ${SMTP_PORT}
    fi

    if [ $TEST_MODE -eq 2 ]; then
        #2.json
        go test -v cmd/book/*.go -run TestIntegrationOnLocalUsingTxtAndBrowserAndJson \
        -toadd ${MAIL_TO_ADDRESS} -fradd ${MAIL_FROM_ADDRESS} \
        -smpass ${SMTP_PASS} -smsvr ${SMTP_SERVER} -smport ${SMTP_PORT}
    fi

    if [ $TEST_MODE -eq 3 ]; then
        #3.saved file test
        go test -v cmd/book/*.go -run TestIntegrationOnLocalUsingTxtAndBrowser \
        -toadd ${MAIL_TO_ADDRESS} -fradd ${MAIL_FROM_ADDRESS} \
        -smpass ${SMTP_PASS} -smsvr ${SMTP_SERVER} -smport ${SMTP_PORT}
    fi

    if [ $TEST_MODE -eq 4 ]; then
        go test -v cmd/book/*.go -run TestIntegrationOnLocalUsingRedisAndMail \
        -toadd ${MAIL_TO_ADDRESS} -fradd ${MAIL_FROM_ADDRESS} \
        -smpass ${SMTP_PASS} -smsvr ${SMTP_SERVER} -smport ${SMTP_PORT}
    fi
fi


###########################################################
# exec
###########################################################
if [ $AUTO_EXEC -eq 1 ]; then
    echo '============== exec =============='
    #when using json file
    #book -f ${PWD}/settings.json

    GOTRACEBACK=all go run ./cmd/book/main.go
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

#heroku run book
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
