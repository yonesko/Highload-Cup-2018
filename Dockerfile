FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/yonesko/Highload-Cup-2018/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN go build -tags=jsoniter -o /go/bin/hello
EXPOSE 80
CMD $GOPATH/bin/hello