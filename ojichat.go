package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/greymd/ojichat/generator"
	"github.com/joho/godotenv"
)

var (
	Token     = ""
	OjichanID = ""
	Command   = ""
	active    = false
	stopBot   = make(chan bool)
)

func main() {

	godotenv.Load()

	Token = os.Getenv("TOKEN")
	OjichanID = os.Getenv("OJICHATID")
	Command = os.Getenv("COMMAND")

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println(err)
	}

	discord.AddHandler(onMessageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("おじちゃんは聞き耳を立てています...")
	<-stopBot
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)

	if err != nil {
		return
	}

	if m.Author.ID == OjichanID {
		return
	}

	if strings.Contains(m.Content, Command) || strings.Contains(m.Content, "死ね") {
		active = !active
		if !active {
			s.ChannelMessageSend(c.ID, "じゃあね")
			return
		}
	}

	if active {
		replyOjichat(s, m.Message, c.ID)
	}
}

func replyOjichat(s *discordgo.Session, m *discordgo.Message, cid string) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}

	name := member.Nick
	if name == "" {
		name = m.Author.Username
	}

	config := generator.Config{TargetName: name}

	result, err := generator.Start(config)
	if err != nil {
		return
	}

	s.ChannelMessageSend(cid, result)
	log.Println(result)
}
