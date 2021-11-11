package bd

import (
	"fmt"
	"ft-bot/logger"
	"ft-bot/store"
)

//-- Получаем очередь игроков на обновление
func GetQueuePlayers() []string {
	var queue []string
	rows, err := bd.Query("SELECT uid FROM discord_queue")
	if err != nil {
		logger.PrintLog("Queue Error: %v",err.Error())
	}
	defer rows.Close()
	var uid string
	for rows.Next() {
		if err := rows.Scan(&uid); err != nil {
			logger.PrintLog("Queue Error: %v",err.Error())
		}
		queue = append(queue, uid)
	}
	// Удаляем всех, потому что они больше не нужны
	trunc, _ := bd.Query("TRUNCATE TABLE discord_queue")
	defer trunc.Close()
	return queue
}

//-- Получить данные определенного игрока в структуре
func GetPlayer(pid string) (store.PlayerStats, error) {
	var player store.PlayerStats
	rows, err := bd.Query("select du.discord_uid, p.uid, p.playerid, p.name from players p inner join discord_users du on p.playerid = du.uid inner join player_hardwares ph on p.playerid = ph.uid where p.playerid = ?",pid)
	if err != nil {
		logger.PrintLog(err.Error())
		return player, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&player.PlayerInfo.DSUid, &player.PlayerInfo.Uid, &player.PlayerInfo.SteamId, &player.PlayerInfo.Name); err != nil {
			logger.PrintLog(err.Error())
		}
	}

	if player.PlayerInfo.DSUid == "" {
		return store.PlayerStats{}, fmt.Errorf("nothing")
	}
	rows, err = bd.Query("SELECT group_id, donorlevel from players where playerid = ?",pid)
	if err != nil {
		logger.PrintLog(err.Error())
		return player, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&player.GroupId, &player.DonatLevel); err != nil {
			logger.PrintLog(err.Error())
		}
	}
	fmt.Println(player)
	return player, nil
}

//-- Получить данные игроков - ID Group, DonatLevel
func GetStatsPlayers() []store.PlayerStats {
	allPlayers := GetAllRegisteredPlayers()
	var players []store.PlayerStats
	for _, elem := range allPlayers {
		rows, err := bd.Query("select group_id, donorlevel from players where playerid = ?", elem.SteamId)
		if err != nil {
			logger.PrintLog("GetStats Error: %v",err.Error())
		}
		var groupId, donatLevel int8
		for rows.Next() {
			if err := rows.Scan(&groupId,&donatLevel); err != nil {
				logger.PrintLog("GetStats Error: %v",err.Error())
			}
		}
		players = append(players, store.PlayerStats{PlayerInfo: elem, GroupId: groupId, DonatLevel: donatLevel})
	}
	return players
}

//-- Получить id ролей организации
func GetGroupsRole(id int8) (string,string) {
	if id == -1 {
		return "",""
	}
	rows, err := bd.Query("SELECT id, ds_role_leader, ds_role_member_id FROM groups WHERE id = ?",id)
	if err != nil {
		logger.PrintLog("GetGroup Error: %v",err.Error())
	}
	defer rows.Close()
	var groupId uint8
	var dsRoleLeader, dsRoleMember string
	for rows.Next() {
		if err := rows.Scan(&groupId,&dsRoleLeader,&dsRoleMember); err != nil {
			logger.PrintLog("GetGroup Error: %v",err.Error())
		}
		return dsRoleLeader, dsRoleMember
	}
	return "",""
}

//-- Получить id всех ролей организации
func GetAllGroupsRole() ([]string) {
	rows, err := bd.Query("SELECT ds_role_leader, ds_role_member_id FROM groups")
	if err != nil {
		logger.PrintLog("GetAllGroup Error: %v",err.Error())
	}
	defer rows.Close()
	var roles []string
	var dsRoleLeader, dsRoleMember string
	for rows.Next() {
		if err := rows.Scan(&dsRoleLeader,&dsRoleMember); err != nil {
			logger.PrintLog("GetAllGroup Error: %v",err.Error())
		}
		roles = append(roles,dsRoleLeader,dsRoleMember)
	}
	return roles
}

//-- Лидер ли он организации
func IsLeaderGroup(id int8, steamId string) bool {
	var owner, leader string
	rows, err := bd.Query("SELECT creator, leader FROM `groups` where id = ?",id)
	if err != nil {
		logger.PrintLog("IsLeader Error: %v",err.Error())
		return false
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&owner, &leader); err != nil {
			logger.PrintLog("IsLeader Error: %v", err.Error())
			return false
		}
	}
	if steamId == owner || steamId == leader {
		return true
	}
	return false
}

//-- Получить всех кто зарегистрирован
func GetAllDiscordUids() []string {
	var uids []string
	rows, err := bd.Query("select discord_uid from discord_users")
	if err != nil {
		return uids
	}
	defer rows.Close()
	var uid string
	for rows.Next() {
		if err := rows.Scan(&uid); err != nil {
			logger.PrintLog("AllDsUids Error: %v",err.Error())
		}
		uids = append(uids, uid)
	}
	return uids
}

//-- Удалить зарегистрированного
func DeleteDiscordUser(pid string) {
	rows, err := bd.Query("delete from discord_users where discord_uid = ?",pid)
	defer rows.Close()
	if err != nil {
		logger.PrintLog("DeleteUser Error: %v",err.Error())
		return
	}
}

//-- Получить всех кто зарегистрирован вольно и невольно
func GetAllRegisteredPlayers() []store.Player {
	rows, err := bd.Query("select du.discord_uid, p.uid, p.playerid, p.name from players p inner join discord_users du on p.playerid = du.uid inner join player_hardwares ph on p.playerid = ph.uid")
	var plr []store.Player
	if err != nil {
		logger.PrintLog("AllRegistered Error: %v",err.Error())
		return plr
	}
	defer rows.Close()

	var Uid uint32
	var DSUid, SteamId, Name string
	for rows.Next() {
		if err := rows.Scan(&DSUid,&Uid,&SteamId,&Name); err != nil {
			logger.PrintLog("AllRegistered Error: %v",err.Error())
		}
		plr = append(plr, store.Player{DSUid: DSUid, Uid:Uid, SteamId:SteamId, Name:Name })
	}
	return plr
}

//-- Получить данные определенного игрока
func GetPlayerStr(pid string) (string, bool) {
	type Player struct {
		Uid uint32
		SteamId string
		Name string
		DonatLevel uint8
		RC uint32
	}
	var plr Player

	rows, err := bd.Query("select p.uid, p.playerid, p.name, p.donorlevel, p.EPoint from players p inner join discord_users du on p.playerid = du.uid inner join player_hardwares ph on p.playerid = ph.uid where ph.discord_id = ? or du.discord_uid = ?", pid, pid)
	if err != nil {
		logger.PrintLog("GetPlayer Error: %v",err.Error())
		return err.Error(), false
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&plr.Uid,&plr.SteamId,&plr.Name,&plr.DonatLevel,&plr.RC); err != nil {
			logger.PrintLog("GetPlayer Error: %v",err.Error())
			return err.Error(), false
		}
	}
	if len(plr.SteamId) == 0 {
		return "Никто не найден", false
	}
	return fmt.Sprintf("Игрок: %v\nID: %d\nPID: %v\nДонат уровень: %d ур.\nRC: %d\n",plr.Name,plr.Uid,plr.SteamId,plr.DonatLevel,plr.RC), true
}
