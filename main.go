package main

import (
	"fmt"
	"ft-bot/bd"
	"ft-bot/bot"
	"ft-bot/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bd.ConnectDatabase()
	fmt.Printf(` 
	________  .__                               .___
	\______ \ |__| ______ ____   ___________  __| _/
	||    |  \|  |/  ___// ___\ /  _ \_  __ \/ __ | 
	||    '   \  |\___ \/ /_/  >  <_> )  | \/ /_/ | 
	||______  /__/____  >___  / \____/|__|  \____ | 
	\_______\/        \/_____/   %-16s\/`+"\n\n", "1.1.0")

	bot.Start()
	return
}