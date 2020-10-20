all: get test build check mv
test:
	go test
build:
	go build
clean:
	go clean
get:
	go get gopkg.in/ini.v1
check:
	./party.sh
party:
	git clone https://github.com/jmhobbs/cultofthepartyparrot.com.git
	mkdir parrots
	mv cultofthepartyparrot.com/parrots/hd/*.gif parrots
	mv cultofthepartyparrot.com/parrots/*.gif parrots
	mv cultofthepartyparrot.com/flags/hd/*.gif parrots
	mv cultofthepartyparrot.com/guests/hd/*.gif parrots
	mv cultofthepartyparrot.com/guests/*.gif parrots
	rm cultofthepartyparrot.com -rf
	mkdir -p ../../bin/cmd/parrot/parrots
	mv parrots ../../bin/cmd/parrot

mv:
	mkdir -p ../../bin/cmd/parrot
	mv parrot ../../bin/cmd/parrot/parrot
	cp README ../../bin/cmd/parrot/README
