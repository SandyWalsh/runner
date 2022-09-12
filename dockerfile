# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /runner

WORKDIR gen
COPY gen/*.go ./

WORKDIR ../library
COPY library/*.go ./
COPY library/go.mod ./
COPY library/go.sum ./
RUN go mod download
RUN go build 

WORKDIR ../wrapper
COPY wrapper/*.go ./
COPY wrapper/go.mod ./
RUN go mod download
RUN go build -o wrapper

WORKDIR ../cert
COPY cert/* ./
