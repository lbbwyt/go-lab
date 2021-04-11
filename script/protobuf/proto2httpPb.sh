#!/bin/bash

echo "执行的文件名：$0";
protoc -I ../googleapis \
--proto_path=.  --proto_path=./..  --go_out=plugins=grpc:. $1

protoc -I ../googleapis \
--proto_path=.  --proto_path=./..  --grpc-gateway_out=logtostderr=true:. $1



#go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
#go get -u github.com/micro/protobuf/{proto,protoc-gen-go}