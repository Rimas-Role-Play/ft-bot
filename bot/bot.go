package bot

import (
	"fmt"
	"ft-bot/bd"
	"ft-bot/config"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"regexp"
	"strings"
)

var BotID string
var goBot *discordgo.Session

var adminIds map[string]int

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

	bd.ConnectDatabase()

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Bot is running!")
}

func IsAllowedToUse(s *discordgo.Session, channelID string, roleID []string) bool {

	var ok bool
	if len(roleID) > 0 {
		_, ok = adminIds[roleID[0]]
		if ok == true {
			return ok
		}
	}
	s.ChannelMessageSend(channelID, fmt.Sprintf("Insuficient permissions"))
	return ok
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
		log.Println(vars)
		switch content {
		case "!help":
			_, _ = s.ChannelMessageSend(m.ChannelID, "Помогу, но потом")
		case "!getHim":
			if len(vars) < 1 {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Тегните пользователя, например ```!getHim @SomeUser```")
				return
			}
			re := regexp.MustCompile(`(?m)[^0-9]`)
			pid := re.ReplaceAllString(vars[0], ``)
			log.Println(pid)
			if len(pid) != 18 {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Неверный ID")
				return
			}
			_, _ = s.ChannelMessageSend(m.ChannelID, bd.GetPlayer(pid))
		case "!removeHim":
			_, _ = s.ChannelMessageSend(m.ChannelID, "Аккаунт отвязан")
		case "!applyRoles":
			_, _ = s.ChannelMessageSend(m.ChannelID, "Roles refreshed")
		}
	}
}
