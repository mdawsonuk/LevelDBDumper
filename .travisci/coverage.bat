@echo off
cd ..\src\LevelDBDumper
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out
del coverage.out
pause