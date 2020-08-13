package main

import (
	"fmt"
	"mashironsrv.visualstudio.com/Public/_git/Mashiron-go/mashiron"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ModuleName    = "vote"
	BacketInfo    = "info"
	BacketChoices = "choices"
	BacketVotes   = "votes"
)

func main() {
	req := mashiron.JSONToRequest(&os.Args[1])
	answer := ""
	if req.Version != 0 {
		fmt.Fprint(os.Stderr, "[vote.go] FATAL:Version error!")
		return
	}
	dir := mashiron.GetDirList(&req, ModuleName)
	conf := mashiron.GetCoreConf(&dir)
	os.MkdirAll(dir.CmdDataDir, 0777)
	if strings.HasPrefix(req.Content, conf.Prefix+ModuleName) {
		mashiron.DB_CreateRootBacket(ModuleName, &dir)
		arr := strings.Split(req.Content, " ")
		if strings.HasPrefix(req.Content, conf.Prefix+ModuleName+".start ") {
			if len(arr) < 4 {
				answer += "Not enough argument.\n"
			} else {
				mashiron.DB_AddBucket(ModuleName, &dir, BacketInfo, [][]string{
					{"title", arr[1]},
					{"time", time.Now().String()},
					{"author", req.User},
				})
				var choices [][]string
				for i, v := range arr[2:] {
					choices = append(choices, []string{strconv.Itoa(i), v})
				}
				mashiron.DB_AddBucket(ModuleName, &dir, BacketChoices, choices)
				answer += "Vote started!\nTitle: " + arr[1] + mashiron.KVArrayPager(true, true, ": ", "", 0, choices) + "\n"
				if len(arr) > 12 {
					answer += "Please note that choices over 10 are not listed. If you want to view more, Just type`" + conf.Prefix + ModuleName + ".status`.\n"
				}
			}
		} else if mashiron.DB_IfBucketExists(ModuleName, &dir, BacketInfo) {
			if strings.HasPrefix(req.Content, conf.Prefix+ModuleName+".end") {
				answer += "Vote ended! Here is the result!\n\n"
				answer += result(mashiron.DB_GetFullKVList(ModuleName,&dir,BacketVotes),mashiron.DB_GetFullKVList(ModuleName,&dir,BacketChoices))
				for _,v := range []string{BacketInfo,BacketVotes,BacketChoices} {
					mashiron.DB_DeleteBucket(ModuleName,&dir,"",v)
				}
			} else if strings.HasPrefix(req.Content, conf.Prefix+ModuleName+".status") {
				answer += "Current vote status:\n\n"
				answer += result(mashiron.DB_GetFullKVList(ModuleName,&dir,BacketVotes),mashiron.DB_GetFullKVList(ModuleName,&dir,BacketChoices))
			} else if strings.HasPrefix(req.Content, conf.Prefix+ModuleName+".list") {
				i := 0
				if len(arr) == 2 {
					t, err := strconv.Atoi(arr[1])
					if err == nil {
						i = t
					}
				}
				answer += mashiron.KVArrayPager(true, true, ": ", "`", i, mashiron.DB_GetFullKVList(ModuleName, &dir, BacketVotes)) + "\n"
			} else if len(arr) >= 2 {
				if mashiron.DB_IfBucketExists(ModuleName, &dir, BacketInfo) {
					_, err := strconv.Atoi(arr[1])
					if err != nil {
						answer += "Invaild request!\n"
					} else {
						chk := mashiron.DB_Regex(ModuleName,BacketVotes,req.User,&dir)
						if len(chk) != 0 {
							mashiron.DB_DeleteBucket(ModuleName,&dir,BacketVotes, req.User)
						}
						mashiron.DB_AddBucket(ModuleName, &dir, BacketVotes, [][]string{{req.User, arr[1]}})
						answer += "Voted!\n"
					}
				}
			}
		} else {
			answer += "Vote not found!\n"
		}
	}
	fmt.Print(mashiron.ResultToJSON(&mashiron.Result{Content: answer}))
}

func result(votes [][]string, choices [][]string) string {
	res := ""
	for _,v := range choices {
		a,_ := strconv.Atoi(v[0])
		res += "No." + v[0] + " (" + v[1] + ") : " + strconv.Itoa(calc(a,votes)) + " votes\n"
	}
	return res
}

func calc(num int,votes [][]string) int {
	i := 0
	for _,v := range votes {
		vote, _ := strconv.Atoi(v[1])
		if num == vote {
			i++
		}
	}
	return i
}