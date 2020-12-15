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
		autodelete := false
		if strings.HasPrefix(req.Content, conf.Prefix+ModuleName+".start ") {
			if len(arr) < 4 {
				answer += "Not enough argument.\n"
			}else if mashiron.DB_IfBucketExists(ModuleName, &dir, BacketInfo) {
				answer += "Another vote is currently ongoing!\n"
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
				answer += "Vote started!\nTitle: `" + arr[1] + "`\n" + mashiron.KVArrayPager(true, true, ": ", "`", 0, choices) + "\n"
				if len(arr) > 12 {
					answer += "Please note that choices over 10 are not listed. If you want to view more, Just type`" + conf.Prefix + ModuleName + ".status`.\n"
				}
			}
		} else if mashiron.DB_IfBucketExists(ModuleName, &dir, BacketInfo) {
			if strings.HasPrefix(req.Content, conf.Prefix+ModuleName+".end") {
				answer += "Vote ended! Here is the result!\nTitle: `" + mashiron.DB_GetBucket(ModuleName,&dir,BacketInfo,[]string{"title"})[0] + "`\n\n"
				answer += result(mashiron.DB_GetFullKVList(ModuleName,&dir,BacketVotes),mashiron.DB_GetFullKVList(ModuleName,&dir,BacketChoices))
				for _,v := range []string{BacketInfo,BacketVotes,BacketChoices} {
					mashiron.DB_DeleteBucket(ModuleName,&dir,"",v)
				}
			} else if strings.HasPrefix(req.Content, conf.Prefix+ModuleName+".status") {
				answer += "Current vote status:\nTitle: `" +  mashiron.DB_GetBucket(ModuleName,&dir,BacketInfo,[]string{"title"})[0] + "`\n\n"
				answer += result(mashiron.DB_GetFullKVList(ModuleName,&dir,BacketVotes),mashiron.DB_GetFullKVList(ModuleName,&dir,BacketChoices))
			} else if strings.HasPrefix(req.Content, conf.Prefix+ModuleName+".list") {
				answer += "Listing voters!\nTitle: `" + mashiron.DB_GetBucket(ModuleName,&dir,BacketInfo,[]string{"title"})[0] + "`\n"
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
						answer += "Invalid request!\n"
						autodelete = true
					} else {
						chk := mashiron.DB_Regex(ModuleName,BacketVotes,req.User,&dir)
						if len(chk) != 0 {
							mashiron.DB_DeleteBucket(ModuleName,&dir,BacketVotes, req.User)
						}
						mashiron.DB_AddBucket(ModuleName, &dir, BacketVotes, [][]string{{req.User, arr[1]}})
						answer += "Voted to `" + mashiron.DB_GetBucket(ModuleName,&dir,BacketInfo,[]string{"title"})[0] + "` as `" + mashiron.DB_GetBucket(ModuleName,&dir,BacketChoices,[]string{arr[1]})[0] + "`!\n"
						autodelete = true
					}
				}
			}
		} else {
			answer += "Vote not found!\n"
		}

		var opt [][]string
		secret := false
		moduleConfig := mashiron.ModuleConfig(&dir, "vote", []string{"secret"})
		if moduleConfig["secret"] != "" {
			s, err := strconv.ParseBool(moduleConfig["secret"])
			if err != nil {
				answer += "**Warning: Invalid config detected.**"
			} else {
				secret = s
			}
		}
		if req.Api == "discord" && autodelete &&  secret == true{
			opt = [][]string{{"TIMEOUT", "3"},{"DELETE", "true"}}
		}
		fmt.Print(mashiron.ResultToJSON(&mashiron.Result{
			Content: answer,
			Options: opt,
		}))
	}
}

func result(votes [][]string, choices [][]string) string {
	res := ""
	for _,v := range choices {
		a,_ := strconv.Atoi(v[0])
		res += "`" + v[0] + " (" + v[1] + ")` : `" + strconv.Itoa(calc(a,votes)) + " vote(s)`\n"
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
