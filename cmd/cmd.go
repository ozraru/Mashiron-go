package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/otiai10/copy"
	"mashironsrv.visualstudio.com/Public/_git/Mashiron-go/mashiron"
)

func main() {
	req := mashiron.JSONToRequest(&os.Args[1])
	dir := mashiron.GetDirList(&req,"")
	var ans []string
	if _, err := os.Stat(dir.RoomDir); os.IsNotExist(err) {
		ans = append(ans,setup(&dir))
	} else {
		conf := mashiron.GetCoreConf(&dir)
		if mashiron.CheckPrivileges(&req, &conf.Privileges_Run) {
			if len(conf.Hooks) != 0 {
				for _, hook := range conf.Hooks {
					req.Ishook = true
					ans = append(ans,start(&req, &dir, true, hook))
				}
			}
			if len(conf.Modules) != 0 {
				if strings.HasPrefix(req.Content, conf.Prefix) {
					for _, module := range conf.Modules {
						if strings.HasPrefix(req.Content, conf.Prefix+module) {
							fl, _ := ioutil.ReadDir(dir.CmdDir)
							for _, f := range  fl {
								if module == f.Name() {
									req.Ishook = false
									ans = append(ans,start(&req, &dir, false, module))
								}
							}
						}
					}
				}
			}
		}
	}
	fmt.Print(mashiron.ResultsToJSON(&ans))
}

//create room data dirs and copy skel dir
func setup(dir *mashiron.Dir) string {
	content := "Mashiron setup:core\n"
	err := os.MkdirAll(dir.CmdDataDir, 500)
	if err != nil {
		content += "WTF>>Can't create data dir...Aborting."
		fmt.Fprint(os.Stderr, err.Error())
	} else {
		content += "DONE>>Created data dir."
		err = copy.Copy(dir.SkeletonDir, dir.RoomDir)
		if err != nil {
			content += "WTF>>Can't copy default files...Aborting"
			fmt.Fprint(os.Stderr, err.Error())
		} else {
			content += "DONE>>Copied default files."
			content += "DONE>>Core setup completed."
		}
	}
	return mashiron.ResultToJSON(&mashiron.Result{
		Attachments: nil,
		Content:     content,
	})
}

//start process
func start(req *mashiron.Request, dir *mashiron.Dir, ishook bool, cmd string) string {
	if strings.Contains(cmd, "..") {
		//For security
		return mashiron.ResultToJSON(&mashiron.Result{
			Attachments: nil,
			Content:     "",
		})
	}
	return mashiron.ExecModule(req,dir,ishook,cmd)
}