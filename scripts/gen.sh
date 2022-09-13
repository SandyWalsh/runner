#!/bin/bash

protoc proto/service.proto --go_out=./sprinter --go_opt=paths=source_relative --go-grpc_out=./sprinter --go-grpc_opt=paths=source_relative \
