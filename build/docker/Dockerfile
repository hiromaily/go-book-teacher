# Dcokerfile for go-book-teacher
FROM golang:1.15.6-alpine3.12

ARG redisHostName=default-redis-server

RUN mkdir -p /go/src/github.com/hiromaily/go-book-teacher
WORKDIR /go/src/github.com/hiromaily/go-book-teacher
COPY . .

ENV REDIS_URL=redis://h:password@${redisHostName}:6379

#RUN go get -u -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/book ./cmd/book/main.go

CMD ["bash"]
