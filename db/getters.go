package db

import (
	"fmt"
	"ft-bot/logger"
	"ft-bot/store"
)

//-- Получить автомобили в продаже
func GetVehiclePriceList() store.PremiumVehicles {
	var veh store.PremiumVehicles
	rows, err := db.Query("SELECT class, name, price, descr, images, sale from lk_secret_shop_price where section = 'vehicle'")
	defer rows.Close()

	if err != nil {
		logger.PrintLog("Vehicle Price Error: %v\n", err.Error())
	}
	for rows.Next() {
		if err := rows.Scan(&veh.Classname, &veh.Name, &veh.Price, &veh.Description, &veh.Images, &veh.Discount); err != nil {
			logger.PrintLog("%v\n", err.Error())
		}
	}
	return veh
}

//-- Получаем очередь игроков на обновление
func GetQueuePlayers() []string {
	var queue []string
	rows, err := db.Query("SELECT uid FROM discord_queue")
	defer rows.Close()

	if err != nil {
		logger.PrintLog("Queue Error: %v", err.Error())
		return queue
	}
	var uid string
	for rows.Next() {
		if err := rows.Scan(&uid); err != nil {
			logger.PrintLog("Queue Error: %v", err.Error())
			return queue
		}
		queue = append(queue, uid)
	}
	// Удаляем всех, потому что они больше не нужны
	trunc, _ := db.Query("TRUNCATE TABLE discord_queue")
	defer trunc.Close()
	return queue
}

//-- Получить данные определенного игрока в структуре
func GetPlayer(pid string) (store.PlayerStats, error) {
	var player store.PlayerStats
	rows, err := db.Query(`select du.discord_uid, p.uid, p.playerid, p.name, CONCAT(p.first_name," \"", p.nick_name, "\" ", p.last_name) from players p inner join discord_users du on p.playerid = du.uid inner join player_hardwares ph on p.playerid = ph.uid where p.playerid = ?`, pid)
	defer rows.Close()

	if err != nil {
		logger.PrintLog(err.Error())
		return player, err
	}
	for rows.Next() {
		if err := rows.Scan(&player.PlayerInfo.DSUid, &player.PlayerInfo.Uid, &player.PlayerInfo.SteamId, &player.PlayerInfo.Name, &player.PlayerInfo.Names); err != nil {
			logger.PrintLog(err.Error())
		}
	}

	if player.PlayerInfo.DSUid == "" {
		return store.PlayerStats{}, fmt.Errorf("nothing")
	}
	rows, err = db.Query("SELECT group_id, donorlevel from players where playerid = ?", pid)
	defer rows.Close()

	if err != nil {
		logger.PrintLog(err.Error())
		return player, err
	}

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
		rows, err := db.Query("select group_id, donorlevel from players where playerid = ?", elem.SteamId)
		defer rows.Close()
		if err != nil {
			logger.PrintLog("GetStats Error: %v", err.Error())
		}
		var groupId, donatLevel int8
		for rows.Next() {
			if err := rows.Scan(&groupId, &donatLevel); err != nil {
				logger.PrintLog("GetStats Error: %v", err.Error())
			}
		}
		players = append(players, store.PlayerStats{PlayerInfo: elem, GroupId: groupId, DonatLevel: donatLevel})
	}
	return players
}

func GetUserByDS(pid string) (store.PlayerStats, error) {
	var player store.PlayerStats
	rows, err := db.Query("select du.discord_uid, p.uid, p.playerid, p.name from players p inner join discord_users du on p.playerid = du.uid inner join player_hardwares ph on p.playerid = ph.uid where du.discord_uid = ? or ph.discord_id = ?", pid, pid)
	defer rows.Close()

	if err != nil {
		logger.PrintLog(err.Error())
		return player, err
	}
	for rows.Next() {
		if err := rows.Scan(&player.PlayerInfo.DSUid, &player.PlayerInfo.Uid, &player.PlayerInfo.SteamId, &player.PlayerInfo.Name); err != nil {
			logger.PrintLog(err.Error())
		}
	}

	if player.PlayerInfo.DSUid == "" {
		return store.PlayerStats{}, fmt.Errorf("nothing")
	}
	rows, err = db.Query("SELECT group_id, donorlevel from players where playerid = ?", player.PlayerInfo.SteamId)
	defer rows.Close()

	if err != nil {
		logger.PrintLog(err.Error())
		return player, err
	}

	for rows.Next() {
		if err := rows.Scan(&player.GroupId, &player.DonatLevel); err != nil {
			logger.PrintLog(err.Error())
		}
	}
	fmt.Println(player)
	return player, nil
}

