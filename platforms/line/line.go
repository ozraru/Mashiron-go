//https://github.com/line/line-bot-sdk-go/blob/master/examples/echo_bot/server.go.

package main

import (
	"bytes"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"mashironsrv.visualstudio.com/Public/_git/Mashiron-go/mashiron"
	"net/http"
	"os"
	"os/exec"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Config struct {
	Token string
	Secret string
	Port string
}

func main() {

	log.Println("Loading config...")
	c, err := ini.Load("mashiron.ini")
	if err != nil {
		panic(err)
	}
	Cnf := Config{
		Token: c.Section("line").Key("token").String(),
		Secret: c.Section("line").Key("secret").String(),
		Port: c.Section("line").Key("port").String(),
	}
	bot, err := linebot.New(Cnf.Secret,Cnf.Token)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(start(event,message.Text))).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+Cnf.Port, nil); err != nil {
		log.Fatal(err)
	}
}

func start(event *linebot.Event, text string) string {
	request := mashiron.Request{
		Api:        "discord",
		Room:       event.Source.RoomID,
		User:       event.Source.UserID,
		Privileges: nil,
		Content:    text,
		Version:    0,
		Ishook:     false,
		Options: [][]string{

		},
	}
	reqj := mashiron.RequestToJSON(&request)
	cmdrun := exec.Command("cmd/cmd", reqj)
	var stdOut bytes.Buffer
	cmdrun.Stdout = &stdOut
	cmdrun.Stderr = os.Stderr
	ch := make(chan bool)
	go func() {
		err := cmdrun.Run()
		if err != nil {
			fmt.Print(err.Error())
		}
		ch <- true
	}()
	<-ch
	o := stdOut.String()
	//fmt.Println(reqj + "\n=>\n" +o) //Just uncomment this for debugging!
	resultArr := mashiron.JSONToResults(&o)
	log.Println(resultArr)
	for _,v := range resultArr {
		r := mashiron.JSONToResult(&v)
		if r.Content != ""{
			return r.Content
		}
	}
	return ""
}