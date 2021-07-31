FROM golang:latest

ENV GO111MODULE=auto

EXPOSE 5000

RUN mkdir -p /go/src/github.com/go-challenge

WORKDIR /go/src/github.com/go-challenge

COPY . /go/src/github.com/go-challenge

RUN chmod a+x /go/src/github.com/go-challenge


RUN apt-get update
RUN apt-get upgrade -y

ENV GOBIN /go/bin

RUN go get github.com/go-redis/redis

RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get go.mongodb.org/mongo-driver/bson


CMD ["go","run","/go/src/github.com/go-challenge/main.go"]