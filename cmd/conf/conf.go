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
	priv_conf []string
	prefix    string
}

func main() {
	req := parse()
	if req.version != "0" {
		fmt.Println("Ask bot admin to update me, This is V0 and request is " + req.version)
		return
	}
	dir := dirstr(req)
	conf := parseconf(dir)
	if strings.HasPrefix(req.CONTENT, conf.prefix+"conf") {
		req_split := strings.SplitN(req.CONTENT, " ", 5)
		if len(req_split) > 1 && check(req, conf.priv_conf) {
			c, err := ini.Load(dir.roomdir + "user.ini")
			if err != nil {
				fmt.Println(err)
			}
			if req_split[1] == "get" || len(req_split) == 4 {
				s := c.Section(req_split[2]).Key(req_split[3]).String()
				if s == "" {
					fmt.Println("> There were some errors while getting conf.(Mayebe not found?)")
				} else {
					fmt.Println("> " + s)
				}
			}
			if req_split[1] == "set" {
				v := strings.Join(req_split[4:], " ")
				c.Section(req_split[2]).Key(req_split[3]).SetValue(v)
				err := c.SaveTo(dir.roomdir + "user.ini")
				if err != nil {
					fmt.Println("> " + err.Error())
				} else {
					fmt.Println("> Set OK")
				}
			}
			if req_split[1] == "group" || len(req_split) == 4 {
				if req_split[2] == "add" {
					_, err := c.NewSection(req_split[3])
					if err != nil {
						fmt.Println(err.Error())
					} else {
						fmt.Println("Created " + req_split[3] + " .")
					}
				}
				if req_split[1] == "rm" {
					c.DeleteSection(req_split[2])
					fmt.Println("Deleted " + req_split[2] + " .")
				}
			}
			if req_split[1] == "default" {
				path := "skel/user.ini"
				_, err := os.Stat(path)
				if err == nil {
					b, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Println("> ERROR: " + err.Error())
					}
					lines := ">>> " + string(b)
					fmt.Print(lines)
				}
			}
			if req_split[1] == "cat" {
				path := dir.roomdir + "user.ini"
				_, err := os.Stat(path)
				if err == nil {
					b, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Println("> ERROR: " + err.Error())
					}
					lines := ">>> " + string(b)
					fmt.Print(lines)
				}
			}
		} else {
			fmt.Println("> This is Mashiron conf module. read man for details, type " + conf.prefix + "conf default for examples.")
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
		prefix:    c.Section("core").Key("prefix").String(),
		priv_conf: c.Section("core").Key("priv_conf").Strings(" "),
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
