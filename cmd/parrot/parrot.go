package main

import (
	"fmt"
	"mashironsrv.visualstudio.com/Public/_git/Mashiron-go/mashiron"
	"os"
	"strings"
)

func main() {
	r := mashiron.JSONToRequest(&os.Args[1])
	d := mashiron.GetDirList(&r,"parrot")
	c := mashiron.GetCoreConf(&d)
	if r.Ishook {
		return
	} else {
		if strings.HasPrefix(r.Content,c.Prefix+"parrot") {
			s := strings.Split(r.Content," ")
			switch len(s) {
			case 1:
				fmt.Print(mashiron.ResultToJSON(&mashiron.Result{
					Attachments: []string{d.CmdDir+"parrots/parrot.gif"},
					Content:     "[HINT] : "+c.Prefix + "parrot <parrotname>\n" + "https://cultofthepartyparrot.com/",
				}))
				return
			case 2:
				if strings.Contains(s[1],"..") {
					return
				}
				if Exists(d.CmdDir+"parrots/"+s[1]+"parrot.gif") {
					fmt.Print(mashiron.ResultToJSON(&mashiron.Result{
						Attachments: []string{d.CmdDir+"parrots/"+s[1]+"parrot.gif"},
						Content:     "",
					}))
				} else if Exists(d.CmdDir+"parrots/"+s[1]+".gif") {
					fmt.Print(mashiron.ResultToJSON(&mashiron.Result{
						Attachments: []string{d.CmdDir+"parrots/"+s[1]+".gif"},
						Content:     "",
					}))
				}
				return
			default:
				return
			}
		}
	}
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}