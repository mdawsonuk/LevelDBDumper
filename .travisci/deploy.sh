cd src/LevelDBDumper
export GOOS="linux"
go build -o LevelDBDumper
export GOOS="darwin"
go build -o LevelDBDumper.app
export GOOS="windows"
go get golang.org/x/sys/windows
go build -o LevelDBDumper.exe
export GOARCH="386"
go build -o LevelDBDumper_x86.exe
zip -q LevelDBDumper.app.zip LevelDBDumper.app