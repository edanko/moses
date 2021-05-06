@echo off
go build -trimpath -ldflags "-s -w" .
@pause