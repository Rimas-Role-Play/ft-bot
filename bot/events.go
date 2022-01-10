package bot

import (
	"ft-bot/config"
	"ft-bot/db"
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// UserConnect event handler
func OnUserConnected(s *discordgo.Session, u *discordgo.GuildMemberAdd) {
	user := u.Member.User
	logger.PrintLog("New user connected %v#%v | ID: %v", user.Username, user.Discriminator, user.ID)
}

// UserDisconnected event handler
func OnUserDisconnected(s *discordgo.Session, u *discordgo.GuildMemberRemove) {
	user := u.Member.User
	logger.PrintLog("User disconnected %v#%v | ID: %v", user.Username, user.Discriminator, user.ID)
	db.DeleteDiscordUser(user.ID)
}

// UserBoosted event handler
func OnUserChanged(s *discordgo.Session, i *discordgo.GuildMemberUpdate) {}

// Messages event handler
func OnMessageHandle(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch m.Message.Type {
	case 7:
		db.InsertMessageLog(m.ChannelID, m.ID, m.Author, m.Message.Type)

	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierOne:
		fallthrough
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierTwo:
		fallthrough
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierThree:
		fallthrough
	case discordgo.MessageTypeUserPremiumGuildSubscription:
		db.InsertMessageLog(m.ChannelID, m.ID, m.Author, m.Message.Type)
		giveBoostPresent(m.ChannelID, m.Author)
	case 0:
		fallthrough
	case 19:
		if isMuted(s, m.Author.ID) {
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				logger.PrintLog(err.Error())
			}
			return
		}
	}

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !isDiscordAdmin(s, m.Author.ID) {
			logger.PrintLog("%v попытался использовать!\n", m.Author)
			return
		}

		var vars []string
		var content string
		inputSplit := strings.Split(m.Content, " ")
		for idx := range inputSplit {
			if idx == 0 {
				content = inputSplit[idx]
			} else {
				vars = append(vars, inputSplit[idx])
			}
		}
		switch content {
		case "!help":
		case "!test":
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
