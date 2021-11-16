
package config

import (
	"encoding/json"
	"fmt"
	"ft-bot/logger"
	"io/ioutil"
	"math/rand"
	"time"
)

var (
	Token     	string
	BotPrefix 	string

	IpDatabase 	string
	Port 		string
	Database 	string
	User 		string
	Password 	string
	GuildId 	string
	cfg *configStruct
)

type configStruct struct {
	Token     	string `json:"Token"`
	BotPrefix 	string `json:"BotPrefix"`

	IpDatabase 	string `json:"IpDatabase"`
	Port 		string `json:"Port"`
	Database 	string `json:"Database"`
	User 		string `json:"User"`
	Password 	string `json:"Password"`
	GuildId 	string `json:"GuildId"`
	AdminRoles 	[]string `json:"adminRoles"`
	VehiclesForBoost 	[]string `json:"vehiclesForBoost"`
}

func GetAdminRoles() []string {
	return cfg.AdminRoles
}

func GetRandomVehicle() string {
	rand.Seed(time.Now().Unix())
	random := rand.Intn(len(cfg.VehiclesForBoost))
	return cfg.VehiclesForBoost[random]
}

func ReadConfig() error {
	fmt.Println("Reading config file...")

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		logger.PrintLog(err.Error())
		return err
	}

	err = json.Unmarshal(file, &cfg)

	if err != nil {
		logger.PrintLog(err.Error())
		return err
	}

	Token = cfg.Token
	BotPrefix = cfg.BotPrefix

	IpDatabase 	= cfg.IpDatabase
	Port 		= cfg.Port
	Database 	= cfg.Database
	User 		= cfg.User
	Password 	= cfg.Password
	GuildId 	= cfg.GuildId
	return nil
}
