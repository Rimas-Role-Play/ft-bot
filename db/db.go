package db

import (
	"database/sql"
	"fmt"
	"ft-bot/env"
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
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		env.E.MySqlUser,
		env.E.MySqlPassword,
		env.E.MySqlHost,
		env.E.MySqlPort,
		env.E.MySqlDatabase))
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
