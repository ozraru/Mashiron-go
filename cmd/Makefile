all: get test build mv
test:
	go test
build:
	go build
clean:
	go clean
get:
	go get gopkg.in/ini.v1
	go get github.com/otiai10/copy
mv:
	mkdir -p ../bin/cmd
	mv cmd ../bin/cmd/cmd
