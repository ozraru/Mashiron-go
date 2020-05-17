package mashiron

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//To add new command,just create new const value like below and define command...
const (
	Command_sh       string = "sh"
	Command_conf     string = "conf"
	Command_help     string = "help"
	Command_man      string = "man"
	Command_ping     string = "ping"
	Command_splatoon string = "splatoon"
	Command_weather  string = "weather"
)

func ExecModule(request *Request, dir *Dir, ishook bool, cmd string) string {
	switch cmd {
	//non-JSON request,Returns non-JSON output (Classic)
	case Command_ping, Command_splatoon, Command_weather:
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
		r := run(dir.CmdDir+cmd+"/"+cmd, []string{"0", strconv.FormatBool(ishook), request.Api, strings.Join(request.Privileges, ","), request.Room, request.User, request.Content})
		return ResultToJSON(&Result{
			Attachments: nil,
			Content:     r,
		})
	case Command_sh, Command_conf, Command_help, Command_man:
		r := run(dir.CmdDir+cmd+"/"+cmd, []string{RequestToJSON(request)})
		return r
	default:
		fmt.Fprint(os.Stderr, "FATAL> Command has been called but there are no definition in modules.go!")
		return ResultToJSON(&Result{
			Attachments: nil,
			Content:     "> Fatal error occurred.Please contact to bot administrator!",
		})
	}
}

func run(cmd string, args []string) string {
	cmdrun := exec.Command(cmd, args...)
	var stdOut bytes.Buffer
	cmdrun.Stdout = &stdOut
	cmdrun.Stderr = os.Stderr
	cmdrun.Run()
	return stdOut.String()
}
