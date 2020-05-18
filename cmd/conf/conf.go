package main

import (
	"fmt"
	"io/ioutil"
	"mashironsrv.visualstudio.com/Public/_git/Mashiron-go/mashiron"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Conf struct {
	priv_conf []string
	prefix    string
}

func main() {
	req := mashiron.JSONToRequest(&os.Args[1])
	if req.Version != 0 {
		fmt.Fprint(os.Stderr,"Version mismatch!")
		return
	}
	dir := mashiron.GetDirList(&req,"conf")
	conf := parseconf(dir)
	if strings.HasPrefix(req.Content, conf.prefix+"conf") {
		var ans string
		req_split := strings.SplitN(req.Content, " ", 5)
		if len(req_split) > 1 && mashiron.CheckPrivileges(&req, &conf.priv_conf) {
			c, err := ini.Load(dir.RoomDir + "user.ini")
			if err != nil {
				fmt.Fprint(os.Stderr,err.Error())
			}
			if req_split[1] == "get" || len(req_split) == 4 {
				s := c.Section(req_split[2]).Key(req_split[3]).String()
				if s == ""  {
					ans += "> There were some errors while getting conf.(Maybe not found?)\n"
				} else {
					ans += "> " + s + "\n"
				}
			}
			if req_split[1] == "set" {
				v := strings.Join(req_split[4:], " ")
				c.Section(req_split[2]).Key(req_split[3]).SetValue(v)
				err := c.SaveTo(dir.RoomDir + "user.ini")
				if err != nil {
					fmt.Fprint(os.Stderr,err.Error())
				} else {
					ans += "> Set OK\n"
				}
			}
			if req_split[1] == "group" || len(req_split) == 4 {
				if req_split[2] == "add" {
					_, err := c.NewSection(req_split[3])
					if err != nil {
						fmt.Fprint(os.Stderr,err.Error())
					} else {
						ans+="Created " + req_split[3] + " .\n"
					}
				}
				if req_split[1] == "rm" {
					c.DeleteSection(req_split[2])
					ans+="Deleted " + req_split[2] + " .\n"
				}
			}
			if req_split[1] == "default" {
				path := "skel/user.ini"
				_, err := os.Stat(path)
				if err == nil {
					b, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Fprint(os.Stderr,err.Error())
					}
					ans+=">>> " + string(b)
				}
			}
			if req_split[1] == "cat" {
				path := dir.RoomDir + "user.ini"
				_, err := os.Stat(path)
				if err == nil {
					b, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Fprint(os.Stderr,err.Error())
					}
					ans += ">>> " + string(b)
				}
			}
		} else {
			ans+="> This is Mashiron conf module. read help for details, type " + conf.prefix + "conf default for examples."
		}
		fmt.Print(mashiron.ResultToJSON(&mashiron.Result{
			Attachments: nil,
			Content:     ans,
		}))
	}
}
func parseconf(dir mashiron.Dir) Conf {
	c, err := ini.Load(dir.RoomDir + "user.ini")
	if err != nil {
		fmt.Println(err)
	}
	return Conf{
		prefix:    c.Section("core").Key("prefix").String(),
		priv_conf: c.Section("core").Key("priv_conf").Strings(" "),
	}
}
