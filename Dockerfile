FROM golang:1.8.3-alpine3.5

RUN apk update && apk upgrade && \
    apk add --no-cache git bash &&\
    apk add --no-cache mysql-client

WORKDIR /go/src/m0ney
COPY . .


RUN go get ./...
RUN go build
RUN cd daemon && \
    go build && \
    cd ..

CMD ["./init.sh"]