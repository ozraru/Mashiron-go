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
	mkdir -p ../../bin/cmd/help
	mv help ../../bin/cmd/help/help
	cp README ../../bin/cmd/help/README
	cp mashiron.png ../../bin/cmd/help/mashiron.png
