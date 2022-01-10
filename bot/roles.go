package bot

import (
	"ft-bot/config"
	"ft-bot/db"
	"ft-bot/logger"
	"ft-bot/store"
	"log"
)

func giveRole(pid string) {
	player, err := db.GetPlayer(pid)
	if err != nil {
		logger.PrintLog(err.Error())
		return
	}
	RoleAction(player)
}

// Giving roles
func giveRoles() {
	var users []store.PlayerStats
	users = db.GetStatsPlayers()
	for _, elem := range users {
		RoleAction(elem)
	}
	logger.PrintLog("Giving role finished")
}

func RoleAction(player store.PlayerStats) {

	_, err := s.GuildMember(config.GuildId, player.PlayerInfo.DSUid)
	if err != nil {
		log.Println(err.Error())
		log.Printf("User: %v will be deleted", player.PlayerInfo.Name)
		db.DeleteDiscordUser(player.PlayerInfo.DSUid)
		return
	}
	groupRoles := db.GetAllGroupsRole()
	if !haveRole(player.PlayerInfo.DSUid, regRoleId) {
		logger.PrintLog("Местный житель not found! %v", player.PlayerInfo.Name)
		setRole(config.GuildId, player.PlayerInfo.DSUid, regRoleId)
	}
	if player.DonatLevel > 0 {
		// Если нет роли выдаем
		if !haveRole(player.PlayerInfo.DSUid, vipRoleId) {
			setRole(config.GuildId, player.PlayerInfo.DSUid, vipRoleId)
			logger.PrintLog("Give user %v | %v VIP Role", player.PlayerInfo.Name, player.PlayerInfo.SteamId)
		}
	} else {
		// Випки нет, роль есть, удаляем
		if haveRole(player.PlayerInfo.DSUid, vipRoleId) {
			remRole(config.GuildId, player.PlayerInfo.DSUid, vipRoleId)
			logger.PrintLog("Remove user %v | %v VIP Role", player.PlayerInfo.Name, player.PlayerInfo.SteamId)
		}
	}

	// удаляем все грп роли
	for _, role := range groupRoles {
		if haveRole(player.PlayerInfo.DSUid, role) {
			logger.PrintLog("Remove role from %v", player.PlayerInfo.Name)
			remRole(config.GuildId, player.PlayerInfo.DSUid, role)
		}
	}

	// Проверяем ид грп
	if player.GroupId > 0 {
		// Получаем роли грп
		lead, member := db.GetGroupsRole(player.GroupId)
		// Если он лидер или владелец выдаем роль главы
		if db.IsLeaderGroup(player.GroupId, player.PlayerInfo.SteamId) {
			setRole(config.GuildId, player.PlayerInfo.DSUid, lead)
			logger.PrintLog("User %v added leader role FOR GroupId %d", player.PlayerInfo.Name, player.GroupId)
		} else { // Если он просто мембер выдаем роль мембера
			setRole(config.GuildId, player.PlayerInfo.DSUid, member)
			logger.PrintLog("User %v added member role FOR GroupId %d", player.PlayerInfo.Name, player.GroupId)
		}
		// Если нет грп, проходимся по ролям грп и удаляем их
	}
}

// Is have role
func haveRole(id string, roleId string) bool {
	member, err := s.GuildMember(config.GuildId, id)
	if err != nil {
		log.Println(err.Error())
		log.Printf("User: %v will be deleted", id)
		db.DeleteDiscordUser(id)
		return false
	}
	for _, role := range member.Roles {
		if role == roleId {
			return true
		}
	}
	return false
}

func setRole(guildId string, uid string, role string) bool {
	_, err := s.GuildMember(guildId, uid)
	if err != nil {
		log.Println(err.Error())
		log.Printf("User: %v will be deleted", uid)
		db.DeleteDiscordUser(uid)
		return false
	}
	s.GuildMemberRoleAdd(guildId, uid, role)
	return true
}
func remRole(guildId string, uid string, role string) bool {
	_, err := s.GuildMember(guildId, uid)
	if err != nil {
		log.Println(err.Error())
		log.Printf("User: %v will be deleted", uid)
		db.DeleteDiscordUser(uid)
		return false
	}
	s.GuildMemberRoleRemove(guildId, uid, role)
	return true
}
