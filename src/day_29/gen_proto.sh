#!/bin/bash

# 扫描proto文件夹，把文件夹里所有的.proto文件都取到，然后执行 protoc --go_out=. --go-grpc_out=. [文件名]

for file in proto/*.proto; do
    protoc --go_out=.  --go-grpc_out=. "$file"
done
