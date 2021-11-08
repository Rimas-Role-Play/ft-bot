package bot

import (
	"ft-bot/bd"
	"ft-bot/config"
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
	players := bd.GetAllNameRegisteredPlayers()
	for idx, _ := range players {

		err := s.GuildMemberNickname(config.GuildId, players[idx].DSUid, players[idx].Name)
		if err != nil {
			log.Printf("****************************************************************")
			log.Println(err.Error())
			log.Printf("User ID: %v | Uid: %v | Name: %v\n",players[idx].DSUid,players[idx].Uid, players[idx].Name)
			log.Printf("****************************************************************")
		}else{
			log.Printf("User: %v, renamed to: %v | IDX: %d/%d", players[idx].DSUid, players[idx].Name, idx, len(players) - 1)
		}
	}
	log.Printf("Renaming is done")
}

func SendMessage(s *discordgo.Session, msg string) {
	s.ChannelMessageSend("864640641891696641", msg)
}

func GetHim(s *discordgo.Session, i *discordgo.InteractionCreate) {
	pid := i.ApplicationCommandData().Options[0].UserValue(nil).ID
	log.Println(i.Interaction.Member.User.ID)
	if !IsDiscordAdmin(s, i.Interaction.Member.User.ID) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: "У вас нет доступа",
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:   1 << 6,
			Content: bd.GetPlayer(pid),
		},
	})
}
