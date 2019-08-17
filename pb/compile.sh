#!/usr/bin/env sh

protoc orders.proto --go_out=plugins=grpc:.