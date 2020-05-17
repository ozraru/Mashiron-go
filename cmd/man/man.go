package main

import (
	"fmt"
	"mashironsrv.visualstudio.com/Public/_git/Mashiron-go/mashiron"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type Conf struct {
	priv_read  []string
	priv_edit  []string
	priv_admin []string
	prefix     string
}
type Man struct {
	name   string
	author string
	time   string
	file   string
}

func main() {
	req := mashiron.JSONToRequest(&os.Args[1])
	if req.Version != 0 {
		fmt.Fprint(os.Stderr, "[man.go] FATAL:Version error!")
		return
	}
	dir := mashiron.GetDirList(&req,"man")
	conf := parseconf(dir)
	if mashiron.CheckPrivileges(&req, &conf.priv_read) {
		cmd(req, conf, dir)
	}
}

func parseconf(dir mashiron.Dir) Conf {
	c, err := ini.Load(dir.RoomDir + "user.ini")
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return Conf{
		priv_edit:  c.Section("man").Key("priv_edit").Strings(" "),
		priv_read:  c.Section("man").Key("priv_read").Strings(" "),
		priv_admin: c.Section("man").Key("priv_admin").Strings(" "),
		prefix:     c.Section("core").Key("prefix").String(),
	}
}

func cmd(req mashiron.Request, conf Conf, dir mashiron.Dir) {
	answer := ""
	os.MkdirAll(dir.CmdDataDir, 0777)
	mashiron.DB_CreateRootBacket("man", &dir)
	if strings.HasPrefix(req.Content, conf.prefix+"man.") {
		//someone calls me
		if mashiron.CheckPrivileges(&req, &conf.priv_edit) {

			if strings.HasPrefix(req.Content, conf.prefix+"man.add") {
				//add command
				req_split := strings.SplitN(req.Content, " ", 4)
				if strings.Contains(req_split[1], "\n") {
					answer += "> Please include man name before file."
				}else if mashiron.DB_IfBucketExists("man",&dir,req_split[1]) {
					answer += "> Man already exists."
				} else {
					man := Man{
						author: req.User,
						time:   time.Now().String(),
						name:   req_split[1],
					}
					//trim
					cmdtmp := strings.SplitN(req.Content, "\n", 5)
					c := 0
					for {
						if strings.HasSuffix(cmdtmp[c], "```") {
							man.file = strings.Join(cmdtmp[c+1:], "")
							break
						} else if c == 5 {
							answer += "> Cannot find man...?"
						} else {
							c++
						}
					}
					man.file = strings.TrimRight(man.file, "```")
					mashiron.DB_AddBucket("man",&dir,man.name,[][]string{
						{"author",man.author},
						{"time",man.time},
						{"file",man.file},
					})
					answer += "> Added man. Type `" + conf.prefix + "man." + man.name + "` for details."
				}
			} else if strings.HasPrefix(req.Content, conf.prefix+"man.rm ") {
				//delete command
				req_split := strings.SplitN(req.Content, " ", 2)
				if len(req_split) != 2 {
					answer += "> Request split error."
				} else if mashiron.DB_IfBucketExists("root",&dir,req_split[1]) {
					info_tmp := mashiron.DB_GetBucket("man",&dir,req_split[1],[]string{"author","time","file"})
					info := Man{
						name: req_split[1],
						author:info_tmp[0],
						time: info_tmp[1],
						file: info_tmp[2],
					}
					if info.author == req.User || mashiron.CheckPrivileges(&req, &conf.priv_admin) {
						mashiron.DB_DeleteBucket("man",&dir,req_split[1])
						answer += "> Deleted `" + info.name + "` ."
					} else {
						answer += "> You are not allowed to delete this command."
					}
				} else {
					answer += "> No such command in database."
				}
			}

			if strings.HasPrefix(req.Content, conf.prefix+"man.ls") {
				list := mashiron.DB_GetBucketList("man",&dir)
				if len(list) == 0 {
					answer += "> There are no man in database."
				} else {
					answer += "> There are " + strconv.Itoa(len(list)) + " man(s) in database.\n"
					req_split := strings.SplitN(req.Content, " ", 2)
					page := 0
					page, err := strconv.Atoi(req_split[1])
					if err != nil {
						fmt.Fprint(os.Stderr, err.Error())
						page = 0
					}
					answer += mashiron.KVArrayPager(true,false," => ","`",page,list)
				}
			}

			if mashiron.CheckPrivileges(&req, &conf.priv_read) {
				req_split := strings.SplitN(req.Content, " ", 2)
				c := strings.Split(req_split[0], ".")
				req_cmd := c[len(c)-1]
				if mashiron.DB_IfBucketExists("man",&dir,req_cmd) {
					r := mashiron.DB_GetBucket("man",&dir,req_cmd, []string{"author","time","file"})
					i := Man{
						name:   req_cmd,
						author: r[0],
						time:   r[1],
						file:   r[2],
					}
					answer += ">>> Man file `" + i.name + "` \n`Author`: `" + i.author + "`\n`Time`: `" + i.time + "`\nMan:```\n" + i.file + "```"
				}
			}
		}
		fmt.Print(mashiron.ResultToJSON(&mashiron.Result{
			Attachments: nil,
			Content:     answer,
		}))
	}
}