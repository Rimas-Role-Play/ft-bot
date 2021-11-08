package bot

import (
	"flag"
	"fmt"
	"ft-bot/config"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
	"strings"
)

var (
	GuildID        = flag.String("guild", config.GuildId, "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", config.Token, "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

var BotID string

func Start() {
	var err error
	s, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := s.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID
	s.AddHandler(messageHandler)
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
			log.Printf("Command %v called\n", i.ApplicationCommandData().Name)
		}
	})

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
		return
	}
	defer s.Close()

	for _, v := range commands {

		_, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		log.Printf("Command %v created", v.Name)
	}

	log.Println("Start goroutine")
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdown")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		if m.Author.ID == BotID {
			return
		}

		if !IsDiscordAdmin(s, m.Author.ID ) {
			log.Printf("%v попытался использовать!\n", m.Author)
			return
		}

		var vars []string
		var content string
		inputSplit := strings.Split(m.Content, " ")
		for idx := range inputSplit {
			if idx == 0 {
				content = inputSplit[idx]
			}else{
				vars = append(vars, inputSplit[idx])
			}
		}

		switch content {
		case "!help":
			SendMessage(s, "Помогу, но потом")
			RenameUsers()
		}
	}
}