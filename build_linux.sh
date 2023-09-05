#!/bin/bash
#CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="x86_64-linux-musl-cc"  go build -ldflags="-L/usr/local/include -lta_lib -Wl" main.go
ServerIp="112.213.97.135"
ServerPwd="Yuanshi20188"
serverPath="/www/binance"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" go build main.go
sshpass -p $ServerPwd scp -r -vvv ./main root@$ServerIp:$serverPath/main_wsj