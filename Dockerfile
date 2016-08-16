# Dcokerfile for go-book-teacher

FROM golang:1.6

ARG redisHostName=default-redis-server

RUN mkdir -p /go/src/github.com/hiromaily/go-book-teacher
#COPY ./go-book-teacher /go/src/github.com/hiromaily/go-book-teacher/

ENV REDIS_URL=redis://h:password@${redisHostName}:6379

WORKDIR /go/src/github.com/hiromaily/go-book-teacher

#RUN go get -d -v ./... && go build -v -o /go/bin/book ./cmd/book/
