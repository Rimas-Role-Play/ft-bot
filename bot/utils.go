package bot

import (
	"fmt"
	"ft-bot/bd"
	"ft-bot/config"
	"ft-bot/logger"
	"ft-bot/store"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

var vipRoleId = "871508004795723787"
var regRoleId = "864630308242849825"

// Listener of discord_queue
func ListenQueue() {
	logger.PrintLog("Check queue")
	queue := bd.GetQueuePlayers()
	for _, elem := range queue {
		logger.PrintLog("%v in queue right now", elem)
		RenameUser(elem)
		GiveRole(elem)
	}
	logger.PrintLog("Queue finished")
}

func RenameUser(pid string) {
	var err error
	player, err := bd.GetPlayer(pid)
	if err != nil {
		logger.PrintLog(err.Error())
		return
	}
	err = s.GuildMemberNickname(config.GuildId, player.PlayerInfo.DSUid, player.PlayerInfo.Name)
	if err != nil {
		logger.PrintLog("Cant rename %v user",player.PlayerInfo.Name)
		logger.PrintLog(err.Error())
	}
}

func GiveRole(pid string) {
	player, err := bd.GetPlayer(pid)
	if err != nil {
		logger.PrintLog(err.Error())
		return
	}
	RoleAction(player)
}

func RoleAction(player store.PlayerStats) {
	groupRoles := bd.GetAllGroupsRole()
	if !haveRole(player.PlayerInfo.DSUid, regRoleId) {
		logger.PrintLog("Местный житель not found! %v",player.PlayerInfo.Name)
		s.GuildMemberRoleAdd(config.GuildId, player.PlayerInfo.DSUid,regRoleId)
	}
	if player.DonatLevel > 0 {
		// Если нет роли выдаем
		if !haveRole(player.PlayerInfo.DSUid, vipRoleId) {
			s.GuildMemberRoleAdd(config.GuildId,player.PlayerInfo.DSUid,vipRoleId)
			logger.PrintLog("Give user %v | %v VIP Role", player.PlayerInfo.Name, player.PlayerInfo.SteamId)
		}
	}else{
		// Випки нет, роль есть, удаляем
		if haveRole(player.PlayerInfo.DSUid, vipRoleId) {
			s.GuildMemberRoleRemove(config.GuildId,player.PlayerInfo.DSUid,vipRoleId)
			logger.PrintLog("Remove user %v | %v VIP Role", player.PlayerInfo.Name, player.PlayerInfo.SteamId)
		}
	}

	// удаляем все грп роли
	for _, role := range groupRoles {
		if haveRole(player.PlayerInfo.DSUid, role) {
			logger.PrintLog("Remove role from %v",player.PlayerInfo.Name)
			s.GuildMemberRoleRemove(config.GuildId, player.PlayerInfo.DSUid, role)
		}
	}

	// Проверяем ид грп
	if player.GroupId > 0 {
		// Получаем роли грп
		lead, member := bd.GetGroupsRole(player.GroupId)
		// Если он лидер или владелец выдаем роль главы
		if bd.IsLeaderGroup(player.GroupId,player.PlayerInfo.SteamId) {
			s.GuildMemberRoleAdd(config.GuildId,player.PlayerInfo.DSUid,lead)
			logger.PrintLog("User %v added leader role FOR GroupId %d",player.PlayerInfo.Name, player.GroupId)
		}else{ // Если он просто мембер выдаем роль мембера
			s.GuildMemberRoleAdd(config.GuildId,player.PlayerInfo.DSUid,member)
			logger.PrintLog("User %v added member role FOR GroupId %d",player.PlayerInfo.Name, player.GroupId)
		}
		// Если нет грп, проходимся по ролям грп и удаляем их
	}
}

// Check is discord admin
func IsDiscordAdmin(s *discordgo.Session, id string) bool {
	g, _ := s.GuildMember("719969719871995958", id)
	roles := g.Roles
	for _, i := range config.GetAdminRoles() {
		for _, y := range roles {
			if i == y {
				log.Println("Admin found")
				return true
			}
		}
	}
	log.Println("Admin not found")
	return false
}

// Delete leaved or banned users
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

// Start ticker routines
func StartRoutine() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				ListenQueue()
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

// Send information from command get-him
func GetHim(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].UserValue(nil)
	pid := user.ID
	sender := i.Interaction.Member.User
	logger.PrintLog("User %v#%v want get information about %v",sender.Username, sender.Discriminator, pid)
	if !IsDiscordAdmin(s, i.Interaction.Member.User.ID) {
		logger.PrintLog("User %v#%v dont have access to /gethim",sender.Username, sender.Discriminator)
		PrintHiddenMessage(s,i,"У вас нет доступа")
		return
	}
	data, found := bd.GetPlayerStr(pid)
	if found {
		logger.PrintLog("User %v found!",pid)
	}else{
		logger.PrintLog(data)
	}
	PrintHiddenMessage(s,i,data)
}

func PrintHiddenMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: msg,
		},
	})
}

// Is have role
func haveRole(id string, roleId string) bool {
	member, _ := s.GuildMember(config.GuildId,id)
	for _, role := range member.Roles {
		if role == roleId {
			return true
		}
	}
	return false
}

// Giving roles
func GiveRoles() {
	var users []store.PlayerStats
	users = bd.GetStatsPlayers()
	for _, elem := range users {
		RoleAction(elem)
	}
	logger.PrintLog("Giving role finished")
}

func giveBoostPresent(channelId string, user *discordgo.User) {
	player, err := bd.GetUserByDS(user.ID)
	if err != nil {
		logger.PrintLog("Cant give vehicle for boost: %v", err.Error())
	}
	vehicle := config.GetRandomVehicle()
	bd.InsertVehicle(vehicle.Classname,player.PlayerInfo.SteamId)
	s.ChannelMessageSend(channelId,pingUser(user.ID))
	s.ChannelMessageSendEmbed(channelId, createEmbedNitroBooster(vehicle))
	logger.PrintLog("%v boosted server and given %v",user.Username,vehicle.DisplayName)
}

func createEmbedNitroBooster(vehicle config.Vehicles) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		URL: "",
		Type: discordgo.EmbedTypeImage,
		Title: "Nitro Booster",
		Description: fmt.Sprintf("Спасибо за буст сервера!\nТвой подарок %v уже доступен на сервере!",vehicle.DisplayName),
		Timestamp: "",
		Color: 0x9300FF,
		Footer: &discordgo.MessageEmbedFooter{
			Text:         "Nitro Boost",
			IconURL:      "",
			ProxyIconURL: "",
		},
		Image: &discordgo.MessageEmbedImage{
			URL:      vehicle.Image,
			ProxyURL: "",
			Width:    0,
			Height:   0,
		},
	}
	return embed
}

func pingUser(id string) string {
	return fmt.Sprintf("<@%v>",id)
}
