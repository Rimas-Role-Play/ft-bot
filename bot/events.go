package bot

import (
	"fmt"
	"ft-bot/config"
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// UserConnect event handler
func OnUserConnected(s *discordgo.Session, u *discordgo.GuildMemberAdd) {
	user := u.Member.User
	logger.PrintLog("New user connected %v#%v | ID: %v",user.Username, user.Discriminator, user.ID )
}

// UserBoosted event handler
func OnUserBoosted(s *discordgo.Session, i *discordgo.GuildMemberUpdate) {
	fmt.Printf("Fired\n")
	fmt.Println(i.Member.User.Username)
}

// Messages event handler
func OnMessageHandle(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !IsDiscordAdmin(s, m.Author.ID ) {
			logger.PrintLog("%v попытался использовать!\n", m.Author)
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
		case "!test":
			go GiveRoles()
		}
	}
}

// Command trigger event handler
func OnCommandsCall(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
		logger.PrintLog("Command %v called\n", i.ApplicationCommandData().Name)
	}
}

// Reaction trigger
func OnReactMessage(s *discordgo.Session, i *discordgo.MessageReactionAdd) {
	log.Println(i)
}
