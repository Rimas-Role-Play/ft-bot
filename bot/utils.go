package bot

import (
	"ft-bot/bd"
	"ft-bot/config"
	"ft-bot/logger"
	"ft-bot/store"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
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
func DeleteUndefinedUsers() {
	uids := bd.GetAllDiscordUids()
	for _, elem := range uids {
		_, err := s.GuildMember(config.GuildId,elem)
		if err != nil {
			log.Println(err.Error())
			log.Printf("User: %v will be deleted",elem)
			bd.DeleteDiscordUser(elem)
		}
	}
	log.Printf("All inactive users deleted")
}
func RenameUsers() {
	ticker := time.NewTicker(4 * 60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				logger.PrintLog("Starting rename")
				players := bd.GetAllRegisteredPlayers()
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
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
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

func haveRole(id string, roleId string) bool {
	member, _ := s.GuildMember(config.GuildId,id)
	for _, role := range member.Roles {
		if role == roleId {
			return true
		}
	}
	return false
}
var vipRoleId = "871508004795723787"
func GiveRoles() {
	var users []store.PlayerStats
	users = bd.GetStatsPlayers()
	for _, elem := range users {
		// Есть ли у него випка
		if elem.DonatLevel > 0 {
			// Если нет роли выдаем
			if !haveRole(elem.PlayerInfo.DSUid, vipRoleId) {
				s.GuildMemberRoleAdd(config.GuildId,elem.PlayerInfo.DSUid,vipRoleId)
				logger.PrintLog("Give user %v | %v VIP Role", elem.PlayerInfo.Name, elem.PlayerInfo.SteamId)
			}
		}else{
			// Випки нет, роль есть, удаляем
			if haveRole(elem.PlayerInfo.DSUid, vipRoleId) {
				s.GuildMemberRoleRemove(config.GuildId,elem.PlayerInfo.DSUid,vipRoleId)
				logger.PrintLog("Remove user %v | %v VIP Role", elem.PlayerInfo.Name, elem.PlayerInfo.SteamId)
			}
		}
		// Проверяем ид грп
		if elem.GroupId > 0 {
			// Получаем роли грп
			lead, member := bd.GetGroupsRole(elem.GroupId)
			// Если он лидер или владелец выдаем роль главы
			if bd.IsLeaderGroup(elem.GroupId,elem.PlayerInfo.SteamId) {
				s.GuildMemberRoleAdd(config.GuildId,elem.PlayerInfo.DSUid,lead)
				logger.PrintLog("User %v added leader role FOR GroupId %d",elem.PlayerInfo.SteamId, elem.GroupId)
			}else{ // Если он просто мембер выдаем роль мембера
				s.GuildMemberRoleAdd(config.GuildId,elem.PlayerInfo.DSUid,member)
				logger.PrintLog("User %v added member role FOR GroupId %d",elem.PlayerInfo.SteamId, elem.GroupId)
			}
		// Если нет грп, проходимся по ролям грп и удаляем их
		}else{
			for _, role := range bd.GetAllGroupsRole() {
				if haveRole(elem.PlayerInfo.DSUid, role) {
					s.GuildMemberRoleRemove(config.GuildId, elem.PlayerInfo.DSUid, role)
				}
			}
		}
	}
}