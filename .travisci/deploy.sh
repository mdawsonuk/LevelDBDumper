cd src/LevelDBDumper
export GOOS="linux"
go build -o LevelDBDumper
export GOOS="darwin"
go build -o LevelDBDumper.app
zip -q LevelDBDumper.app.zip LevelDBDumper.app
export GOOS="windows"
go get golang.org/x/sys/windows
go build -o LevelDBDumper.exe

export GOARCH="386"

export GOOS="windows"
go build -o LevelDBDumper_x86.exe
export GOOS="linux"
go build -o LevelDBDumper_x86