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
	mkdir -p ../../bin/cmd/vote
	mv vote ../../bin/cmd/vote/vote
	cp README ../../bin/cmd/vote/README
