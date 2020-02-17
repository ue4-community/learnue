@echo off

setlocal

if exist install.bat goto ok
echo install.bat must be run from its folder
goto end

:ok

set GOPROXY=https://goproxy.cn
set GO111MODULE=on

if not exist log mkdir log

gofmt -w -s .

go build -o bin/studygolang.exe github.com/ue4-community/learnue/cmd/studygolang

:end
echo finished