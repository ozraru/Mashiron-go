all: get test build mv
test:
	go test
build:
	go build
clean:
	go clean
get:
	go get gopkg.in/ini.v1
	go get go.etcd.io/bbolt
mv:
	mkdir -p ../../bin/cmd/man
	mv man ../../bin/cmd/man/man
	cp README ../../bin/cmd/man/README
