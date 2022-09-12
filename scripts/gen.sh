#!/bin/bash

protoc proto/service.proto --go-grpc_out=/home/sandy/go/src --go_out=/home/sandy/go/src
