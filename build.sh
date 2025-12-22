#!/bin/bash
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o cliface-darwin-amd64 .
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H windowsgui=0" -o cliface-windows-amd64.exe .
