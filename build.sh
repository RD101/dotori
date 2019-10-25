#!/bin/sh
APP="dotori"
GOOS=linux GOARCH=amd64 go build -o ./bin/linux/${APP} check.go dbapi.go dotori.go struct.go
GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/${APP} check.go dbapi.go dotori.go struct.go
GOOS=windows GOARCH=amd64 go build -o ./bin/windows/${APP} check.go dbapi.go dotori.go struct.go

# Github Release에 업로드 하기위해 압축
cd ./bin/linux/ && mkdir thumbnail && tar -zcvf ../${APP}_linux_x86-64.tgz . && cd -
cd ./bin/darwin/ && mkdir thumbnail && tar -zcvf ../${APP}_darwin_x86-64.tgz . && cd -
cd ./bin/windows/ && mkdir thumbnail && tar -zcvf ../${APP}_windows_x86-64.tgz . && cd -

# 삭제
rm -rf ./bin/linux
rm -rf ./bin/darwin
rm -rf ./bin/windows
