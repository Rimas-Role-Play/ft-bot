package bd

import (
	"database/sql"
	"fmt"
	"ft-bot/config"
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
		fmt.Println("Database not connected")
		return nil, err
	}

	statement, err := db.Prepare("SELECT VERSION()")
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println("Database version: ", title)
	}
	return bd, nil
}



func GetPlayer(pid string) string {

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
		log.Println(err.Error())
		return err.Error()
	}
	log.Println(rows.Err())
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&plr.Uid,&plr.SteamId,&plr.Name,&plr.DonatLevel,&plr.RC); err != nil {
			log.Println(err.Error())
			return err.Error()
		}
	}
	if len(plr.SteamId) == 0 {
		return "Никто не найден"
	}
	return fmt.Sprintf("Игрок: %v\nID: %d\nPID: %v\nДонат уровень: %d ур.\nRC: %d\n",plr.Name,plr.Uid,plr.SteamId,plr.DonatLevel,plr.RC)
}
