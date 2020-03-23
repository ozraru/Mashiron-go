all: get test build mv
test:
	go test
build:
	go build discord.go
	bash cmd/mkall.sh
clean:
	go clean
	sudo rm -rf bin
get:
	go get gopkg.in/ini.v1
	go get github.com/bwmarrin/discordgo
mv:
	mkdir -p bin
	mv discord bin/discord
	cp -r skel bin/skel
	cp mashiron.ini bin/
	cp mashiron.service bin/
release: clean all
	sudo tar -Jcvf Mashiron-go.tar.xz bin/*
