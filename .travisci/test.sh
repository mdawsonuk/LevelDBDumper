cd src/LevelDBDumper
go test -v -cover -count=1 -coverprofile=coverage.out
go tool cover -func=coverage.out