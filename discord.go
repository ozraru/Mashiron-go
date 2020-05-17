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
		Token: c.Section("token").Key("secret").String(),
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

//Start core cmd when recieved message
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
			sendMessage(session, channel, r)
	}
}

//Send message
func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg mashiron.Result) {
	//Don't send enpty message-will be denyed
	if msg.Content != ""{
		_, err := s.ChannelMessageSend(c.ID, msg.Content)
		if err != nil {
			log.Println("Error sending message: ", err)
		}
	}
	if msg.Attachments != nil {
		for _,v := range msg.Attachments {
			f, err := ioutil.ReadFile(v)
			if err != nil {
				fmt.Fprint(os.Stderr,err)
				return
			}
			s.ChannelFileSend(c.ID,v,bytes.NewBuffer(f))
		}
	}
}
