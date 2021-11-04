
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
}

func ReadConfig() error {
	fmt.Println("Reading config file...")

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &cfg)

	if err != nil {
		fmt.Println(err.Error())
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
