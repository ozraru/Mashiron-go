all: get test build mv
test:
	go test
build:
	go build platforms/discord/discord.go
	go build platforms/line/line.go
	bash cmd/mkall.sh
clean:
	go clean
#	sudo rm -rf bin
get:
	go get gopkg.in/ini.v1
	go get github.com/bwmarrin/discordgo
	go get go.etcd.io/bbolt
mv:
	mkdir -p bin
	mv discord bin/discord
	mv line bin/line
	cp -r skel bin/skel
	cp mashiron.ini bin/
	cp mashiron.service bin/
	cp mashiron.sh bin/
	cp ascii.txt bin/
release: clean all
	sudo tar -Jcvf Mashiron-go.tar.xz bin/*

