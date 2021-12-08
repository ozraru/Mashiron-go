package main

import (
	"bytes"
	"fmt"
	"mashironsrv.visualstudio.com/Public/_git/Mashiron-go/mashiron"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)
type Conf struct {
	priv_run     []string
	priv_edit    []string
	priv_regex   []string
	priv_admin   []string
	global_cache bool
	prefix       string
}
type Cmd struct {
	name   string
	author string
	time   string
	file   string
	cache  string
}

var defKeys = []string{"author", "time", "file", "cache"}
const ModuleName = "sh"

func main() {
	req := mashiron.JSONToRequest(&os.Args[1])
	if req.Version != 0 {
		fmt.Fprint(os.Stderr, "[sh.go] FATAL:version error!")
		return
	}
	dir := mashiron.GetDirList(&req,ModuleName)
	conf := parseconf(&dir)
	if mashiron.CheckPrivileges(&req, &conf.priv_run) {
		cmd(&req, &conf, &dir)
	}
}
func parseconf(dir *mashiron.Dir) Conf {
	c, err := ini.Load(dir.RoomDir + "user.ini")
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	c2, err := ini.Load("mashiron.ini")
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return Conf{
		priv_edit:    c.Section(ModuleName).Key("priv_conf").Strings(" "),
		priv_run:     c.Section(ModuleName).Key("priv_run").Strings(" "),
		priv_regex:   c.Section(ModuleName).Key("priv_regex").Strings(" "),
		priv_admin:   c.Section(ModuleName).Key("priv_admin").Strings(" "),
		global_cache: c2.Section(ModuleName).Key("cache").MustBool(),
		prefix:       c.Section("core").Key("prefix").String(),
	}
}

