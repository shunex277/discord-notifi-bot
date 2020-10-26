package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type UserState struct{
	Name string
	CurrentVC string
}

var (
	userMap         = map[string]*UserState{}
	generalChanelId = os.Getenv("DISCORD_CHANNEL_ID")
	Token 			= os.Getenv("DISCORD_TOKEN")
)

func main() {
	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register ready as a callback for the ready events.
	session.AddHandler(ready)

	// Register the messageCreate func as a callback for MessageCreate events.
	//session.AddHandler(messageCreate)

	session.AddHandler(onVoiceStateUpdate)

	// In this example, we only care about receiving message events.
	//session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready)  {
	// Set the playing status.
	log.Println("BotName: ",event.User.ID)
	log.Println("BotID: ",event.User.Username)
	s.UpdateStatus(0, "bot test!")
}


// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
//func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
//
//	// Ignore all messages created by the bot itself
//	// This isn't required in this specific example but it's a good practice.
//	// ボットからのメッセージの場合は返さないように判定します。
//	if m.Author.ID == s.State.User.ID {
//		return
//	}
//
//	// Server名を取得して返します。
//	if m.Content == "ServerName" {
//		g, err := s.Guild(m.GuildID)
//		if err != nil {
//			log.Fatal(err)
//		}
//		log.Println(g.Name)
//		s.ChannelMessageSend(m.ChannelID, g.Name)
//	}
//
//	// !Helloというチャットがきたら　「Hello」　と返します
//	if m.Content == "!Hello" {
//		s.ChannelMessageSend(m.ChannelID, "Hello")
//	}
//}

func onVoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate ) {

	_, ok := userMap[vs.UserID]
	if !ok {
		//Userが居ない VC 未設定の User を追加しておく
		userMap[vs.UserID] = new(UserState)
		user, _ := s.User(vs.UserID)
		userMap[vs.UserID].Name = user.Username
		log.Print("new user added : "+user.Username)
	}

	if len(vs.ChannelID) > 0 && userMap[vs.UserID].CurrentVC != vs.ChannelID {
		channel, _ := s.Channel(vs.ChannelID)
		message := "@everyone "+ userMap[vs.UserID].Name + "さんが" + channel.Name + "で作業を開始しました"
		log.Print(message)
		s.ChannelMessageSend(generalChanelId, message)
	}

	userMap[vs.UserID].CurrentVC = vs.ChannelID

	//fmt.Printf("%+v", vs.VoiceState)
}