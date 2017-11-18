FROM golang:1.9-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git bash


WORKDIR /go/src/daemon
COPY . .

RUN go get -v ./...
RUN go build
CMD ["./daemon"]