package bd

import (
	"database/sql"
	"fmt"
	"ft-bot/config"
	"ft-bot/logger"
	"ft-bot/store"
	"log"
)

var bd *sql.DB

func ConnectDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.User,
		config.Password,
		config.IpDatabase,
		config.Port,
		config.Database))
	bd = db
	if err != nil {
		logger.PrintLog("Database not connected\n")
		return nil, err
	}

	statement, err := db.Prepare("SELECT VERSION()")
	if err != nil {
		logger.PrintLog(err.Error())
		return nil, err
	}
	rows, err := statement.Query() // execute our select statement

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for rows.Next() {
		var title string
		rows.Scan(&title)
		logger.PrintLog("Database version: %v\n", title)
	}
	return bd, nil
}

func RowsCounter(rows *sql.Rows) uint {
	var i uint
	for i = 0; rows.Next(); i++ {}
	return i
}

func GetAllNameRegisteredPlayers() []store.Player {


	rows, err := bd.Query("select du.discord_uid, p.playerid, p.name from players p inner join discord_users du on p.playerid = du.uid inner join player_hardwares ph on p.playerid = ph.uid")
	var plr []store.Player

	if err != nil {
		log.Println(err.Error())
		return plr
	}
	defer rows.Close()

	var DSUid, Uid, Name string
	for rows.Next() {
		if err := rows.Scan(&DSUid,&Uid,&Name); err != nil {
			log.Println(err.Error())
			return plr
		}
		plr = append(plr, store.Player{DSUid: DSUid, SteamId:Uid, Name:Name })
	}
	return plr
}

func GetPlayer(pid string) (string, bool) {

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
		logger.PrintLog(err.Error())
		return err.Error(), false
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&plr.Uid,&plr.SteamId,&plr.Name,&plr.DonatLevel,&plr.RC); err != nil {
			log.Println(err.Error())
			return err.Error(), false
		}
	}
	if len(plr.SteamId) == 0 {
		return "Никто не найден", false
	}
	return fmt.Sprintf("Игрок: %v\nID: %d\nPID: %v\nДонат уровень: %d ур.\nRC: %d\n",plr.Name,plr.Uid,plr.SteamId,plr.DonatLevel,plr.RC), true
}
