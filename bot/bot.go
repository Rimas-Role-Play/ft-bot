package bot

import (
	"flag"
	"fmt"
	"ft-bot/bd"
	"ft-bot/config"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

var (
	GuildID        = flag.String("guild", config.GuildId, "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", config.Token, "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

var BotID string
var goBot *discordgo.Session

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")
}

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	bd.ConnectDatabase()
	fmt.Println("Bot is running!")
}

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID
	goBot.AddHandler(messageHandler)
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func IsDiscordAdmin(s *discordgo.Session, id string) bool {
	g, _ := s.GuildMember("719969719871995958", id)
	roles := g.Roles
	for idx, _ := range roles {
		if roles[idx] == "775499720222310411" || roles[idx] == "878824238075748372" {
			return true
		}
	}
	return false
}

func SendMessage(s *discordgo.Session, msg string) {
	s.ChannelMessageSend("864640641891696641", msg)
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
		case "!getHim":
			if len(vars) == 0 || len(vars) > 1 {
				SendMessage(s, "Тегните пользователя, например ```!getHim @SomeUser```")
				return
			}
			re := regexp.MustCompile(`(?m)[^0-9]`)
			pid := re.ReplaceAllString(vars[0], ``)
			if len(pid) != 18 {
				SendMessage(s, "Неверный ID")
				return
			}
			SendMessage(s, bd.GetPlayer(pid))
		case "!removeHim":
			SendMessage(s, "Аккаунт отвязан")
		case "!applyRoles":
			SendMessage(s, "Roles refreshed")
		}
		s.ChannelMessageDelete(m.Message.ChannelID, m.Message.ID)
	}
}
