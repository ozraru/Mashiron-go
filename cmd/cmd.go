package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/otiai10/copy"
	"gopkg.in/ini.v1"
)

//Define objs
type Request struct {
	API     string
	ROOM    string
	USER    string
	PRIV    string
	CONTENT string
}
type Dir struct {
	basedir    string
	datadir    string
	apidir     string
	roomdir    string
	cmddatadir string
	skeldir    string
	cmddir     string
}
type Confdata struct {
	prefix    string
	modules   []string
	hooks     []string
	priv_run  []string
	priv_conf []string
}

func main() {
	req := parse()
	dir := dirstr(&req)
	if _, err := os.Stat(dir.roomdir); os.IsNotExist(err) {
		fmt.Println("Setup: mashiron core...")
		setup(&dir)
		return
	} else {
		conf := readconf(&dir)
		if check(&req, &conf.priv_run) {
			cmd(&req, &dir, &conf)
		}
	}
}

//read cmdline flags
func parse() Request {
	/* flag.arg(0):Program name
	 * flag.arg(1):Req API
	 * flag.arg(2):Req chat room id
	 * flag.arg(3):Who?
	 * flag.arg(4):Priv of flag.arg(3)
	 * flag.arg(5):Content
	 */
	return Request{
		API:     os.Args[1],
		ROOM:    os.Args[2],
		USER:    os.Args[3],
		PRIV:    os.Args[4],
		CONTENT: os.Args[5],
	}
}

//check privs
func check(req *Request, privs *[]string) bool {
	//PRIV_LIST := []string{"PRIV_RUN_CMD", "PRIV_CHANGE_CONFIG"}
	if len(*privs) == 0 {
		return true
	}
	for _, priv := range *privs {
		for _, req_priv := range strings.Split(req.PRIV, ",") {
			if req_priv == priv {
				return true
			}
		}
	}
	return false
}

//gen dir strs
func dirstr(req *Request) Dir {
	basedir := "./"
	datadir := basedir + "data/"
	apidir := datadir + req.API + "/"
	roomdir := apidir + req.ROOM + "/"
	cmddatadir := roomdir + "cmd/"
	skeldir := basedir + "skel/"
	cmddir := basedir + "cmd/"
	return Dir{
		basedir:    basedir,
		datadir:    datadir,
		apidir:     apidir,
		roomdir:    roomdir,
		cmddatadir: cmddatadir,
		skeldir:    skeldir,
		cmddir:     cmddir,
	}
}

//create room data dirs and copy skel dir
func setup(dir *Dir) bool {
	fmt.Println("Mashiron setop:core")
	err := os.MkdirAll(dir.cmddatadir, 500)
	if err != nil {
		fmt.Println("WTF>>Can't create data dir...Aborting.")
		fmt.Println(err)
		return false
	} else {
		fmt.Println("DONE>>Created data dir.")
	}
	err = copy.Copy(dir.skeldir, dir.roomdir)
	if err != nil {
		fmt.Println("WTF>>Can't copy default files...Aborting")
		fmt.Println(err)
		return false
	} else {
		fmt.Println("DONE>>Copied default files.")
	}
	fmt.Println("DONE>>Core setup completed.")
	return true
}

//read config
func readconf(dir *Dir) Confdata {
	c, _ := ini.Load(dir.roomdir + "user.ini")
	return Confdata{
		prefix:    c.Section("core").Key("prefix").String(),
		modules:   c.Section("core").Key("module").Strings(" "),
		hooks:     c.Section("core").Key("hook").Strings(" "),
		priv_conf: c.Section("core").Key("priv_conf").Strings(" "),
		priv_run:  c.Section("core").Key("priv_run").Strings(" "),
	}
}

//hook
func cmd(req *Request, dir *Dir, conf *Confdata) {
	if len(conf.hooks) != 0 {
		for _, hook := range conf.hooks {
			start(req, dir, conf, true, hook)
		}
	}
	if len(conf.modules) != 0 {
		if strings.HasPrefix(req.CONTENT, conf.prefix) {
			for _, module := range conf.modules {
				if strings.HasPrefix(req.CONTENT, conf.prefix+module) {
					fl, _ := ioutil.ReadDir(dir.cmddir)
					for _, f := range fl {
						if module == f.Name() {
							start(req, dir, conf, false, module)
						}
					}
				}
			}
		}
	}
}

//start process
func start(req *Request, dir *Dir, conf *Confdata, ishook bool, cmd string) {
	/* Args
	0:cmd
	1:ver
	2:ishook
	3:API
	4:PRIV
	5:ROOM
	6:USER
	7:CONTENT
	*/
	if strings.Contains(cmd, "..") {
		//For security
		return
	}
	cmdrun := exec.Command(dir.cmddir+cmd+"/"+cmd, "0", strconv.FormatBool(ishook), req.API, req.PRIV, req.ROOM, req.USER, req.CONTENT)
	cmdrun.Stdout = os.Stdout
	cmdrun.Stderr = os.Stdout
	cmdrun.Run()
}
