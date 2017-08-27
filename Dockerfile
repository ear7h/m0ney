FROM golang:1.8.3-alpine3.5

RUN apk update && apk upgrade && \
    apk add --no-cache git bash

WORKDIR /go/src/m0ney
COPY . .

RUN go get ./...
RUN go build

CMD ["./init.sh"]