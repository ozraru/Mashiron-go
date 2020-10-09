package main

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"mashironsrv.visualstudio.com/Public/_git/Mashiron-go/mashiron"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var stopBot = make(chan bool)

type Config struct {
	Token string
}

//Read config file and start discord routine
func main() {

	log.Println("Loading config...")
	c, err := ini.Load("mashiron.ini")
	if err != nil {
		panic(err)
	}
	Cnf := Config{
		Token: c.Section("discord").Key("secret").String(),
	}

	log.Println("Connecting to discord...")
	discord, err := discordgo.New()
	if err != nil {
	panic(err)
	}
	discord.Token = Cnf.Token
	discord.AddHandler(onMessageCreate)

	err = discord.Open()
	if err != nil {
		panic(err)
	}
	log.Println("Listening...")
	<-stopBot
	return
}

//Start core cmd when received message
func onMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	channel, err := session.State.Channel(message.ChannelID) //チャンネル取得
	if err != nil {
		log.Println("Error getting channel: ", err)
		return
	}

	if message.Author.Bot {
		return
	}

	request := mashiron.Request{
		Api:        "discord",
		Room:       message.GuildID,
		User:       message.Author.ID,
		Privileges: nil,
		Content:    message.Content,
		Version:    0,
		Ishook:     false,
		Options: [][]string{
			{"CATEGORY",channel.ParentID},
			{"USERNAME",message.Author.Username},
		},
	}
	reqj := mashiron.RequestToJSON(&request)
	cmdrun := exec.Command("cmd/cmd", reqj)
	var stdOut bytes.Buffer
	cmdrun.Stdout = &stdOut
	cmdrun.Stderr = os.Stderr
	ch := make(chan bool)
	go func() {
		err = cmdrun.Run()
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
			sendMessage(session, message, r)
	}
}

//Send message
func sendMessage(session *discordgo.Session, message *discordgo.MessageCreate, result mashiron.Result) {
	timeout := 0
	if result.Options != nil {
		for _, option := range result.Options {
			switch option[0] {
			case "TIMEOUT":
				v, err := strconv.Atoi(option[1])
				if err != nil {
					log.Println("Error while converting timeout string to int: ", err)
					break
				}
				timeout = v
			case "DELETE":
				v, err := strconv.ParseBool(option[1])
				if err != nil {
					log.Println("Error while converting timeout string to int: ", err)
					break
				}
				if v == true {
					err = session.ChannelMessageDelete(message.ChannelID,message.ID)
					if err != nil {
						log.Println("Error while deleting message: ", err)
						break
					}
				}
			default:
				break
			}
		}
	}
	//Don't send empty message - Will be denied
	if result.Content != ""{
		result, err := session.ChannelMessageSend(message.ChannelID, result.Content)
		if err != nil {
			log.Println("Error while sending message: ", err)
		} else if timeout != 0 {
			ch := make(chan bool, 1)
			go func() {
				time.Sleep(time.Duration(timeout) * time.Second)
				ch <- true
			}()
			<-ch
			err = session.ChannelMessageDelete(message.ChannelID,result.ID)
			if err != nil {
				log.Println("Error while deleting message: ", err)
			}
		}
	}
	if result.Attachments != nil {
		for _,v := range result.Attachments {
			if v == "" {
				break
			}
			f, err := ioutil.ReadFile(v)
			if err != nil {
				fmt.Fprint(os.Stderr,err)
				return
			}
			session.ChannelFileSend(message.ChannelID,v,bytes.NewBuffer(f))
		}
	}
}
