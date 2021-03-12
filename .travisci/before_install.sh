chmod +x ./.travisci/test.sh
chmod +x ./.travisci/deploy.sh
export GO111MODULE=auto
go get github.com/syndtr/goleveldb/leveldb
go get github.com/hashicorp/go-version
go get github.com/gookit/color
if [ "$TRAVIS_OS_NAME" = "windows" ]; then go get golang.org/x/sys/windows; fi