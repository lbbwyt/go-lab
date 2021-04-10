#!/bin/bash

echo "执行的文件名：$0";
protoc --proto_path=.  --proto_path=./..  --go_out=plugins=grpc:. $1