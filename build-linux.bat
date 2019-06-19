@echo off

if exist build-linux.bat goto ok
echo install.bat must be run from its folder
goto end

: ok
SET dist=go

set GOPATH=%~dp0;%GOPATH%

SET GOOS=linux
SET GOARCH=amd64

go build -o bin/%dist%_%GOOS%_%GOARCH% main

:end
echo finished