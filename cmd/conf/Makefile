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
	mkdir -p ../../bin/cmd/conf
	mv conf ../../bin/cmd/conf/conf
	cp README ../../bin/cmd/conf/README
