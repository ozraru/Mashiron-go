all: test mv
test:
	bash -n splatoon
	bash -n help.sh
	bash -n stage.sh
mv:
	mkdir -p ../../bin/cmd/splatoon
	cp splatoon ../../bin/cmd/splatoon/splatoon
	cp stage.sh ../../bin/cmd/splatoon/stage.sh
	cp help.sh ../../bin/cmd/splatoon/help.sh
	cp README ../../bin/cmd/splatoon/README
