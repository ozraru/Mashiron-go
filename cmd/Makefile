all: get test build mv
test:
	go test
build:
	go build
clean:
	go clean
get:
	go get gopkg.in/ini.v1
mv:
	mkdir -p ../bin/cmd
	mv cmd ../bin/cmd/cmd
