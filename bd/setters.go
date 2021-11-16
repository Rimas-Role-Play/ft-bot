package bd

import (
	"fmt"
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"time"
)

//-- Удалить зарегистрированного
func DeleteDiscordUser(pid string) {
	rows, err := bd.Query("delete from discord_users where discord_uid = ?",pid)
	defer rows.Close()
	if err != nil {
		logger.PrintLog("DeleteUser Error: %v",err.Error())
		return
	}
}

func InsertMessageLog(channelId string, messageId string, author *discordgo.User, mType discordgo.MessageType) {
	rows, err := bd.Query("INSERT INTO discord_messages_logs (channel_id, message_id, author_id, author_name, author_discriminator, message_type, time) VALUES (?,?,?,?,?,?,NOW())",channelId,messageId,author.ID,author.Username,author.Discriminator,mType)
	defer rows.Close()
	if err != nil {
		logger.PrintLog("InsertMessageLog Error: %v",err.Error())
		return
	}
}

var letters = []rune("ABEIKMHOPCTXZ")

func randStringRune(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randInt() int {
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999
	return rand.Intn(max - min + 1) + min
}

func generatePlateNumber() string {
	return fmt.Sprintf("DS %d %v",randInt(),randStringRune(2))
}

func InsertVehicle(pid string, classname string) {
	rows, err := bd.Query("INSERT INTO vehicles SET servermap = 'RRpMap',classname = ?, pid = ?, plate = ?," +
		"type = 'Car', alive = '1', active = '1', inventory = '[[],0]',color = 'default', material = 'default', gear = '[]', damage = '0', hitpoints = '[]', baseprice = 10000, spname = 'none', parking = '[]', maxslots = 60, tuning_data = '[[\"nitro\"],[\"tracker\"],[\"breaking\"],[\"seatbelt\"]]', distance = '0', deleted_at = NULL, comment = ''",
		pid,classname,generatePlateNumber())

	defer rows.Close()
	if err != nil {
		logger.PrintLog("InsertMessageLog Error: %v",err.Error())
		return
	}
}