//-- Получить id ролей организации
func GetGroupsRole(id int8) (string, string) {
	if id == -1 {
		return "", ""
	}
	rows, err := db.Query("SELECT id, ds_role_leader, ds_role_member_id FROM groups WHERE id = ?", id)
	defer rows.Close()

	if err != nil {
		logger.PrintLog("GetGroup Error: %v", err.Error())
	}
	var groupId uint8
	var dsRoleLeader, dsRoleMember string
	for rows.Next() {
		if err := rows.Scan(&groupId, &dsRoleLeader, &dsRoleMember); err != nil {
			logger.PrintLog("GetGroup Error: %v", err.Error())
		}
		return dsRoleLeader, dsRoleMember
	}
	return "", ""
}

//-- Получить id всех ролей организации
func GetAllGroupsRole() []string {
	rows, err := db.Query("SELECT ds_role_leader, ds_role_member_id FROM groups")
	defer rows.Close()

	if err != nil {
		logger.PrintLog("GetAllGroup Error: %v", err.Error())
	}
	var roles []string
	var dsRoleLeader, dsRoleMember string
	for rows.Next() {
		if err := rows.Scan(&dsRoleLeader, &dsRoleMember); err != nil {
			logger.PrintLog("GetAllGroup Error: %v", err.Error())
		}
		roles = append(roles, dsRoleLeader, dsRoleMember)
	}
	return roles
}

//-- Лидер ли он организации
func IsLeaderGroup(id int8, steamId string) bool {
	var owner, leader string
	rows, err := db.Query("SELECT creator, leader FROM `groups` where id = ?", id)
	defer rows.Close()

	if err != nil {
		logger.PrintLog("IsLeader Error: %v", err.Error())
		return false
	}
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
	rows, err := db.Query("select discord_uid from discord_users")
	defer rows.Close()

	if err != nil {
		return uids
	}
	var uid string
	for rows.Next() {
		if err := rows.Scan(&uid); err != nil {
			logger.PrintLog("AllDsUids Error: %v", err.Error())
		}
		uids = append(uids, uid)
	}
	return uids
}

//-- Получить всех кто зарегистрирован вольно и невольно
func GetAllRegisteredPlayers() []store.Player {
	rows, err := db.Query("select du.discord_uid, p.uid, p.playerid, p.name from players p inner join discord_users du on p.playerid = du.uid inner join player_hardwares ph on p.playerid = ph.uid")
	defer rows.Close()

	var plr []store.Player
	if err != nil {
		logger.PrintLog("AllRegistered Error: %v", err.Error())
		return plr
	}

	var Uid uint32
	var DSUid, SteamId, Name string
	for rows.Next() {
		if err := rows.Scan(&DSUid, &Uid, &SteamId, &Name); err != nil {
			logger.PrintLog("AllRegistered Error: %v", err.Error())
		}
		plr = append(plr, store.Player{DSUid: DSUid, Uid: Uid, SteamId: SteamId, Name: Name})
	}
	return plr
}

//-- Получить данные определенного игрока
func GetPlayerStr(pid string) (string, bool) {
	type Player struct {
		Uid        uint32
		SteamId    string
		Name       string
		Names      string
		DonatLevel uint8
		RC         uint32
	}
	var plr Player

	rows, err := db.Query(`select p.uid, p.playerid, p.name, CONCAT(p.first_name," \"", p.nick_name, "\" ", p.last_name), p.donorlevel, p.EPoint `+
		"from players p "+
		"left join discord_users du on p.playerid = du.uid "+
		"right join player_hardwares ph on p.playerid = ph.uid "+
		"where ph.discord_id = ? or du.discord_uid = ?", pid, pid)
	defer rows.Close()
	if err != nil {
		logger.PrintLog("GetPlayer Error: %v", err.Error())
		return err.Error(), false
	}

	for rows.Next() {
		if err := rows.Scan(&plr.Uid, &plr.SteamId, &plr.Name, &plr.Names, &plr.DonatLevel, &plr.RC); err != nil {
			logger.PrintLog("GetPlayer Error: %v", err.Error())
			return err.Error(), false
		}
	}
	if len(plr.SteamId) == 0 {
		return "Никто не найден", false
	}
	return fmt.Sprintf("Имя профиля: %v\nИмя: %s\nID: %d\nPID: %v\nДонат уровень: %d ур.\nRC: %d\n", plr.Name, plr.Names, plr.Uid, plr.SteamId, plr.DonatLevel, plr.RC), true
}

func GetRandomVehicle() *store.Vehicles {
	var veh store.Vehicles
	rows, err := db.Query("SELECT d.classname, d.image, m.displayName FROM discord_boosters d " +
		"INNER JOIN lk_mapobjects m ON m.classname = d.classname " +
		"WHERE active = 1 " +
		"ORDER BY RAND() " +
		"LIMIT 1")
	defer rows.Close()

	if err != nil {
		logger.PrintLog("GetRandomVehicle Error: %s\n", err.Error())
		return nil
	}
	for rows.Next() {
		if err := rows.Scan(&veh.Classname, &veh.Image, &veh.DisplayName); err != nil {
			logger.PrintLog("GetRandomVehicle Error: %s\n", err.Error())
			return nil
		}
	}
	return &veh
}
