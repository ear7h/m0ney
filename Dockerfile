FROM golang:1.9-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git bash


WORKDIR /go/src/github.com/ear7h/m0ney
COPY . .

WORKDIR /go/src/github.com/ear7h/m0ney/daemon
RUN go get ./...
RUN go build .

WORKDIR /go/src/github.com/ear7h/m0ney/
RUN go get ./...
RUN go build .

CMD daemon/daemon & ./m0ney