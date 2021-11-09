package bot

import (
	"ft-bot/bd"
	"ft-bot/config"
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
	"log"
)

func IsDiscordAdmin(s *discordgo.Session, id string) bool {
	g, _ := s.GuildMember("719969719871995958", id)
	log.Println(g.User.Username)
	roles := g.Roles
	for idx, _ := range roles {
		if roles[idx] == "775499720222310411" || roles[idx] == "878824238075748372" || roles[idx] == "866252450234630204" {
			log.Println("Admin found")
			return true
		}
	}
	log.Println("Admin not found")
	return false
}

func RenameUsers() {
	logger.PrintLog("Starting rename")
	players := bd.GetAllNameRegisteredPlayers()
	for idx, _ := range players {

		err := s.GuildMemberNickname(config.GuildId, players[idx].DSUid, players[idx].Name)
		if err != nil {
			logger.PrintLog("****************************************************************\n")
			logger.PrintLog(err.Error())
			logger.PrintLog("User ID: %v | Uid: %v | Name: %v\n",players[idx].DSUid,players[idx].Uid, players[idx].Name)
			logger.PrintLog("****************************************************************\n")
		}else{
			logger.PrintLog("User: %v, renamed to: %v | IDX: %d/%d", players[idx].DSUid, players[idx].Name, idx, len(players) - 1)
		}
	}
	logger.PrintLog("Renaming is done")
}

func SendMessage(s *discordgo.Session, msg string) {
	s.ChannelMessageSend("864640641891696641", msg)
}

func GetHim(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].UserValue(nil)
	pid := user.ID
	sender := i.Interaction.Member.User
	logger.PrintLog("User %v#%v want get information about %v",sender.Username, sender.Discriminator, pid)
	if !IsDiscordAdmin(s, i.Interaction.Member.User.ID) {
		logger.PrintLog("User %v#%v dont have access to /gethim",sender.Username, sender.Discriminator)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: "У вас нет доступа",
			},
		})
		return
	}
	data, found := bd.GetPlayer(pid)
	if found {
		logger.PrintLog("User %v found!",pid)
	}else{
		logger.PrintLog(data)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:   1 << 6,
			Content: data,
		},
	})
}
