#!/usr/bin/env bash
cd contracts;
docker run --rm -u $(id -u):$(id -g) -v $PWD:/contracts -w /contracts thethingsindustries/protoc --go_out=plugins=grpc:. -I. ./*.proto
docker run --rm -u $(id -u):$(id -g) -v $PWD:/contracts -w /contracts thethingsindustries/protoc --grpc-gateway_out=logtostderr=true:. -I. ./*.proto
docker run --rm -u $(id -u):$(id -g) -v $PWD:/contracts -w /contracts thethingsindustries/protoc --swagger_out=logtostderr=true:. -I. ./*.proto
