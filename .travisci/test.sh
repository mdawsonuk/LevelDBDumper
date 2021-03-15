cd src/LevelDBDumper
gosec -tests ./...
go test -v -cover -count=1 -coverprofile=coverage.out
go tool cover -func=coverage.out