package bd

import (
	"database/sql"
	"errors"
	"fmt"
	"ft-bot/config"
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

func GetUidByDiscordCode(uid string) (string, error) {

	statement, err := bd.Prepare("select playerid from players where discord_code = ?")

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	defer statement.Close()

	rows, err := statement.Query(uid)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	for rows.Next() {
		var title string
		rows.Scan(&title)
		if title == "" {
			return "", errors.New("Код не определен")
		}
		return title, nil
	}
	return "", nil
}

func FireNewRegisteredUser(uid string, discordUid string, discordName string, discordEmail string) error {

	_, err := bd.Exec("insert into discord_users (uid, discord_uid, discord_name, discord_email) values (?,?,?,?)",uid,discordUid, discordName, discordEmail)
	if err != nil {
		return err
	}
	return nil
}

func CheckRegistered(user string) bool {

	statement, err := bd.Prepare(`select case when exists (select id from discord_users where discord_uid = ?) then "true" else "false" end`)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer statement.Close()

	rows, err := statement.Query(user)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for rows.Next() {
		var title string
		rows.Scan(&title)
		fmt.Println(title)
		if title == "false" {
			return true
		}
	}
	return false
}
