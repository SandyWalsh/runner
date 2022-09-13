# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /runner


WORKDIR ../sprinter
COPY sprinter/*.go ./
COPY sprinter/go.mod ./
COPY sprinter/go.sum ./
COPY sprinter/proto/*.go ./proto/
RUN go mod download
RUN go build 

WORKDIR ../wrapper
COPY wrapper/*.go ./
COPY wrapper/go.mod ./
RUN go mod download
RUN go build -o wrapper

WORKDIR ../cert
COPY cert/* ./
