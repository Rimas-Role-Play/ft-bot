package bd

import (
	"database/sql"
	"fmt"
	"ft-bot/config"
	"ft-bot/logger"
)

// Global variable in bd package
var bd *sql.DB

// Connect to database
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
