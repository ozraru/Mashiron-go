package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/ini.v1"
)

var stopBot = make(chan bool)

type Config struct {
	Token string
}

//Read config file and start discord routine
func main() {
	fmt.Println("Loading config...")
	c, err := ini.Load("mashiron.ini")
	if err != nil {
		fmt.Println("Error loading config file! Aborting...")
		return
	}
	Cnf := Config{
		Token: c.Section("token").Key("secret").String(),
	}
	fmt.Println("Connecting to discord...")
	discord, err := discordgo.New()
	discord.Token = Cnf.Token
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
		return
	}

	discord.AddHandler(onMessageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Listening...")
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
	out, err := exec.Command("cmd/cmd", "discord", message.GuildID, message.Author.ID, strings.Join(message.Member.Roles, ","), message.Content).Output()
	sendMessage(session, channel, string(out))
}

//Send message
func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	//Don't send enpty message-will be denyed
	if msg == "" {
		return
	}
	_, err := s.ChannelMessageSend(c.ID, msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}
