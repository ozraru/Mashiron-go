package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Request struct {
	version string
	API     string
	ROOM    string
	USER    string
	PRIV    []string
	CONTENT string
}
type Dir struct {
	roomdir    string
	cmddatadir string
}
type Conf struct {
	modules []string
	prefix  string
}

func main() {
	req := parse()
	if req.version != "0" {
		fmt.Println("Ask bot admin to update me, This is V0 and request is " + req.version)
		return
	}
	dir := dirstr(req)
	conf := parseconf(dir)
	if strings.HasPrefix(req.CONTENT, conf.prefix+"help") {
		req_split := strings.SplitN(req.CONTENT, " ", 2)
		if len(req_split) != 2 {
			fmt.Print(">>> Welcome to Mashiron!\nIf you want to read help of each commands, type module name after this command!\nCurrently enabled commands are: ", conf.modules)
			return
		} else {
			path := "cmd/" + req_split[1] + "/README"
			_, err := os.Stat(path)
			if err == nil {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Println("> ERROR: " + err.Error())
				}
				lines := ">>> " + string(b)
				fmt.Print(lines)
			}
			if os.IsNotExist(err) {
				fmt.Println("> Module " + req_split[1] + " does not exists.")
				return
			}
		}
	}
}
func parse() Request {
	priv := strings.Split(os.Args[4], ",")
	req := Request{
		version: os.Args[1],
		API:     os.Args[3],
		ROOM:    os.Args[5],
		USER:    os.Args[6],
		PRIV:    priv,
		CONTENT: os.Args[7],
	}
	return req
}
func dirstr(req Request) Dir {
	roomdir := "data/" + req.API + "/" + req.ROOM + "/"
	cmddir := "cmd/man/"
	return Dir{
		roomdir:    roomdir,
		cmddatadir: roomdir + cmddir,
	}
}
func parseconf(dir Dir) Conf {
	c, err := ini.Load(dir.roomdir + "user.ini")
	if err != nil {
		fmt.Println(err)
	}
	return Conf{
		prefix:  c.Section("core").Key("prefix").String(),
		modules: c.Section("core").Key("module").Strings(" "),
	}
}
func check(req Request, privs []string) bool {
	if len(privs) == 0 {
		return true
	}
	for _, priv := range privs {
		for _, req_priv := range req.PRIV {
			if req_priv == priv {
				return true
			}
		}
	}
	return false
}
