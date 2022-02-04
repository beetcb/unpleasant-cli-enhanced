#!/bin/sh

GOOS=windows GOARCH=amd64 go build . &&
cp enclash.exe /mnt/share/bin &&
rm enclash.exe

