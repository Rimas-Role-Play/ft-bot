package bot

import (
	"fmt"
	"ft-bot/bd"
	"ft-bot/config"
	"github.com/bwmarrin/discordgo"
	logger "github.com/sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var BotID string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bd.ConnectDatabase()

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Bot is running!")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		if m.Author.ID == BotID {
			return
		}

		var vars []string
		var content string
		inputSplit := strings.Split(m.Content, " ")
		for idx := range inputSplit {
			if idx == 0 {
				content = inputSplit[idx]
			}else{
				vars = append(vars, inputSplit[idx])
			}
		}


		switch content {
			case "!help":
				_, _ = s.ChannelMessageSend(m.ChannelID, "Помогу, но потом")
			case "!register":
				if len(inputSplit) < 2 {
					_, _ = s.ChannelMessageSend(m.ChannelID, "После команды !register требуется ввести код, который вы можете увидеть в личном кабинете https://lk.rimasrp.life/")
					return
				}
				playerUid, err := bd.GetUidByDiscordCode(inputSplit[1])

				if err != nil {
					fmt.Println(err.Error())
					_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
					return
				}
				if playerUid == "" {
					log := logger.WithFields(logger.Fields{"Code": inputSplit[1], "Discord":m.Author.Username})
					log.Info("Undefined code: ")
					_, _ = s.ChannelMessageSend(m.ChannelID, "Код не определен! Ваш личный код, вы можете увидеть в личном кабинете https://lk.rimasrp.life/")
					return
				}

				if !bd.CheckRegistered(m.Author.ID) {
					_, _ = s.ChannelMessageSend(m.ChannelID, "Этот аккаунт уже закреплен! При необходимости открепления, обратитесь к администратору!")
					return
				}
				_ = bd.FireNewRegisteredUser(playerUid, m.Author.ID, m.Author.Username, m.Author.Discriminator, m.Author.Email)
				_, _ = s.ChannelMessageSend(m.ChannelID, "Вы успешно привязали свой аккаунт, ваш аккаунт на сервере синхронизирован!")

				log := logger.WithFields(logger.Fields{"Steam":playerUid, "Discord":m.Author.Username})
				log.Info("Registered new user")
		}
	}
}
