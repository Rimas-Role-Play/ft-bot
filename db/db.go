package db

import (
	"database/sql"
	"fmt"
	"ft-bot/config"
	"ft-bot/logger"
)

// Global variable in db package
var db *sql.DB

func init() {
	var err error
	db, err = connectDatabase()
	if err != nil {
		logger.PrintLog("cant open database %s\n", err.Error())
	}
	logger.PrintLog("Database connected")
}

// Connect to database
func connectDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.User,
		config.Password,
		config.IpDatabase,
		config.Port,
		config.Database))
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
	return db, nil
}
