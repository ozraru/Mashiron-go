package main

import (
	"fmt"
	"io/ioutil"
	"mashironsrv.visualstudio.com/Public/_git/Mashiron-go/mashiron"
	"os"
	"strings"
)

func main() {
	req := mashiron.JSONToRequest(&os.Args[1])
	answer:= ""
	var attachment string
	if req.Version != 0 {
		fmt.Fprint(os.Stderr, "[help.go] FATAL:Version error!")
		return
	}
	dir := mashiron.GetDirList(&req,"help")
	conf := mashiron.GetCoreConf(&dir)
	if strings.HasPrefix(req.Content, conf.Prefix+"help") {
		req_split := strings.SplitN(req.Content, " ", 2)
		if len(req_split) != 2 {
			answer += ">>> Welcome to Mashiron!\nIf you want to read help of each commands, type module name after this command!\nCurrently enabled commands are: "+ strings.Join(conf.Modules,",")
			attachment = dir.CmdDir + "mashiron.png"
		} else {
			if !(strings.Contains(req_split[1], "..") || strings.Contains(req_split[1], "/")) {
				path := "cmd/" + req_split[1] + "/README"
				_, err := os.Stat(path)
				if err == nil {
					b, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Fprint(os.Stderr, "[help.go] ERROR:" + err.Error())
					}
					answer += string(b)
				}
				if os.IsNotExist(err) {
					answer += "> Module `" + req_split[1] + "` does not exists."
				}
			}
		}
	}
	fmt.Print(mashiron.ResultToJSON(&mashiron.Result{
		Attachments: []string{attachment},
		Content:     answer,
	}))
}
