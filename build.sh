#!/bin/sh

###########################################################
# Variable
###########################################################
#export GOTRACEBACK=single
GOTRACEBACK=all
CURRENTDIR=`pwd`

###########################################################
# Update all package
###########################################################
#go get -u -v ./...


###########################################################
# Adjust version dependency of projects
###########################################################
#cd ${GOPATH}/src/github.com/aws/aws-sdk-go
#git checkout v0.9.17
#git checkout master


###########################################################
# go fmt and go vet
###########################################################
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
#go list -f '{{.ImportPath}} -> {{join .Imports "\n"}}' ./ginserver.go


###########################################################
# go build and install
###########################################################
#-n show just command for build
#go build -i -n -o ./ginserver ./ginserver.go

#rebuild dependent packages (rebuild all package)
#go build -a -v -o ./ginserver ./ginserver.go

#build and install
#go build -i -v -o book ./cmd/book/
#go build -i -v ./cmd/book/
go build -i -v -o ${GOPATH}/bin/book ./cmd/book/
EXIT_STATUS=$?

if [ $EXIT_STATUS -gt 0 ]; then
    exit $EXIT_STATUS
fi

exit 0

###########################################################
# go test
###########################################################
#MAIL_TO_ADDRESS=xxxxx@gmail.com
#MAIL_FROM_ADDRESS=xxxx@xxxx.com
#SMTP_ADDRESS=xxxx@xxxx.com
#SMTP_PASS=xxxxx
#SMTP_SERVER=smtp.xxxx.com
#SMTP_PORT=587

#1.All
#go test -v cmd/book/*.go -toadd ${MAIL_TO_ADDRESS} -fradd ${MAIL_FROM_ADDRESS} -smpass ${SMTP_PASS} -smsvr ${SMTP_SERVER} -smport ${SMTP_PORT}
#2.json
#go test -v cmd/book/*.go -run TestIntegrationOnLocalUsingTxtAndBrowserAndJson -toadd ${MAIL_TO_ADDRESS} -fradd ${MAIL_FROM_ADDRESS} -smpass ${SMTP_PASS} -smsvr ${SMTP_SERVER} -smport ${SMTP_PORT}
#3.saved file test
#go test -v cmd/book/*.go -run TestIntegrationOnLocalUsingTxtAndBrowser -toadd ${MAIL_TO_ADDRESS} -fradd ${MAIL_FROM_ADDRESS} -smpass ${SMTP_PASS} -smsvr ${SMTP_SERVER} -smport ${SMTP_PORT}
#go test -v cmd/book/*.go -run TestIntegrationOnLocalUsingRedisAndMail -toadd ${MAIL_TO_ADDRESS} -fradd ${MAIL_FROM_ADDRESS} -smpass ${SMTP_PASS} -smsvr ${SMTP_SERVER} -smport ${SMTP_PORT}


###########################################################
# exec
###########################################################
#./book -f ${GOPATH}/src/github.com/hiromaily/booking-teacher/settings.json


###########################################################
# godep
###########################################################
#go get -u github.com/tools/godep

#Save
rm -rf Godeps
rm -rf vendor

godep save ./...
EXIT_STATUS=$?

if [ $EXIT_STATUS -gt 0 ]; then
    exit $EXIT_STATUS
fi

#Build
#godep go build -o book ./cmd/book/

#Restore
#godep restore

exit 0

###########################################################
# git add, commit, push
###########################################################
git recm
git pufom


###########################################################
# heroku
###########################################################
#heroku config:add HEROKU_FLG=1
#heroku addons:create scheduler:standard

git push -f heroku master

#heroku run book
#heroku run bash
#heroku logs -t
#heroku ps -a book
#heroku ps
#heroku config


###########################################################
# cross-compile for linux
###########################################################
#GOOS=linux go install -v ./...


###########################################################
# godoc
###########################################################
#godoc -http :8000
#http://localhost:8000/pkg/
