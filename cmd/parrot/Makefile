all: get test build party mv
test:
	go test
build:
	go build
clean:
	go clean
get:
	go get gopkg.in/ini.v1
party:
	git clone https://github.com/jmhobbs/cultofthepartyparrot.com.git
	mkdir parrots
	rm -rf ../../bin/cmd/parrot
	mv cultofthepartyparrot.com/parrots/hd/*.gif parrots
	mv cultofthepartyparrot.com/parrots/*.gif parrots
	mv cultofthepartyparrot.com/flags/hd/*.gif parrots
	mv cultofthepartyparrot.com/guests/hd/*.gif parrots
	mv cultofthepartyparrot.com/guests/*.gif parrots
	rm cultofthepartyparrot.com -rf
mv:
	mkdir -p ../../bin/cmd/parrot
	mv parrot ../../bin/cmd/parrot/parrot
	cp README ../../bin/cmd/parrot/README
	mv parrots ../../bin/cmd/parrot/parrots
