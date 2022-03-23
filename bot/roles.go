package bot

import (
	"fmt"
	"ft-bot/db"
	"ft-bot/env"
	"ft-bot/logger"
	"ft-bot/store"
	"github.com/bwmarrin/discordgo"
	"log"
)

func RoleAction(player store.PlayerStats) {

	_, err := s.GuildMember(env.E.GuildId, player.PlayerInfo.DSUid)
	if err != nil {
		log.Println(err.Error())
		log.Printf("User: %v will be deleted", player.PlayerInfo.Name)
		db.DeleteDiscordUser(player.PlayerInfo.DSUid)
		return
	}
	groupRoles := db.GetAllGroupsRole()
	if !haveRole(player.PlayerInfo.DSUid, regRoleId) {
		logger.PrintLog("Местный житель not found! %v", player.PlayerInfo.Name)
		setRole(env.E.GuildId, player.PlayerInfo.DSUid, regRoleId)
	}
	if player.DonatLevel > 0 {
		// Если нет роли выдаем
		if !haveRole(player.PlayerInfo.DSUid, vipRoleId) {
			setRole(env.E.GuildId, player.PlayerInfo.DSUid, vipRoleId)
			logger.PrintLog("Give user %v | %v VIP Role", player.PlayerInfo.Name, player.PlayerInfo.SteamId)
		}
	} else {
		// Випки нет, роль есть, удаляем
		if haveRole(player.PlayerInfo.DSUid, vipRoleId) {
			remRole(env.E.GuildId, player.PlayerInfo.DSUid, vipRoleId)
			logger.PrintLog("Remove user %v | %v VIP Role", player.PlayerInfo.Name, player.PlayerInfo.SteamId)
		}
	}

	// удаляем все грп роли
	for _, role := range groupRoles {
		if haveRole(player.PlayerInfo.DSUid, role) {
			logger.PrintLog("Remove role from %v", player.PlayerInfo.Name)
			remRole(env.E.GuildId, player.PlayerInfo.DSUid, role)
		}
	}

	// Проверяем ид грп
	if player.GroupId > 0 {
		// Получаем роли грп
		lead, member := db.GetGroupsRole(player.GroupId)
		// Если он лидер или владелец выдаем роль главы
		if db.IsLeaderGroup(player.GroupId, player.PlayerInfo.SteamId) {
			setRole(env.E.GuildId, player.PlayerInfo.DSUid, lead)
			logger.PrintLog("User %v added leader role FOR GroupId %d", player.PlayerInfo.Name, player.GroupId)
		} else { // Если он просто мембер выдаем роль мембера
			setRole(env.E.GuildId, player.PlayerInfo.DSUid, member)
			logger.PrintLog("User %v added member role FOR GroupId %d", player.PlayerInfo.Name, player.GroupId)
		}
		// Если нет грп, проходимся по ролям грп и удаляем их
	}
}

//-- Unexported functions

func findRoleById(guildId, roleId string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildId)
	if err != nil {
		return &discordgo.Role{}, err
	}
	for _, elem := range roles {
		if elem.ID == roleId {
			return elem, nil
		}
	}
	return &discordgo.Role{}, fmt.Errorf("role not found")
}

func copyRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	roleId := i.ApplicationCommandData().Options[0].RoleValue(s, "").ID
	role, err := findRoleById(env.E.GuildId, roleId)
	if err != nil {
		printHiddenMessage(s, i, "ошибка при создании роли "+err.Error())
		return
	}
	nameNewRole := i.ApplicationCommandData().Options[1].StringValue()
	newRole, err := s.GuildRoleCreate(env.E.GuildId)
	if err != nil {
		printHiddenMessage(s, i, "ошибка при создании роли "+err.Error())
		return
	}
	fmt.Println(env.E.GuildId, role.ID, role.Name, role.Color, role.Hoist, role.Permissions, role.Mentionable)
	fmt.Println(env.E.GuildId, newRole.ID, nameNewRole, role.Color, role.Hoist, role.Permissions, role.Mentionable)
	newRole, err = s.GuildRoleEdit(env.E.GuildId, newRole.ID, nameNewRole, role.Color, role.Hoist, role.Permissions, role.Mentionable)
	if err != nil {
		printHiddenMessage(s, i, "ошибка при изменении роли "+err.Error())
		return
	}
	printHiddenMessage(s, i, "Роль "+role.Mention()+" скопирована в "+newRole.Mention())
}

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

// Is have role
func haveRole(id string, roleId string) bool {
	member, err := s.GuildMember(env.E.GuildId, id)
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
