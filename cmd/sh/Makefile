all: get test build mv
test:
	go test
build:
	mkdir -p ../../bin/cmd/sh
	go build
	sudo ./init.sh
clean:
	go clean
get:
	go get gopkg.in/ini.v1
	go get go.etcd.io/bbolt
mv:
	mv sh ../../bin/cmd/sh/sh
	cp init.sh ../../bin/cmd/sh/init.sh
	cp mashironrc ../../bin/cmd/sh/mashironrc
	cp run.sh ../../bin/cmd/sh/run.sh
	cp README ../../bin/cmd/sh/README
	cp Makefile ../../bin/cmd/sh/Makefile