func cmd(req *mashiron.Request, conf *Conf, dir *mashiron.Dir) {
	answer := ""
	os.MkdirAll(dir.CmdDataDir, 0777)
	for _, v := range []string{"hook", "cmd", "cache"} {
		mashiron.DB_CreateRootBacket(v, dir)
	}
	if req.Ishook {
		for _, i := range mashiron.DB_Regex("hook","",req.Content,dir) {
			//run
			if mashiron.DB_IfBucketExists("cmd",dir,i) {
				b := mashiron.DB_GetBucket("cmd", dir, i, defKeys)
				answer += vm(&req.Content, dir, &Cmd{
					name:   i,
					author: b[0],
					time:   b[1],
					file:   b[2],
					cache:  b[3],
				}, options(req))
			} else {
				answer += "No such script in database.\n"
			}
		}
	} else if strings.HasPrefix(req.Content, conf.prefix+ModuleName+".") {
		//someone calls me
		if mashiron.CheckPrivileges(req, &conf.priv_edit) {
			if strings.HasPrefix(req.Content, conf.prefix+ModuleName+".add ") {
				//add script
				addexp, _ := regexp.Compile("(?s)" + conf.prefix + ModuleName + ".add ([^ ]*?)( --no-cache)?( --overwrite)?\n?```(bash\n|sh\n|\n)?(.*)```") // 1:name, 5:script
				matched := addexp.FindStringSubmatch(req.Content)
				if matched == nil {
					answer += "> Cannot find script...? (Syntax Error)"
				} else if matched[1] == "" {
					answer += "> Please include script name before file."
				} else if matched[3] == "" && mashiron.DB_IfBucketExists("cmd", dir, matched[1]) {
					answer += "> Script already exists."
				} else {
					cmd := Cmd{
						author: req.User,
						cache:  "true",
						time:   time.Now().String(),
						name:   matched[1],
					}
					if matched[2] == " --no-cache" {
						cmd.cache = "false"
					} else {
						cmd.cache = "true"
					}
					cmd.file = matched[5]
					out, _ := exec.Command(dir.CmdDir+"shchk.sh", cmd.file).Output()
					answer += string(out)
					mashiron.DB_AddBucket("cmd", dir, cmd.name, [][]string{
						{"author", cmd.author},
						{"cache", cmd.cache},
						{"time", cmd.time},
						{"name", cmd.name},
						{"file", cmd.file},
					})
					answer += "> Added script. Type `" + conf.prefix + "sh.info " + cmd.name + "` for details."
				}
			}
			if strings.HasPrefix(req.Content, conf.prefix+ModuleName+".rm ") {
				//delete script
				req_split := strings.SplitN(req.Content, " ", 2)
				if len(req_split) != 2 {
					answer += "> Request split error."
				} else if mashiron.DB_IfBucketExists("cmd", dir, req_split[1]) {
					i := mashiron.DB_GetBucket("cmd", dir, req_split[1], defKeys)
					info := Cmd{
						name:   req_split[1],
						author: i[0],
						time:  i[1],
						file:   i[2],
						cache:   i[3],
					}
					if info.author == req.User || mashiron.CheckPrivileges(req, &conf.priv_admin) {
						mashiron.DB_DeleteBucket("cmd", dir, "", req_split[1])
						answer += "> Deleted `" + info.name + "` ."
					} else {
						answer += "> You are not allowed to delete this command."
					}
				} else {
					answer += "> No such command in database."
				}
			}
			if strings.HasPrefix(req.Content, conf.prefix+ModuleName+".info ") {
				req_split := strings.SplitN(req.Content, " ", 2)
				if len(req_split) != 2 {
					answer += "> Request split error."
				} else if mashiron.DB_IfBucketExists("cmd", dir, req_split[1]) {
					i := mashiron.DB_GetBucket("cmd", dir, req_split[1], defKeys)
					info := Cmd{
						name:   req_split[1],
						author: i[0],
						time:  i[1],
						file:   i[2],
						cache:   i[3],
					}
					answer += ">>> Name: `" + info.name + "`\nBy: `" + info.author + "`\nAt: `" + info.time + "`\nCache: `" + info.cache + "`\n File:```sh\n" + info.file + "```\n"
				} else {
					answer += "> No such script in database."
				}
			}
			if strings.HasPrefix(req.Content, conf.prefix+ModuleName+".ls") {
				list := mashiron.DB_GetBucketList("cmd", dir)
				if len(list) == 0 {
					answer += "> There are no script in database."
				} else {
					answer += "> There are " + strconv.Itoa(len(list)) + " command(s) in database.\n"
					page, err := strconv.Atoi(strings.TrimLeft(req.Content, conf.prefix+"sh.ls "))
					if err != nil {
						page = 0
					}
					answer += mashiron.KVArrayPager(true, false, " => ", "`", page, list)
				}
			}
			if mashiron.CheckPrivileges(req, &conf.priv_regex) {
				if strings.HasPrefix(req.Content, conf.prefix+ModuleName+".hook.add ") {
					//add hook regex
					req_split := strings.SplitN(req.Content, " ", 3)
					if len(req_split) != 3 {
						answer += "> Request split error."
					} else if mashiron.DB_IfBucketExists("cmd", dir, req_split[2]) {
						_, err := regexp.Compile(req_split[1])
						if err != nil {
							answer += "> Regex error.\n" + err.Error()
						} else {
							mashiron.DB_AddBucket("hook",dir,"",[][]string{{req_split[1], req_split[2]}})
							answer += ">>> Added to DB.\nRegEx: `" + req_split[1] + "`\nCmd: `" + req_split[2] + "`"
						}
					} else {
						answer += "> No such script in database."
					}
				}
				if strings.HasPrefix(req.Content, conf.prefix+ModuleName+".hook.rm ") {
					//delete hook regex
					req_split := strings.SplitN(req.Content, " ", 2)
					if len(req_split) != 2 {
						answer += "> Request split error."
					} else if mashiron.DB_IfBucketExists("hook",dir,req_split[1]) {
						mashiron.DB_DeleteBucket("hook",dir, "",req_split[1])
						answer += "> Deleted `" + req_split[1] + "`."
					} else {
						answer += "> No such regex in database."
					}
				}
				if strings.HasPrefix(req.Content, conf.prefix+ModuleName+".hook.ls") {
					//hook listing
					list := mashiron.DB_GetBucketList("hook", dir)
					if len(list) == 0 {
						answer += "> There are no regex in database."
					} else {
						answer += "> There are " + strconv.Itoa(len(list)) + " regex(s) in database.\n"
						i,err := strconv.Atoi(strings.TrimLeft(req.Content, conf.prefix+ModuleName+".hook.ls "))
						if err != nil {
							i = 0
						}
						answer += mashiron.KVArrayPager(true,true," => ","`",i,list)
					}
				}
			}
			if mashiron.CheckPrivileges(req, &conf.priv_run) {
				req_split := strings.SplitN(req.Content, " ", 2)
				req_cmd := strings.TrimPrefix(req_split[0],conf.prefix+ModuleName+".")
				if mashiron.DB_IfBucketExists("cmd", dir, req_cmd) {
					i := mashiron.DB_GetBucket("cmd", dir, req_cmd, defKeys)
					info := Cmd{
						name:   req_cmd,
						author: i[0],
						time:  i[1],
						file:   i[2],
						cache:   i[3],
					}
					v := ""
					if len(req_split) > 1 {
						v = req_split[1]
					}
					answer += vm(&v, dir, &info, options(req))
				}
			}
		}
	}
	fmt.Print(mashiron.ResultToJSON(&mashiron.Result{
		Attachments: nil,
		Content:     answer,
	}))
}

func vm(req *string, dir *mashiron.Dir, cmd *Cmd, opts string) string{
	//Systemd-nspawn needs root priv.
	c := exec.Command("sudo", append([]string{dir.CmdDir + "run.sh", cmd.file, opts}, strings.Split(*req, " ")...)...)
	var stdOut bytes.Buffer
	c.Stdout = &stdOut
	c.Stderr = &stdOut
	err := c.Run()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	} else {
		return stdOut.String()
	}
	return ""
}

//Convert mashiron.Request opts to string
func options(req *mashiron.Request) string {
	r := ""
	for _,v :=  range req.Options {
		r += "export MSR_" + strings.ToUpper(req.Api + "_" + v[0]) + "=" + "'" + v[1] + "'\n"
	}
	fmt.Fprint(os.Stderr, r)
	return r
}
