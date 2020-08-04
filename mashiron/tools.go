package mashiron

import (
	"encoding/json"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strconv"
	"strings"
)

//Type definitions
type Request struct {
	Api        string `json:"api"`
	Room       string `json:"room"`
	User       string `json:"user"`
	Privileges []string `json:"privileges"`
	Content    string `json:"content"`
	Options	[][]string `json:"options"`

	//These are for core to each commands.
	Version int
	Ishook  bool
}

type Result struct {
	Attachments []string `json:"attachments"`
	Content    string `json:"content"`
}
type Results struct {
	Jsons []string `json:"jsons"`
}

type Dir struct {
	BaseDir     string
	DataDir     string
	ApiDir      string
	RoomDir     string
	CmdDataDir  string
	SkeletonDir string
	CmdDir      string
}

type CoreConf struct {
	Prefix   string
	Modules  []string
	Hooks    []string
	Privileges_Run  []string
	Privileges_Conf []string
}

func JSONToRequest(data *string) Request {
	var request Request
	if err := json.Unmarshal([]byte(*data), &request); err != nil {
		fmt.Fprint(os.Stderr, "Request parse error:" + err.Error())
	}
	return request
}

func RequestToJSON(request *Request) string {
	result, err := json.Marshal(request)
	if err != nil {
		fmt.Fprint(os.Stderr, "Request encode error:" + err.Error())
	}
	return string(result)
}

func JSONToResult(data *string) Result {
	var result Result
	if err := json.Unmarshal([]byte(*data), &result); err != nil {
		fmt.Fprint(os.Stderr, "Request parse error:" + err.Error())
	}
	return result
}

func ResultToJSON(result *Result) string {
	str, err := json.Marshal(result)
	if err != nil {
		fmt.Fprint(os.Stderr, "Request encode error:" + err.Error())
	}
	return string(str)
}

func JSONToResults(data *string) []string {
	var results Results
	if err := json.Unmarshal([]byte(*data), &results); err != nil {
		fmt.Fprint(os.Stderr, "Request parse error:" + err.Error())
	}
	return results.Jsons
}

func ResultsToJSON(results *[]string) string {
	str, err := json.Marshal(Results{Jsons:*results})
	if err != nil {
		fmt.Fprint(os.Stderr, "Request encode error:" + err.Error())
	}
	return string(str)
}

func GetDirList(req *Request,CommandName string) Dir {
	basedir := "./"
	datadir := basedir + "data/"
	apidir := datadir + req.Api + "/"
	roomdir := apidir + req.Room + "/"
	skeldir := basedir + "skel/"
	cmddatadir := roomdir + "cmd/" + CommandName + "/"
	cmddir := basedir + "cmd/" + CommandName + "/"
	return Dir{
		BaseDir:    basedir,
		DataDir:    datadir,
		ApiDir:     apidir,
		RoomDir:    roomdir,
		CmdDataDir: cmddatadir,
		SkeletonDir:    skeldir,
		CmdDir:     cmddir,
	}
}

func Old_CmdParser() Request {
	ishook, _ := strconv.ParseBool(os.Args[2])
	priv := strings.Split(os.Args[4], ",")
	version,_ := strconv.Atoi(os.Args[1])
	req := Request{
		Version: version,
		Ishook:  ishook,
		Api:     os.Args[3],
		Room:    os.Args[5],
		User:    os.Args[6],
		Privileges:    priv,
		Content: os.Args[7],
	}
	return req
}

//check privs
func CheckPrivileges(req *Request, privs *[]string) bool {
	if len(*privs) == 0 {
		return true
	}
	for _, priv := range *privs {
		for _, req_priv := range req.Privileges {
			if req_priv == priv {
				return true
			}
		}
	}
	return false
}

//KV means Key-Value
//[][]string{{"a","b"},{"a","b"},{"a","b"}...}
func KVArrayPager(IfPrintKey bool,IfPrintValue bool,separator string,wrapper string,page int,array [][]string) string{
	var pages [][][]string
	result := "\n[Pager] Current page (Starts from 0) : " + strconv.Itoa(page) + "\n\n"
	for i := 0; i < len(array); i += 10 {
		end := i + 10
		if len(array) < end {
			end = len(array)
		}
		pages = append(pages, array[i:end])
	}
	if len(pages)-1 < page {
		result += "<NOTHING>\n"
	} else {
		for _,kv := range pages[page] {
			if IfPrintKey && IfPrintValue {
				result += wrapper + kv[0] + wrapper + separator + wrapper + kv[1]
			} else if IfPrintKey {
				result += wrapper + kv[0]
			} else if IfPrintValue {
				result += wrapper + kv[1]
			}
			result += wrapper + "\n"
		}
	}
	return result
}

//read config
func GetCoreConf(dir *Dir) CoreConf {
	c, _ := ini.Load(dir.RoomDir + "user.ini")
	return CoreConf{
		Prefix:   c.Section("core").Key("prefix").String(),
		Modules:  c.Section("core").Key("module").Strings(" "),
		Hooks:    c.Section("core").Key("hook").Strings(" "),
		Privileges_Conf: c.Section("core").Key("priv_conf").Strings(" "),
		Privileges_Run:  c.Section("core").Key("priv_run").Strings(" "),
	}
}
