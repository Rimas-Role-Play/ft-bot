package bot

import (
	"ft-bot/bd"
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "help",
			Description: "Бот для администрирования сервера Rimas, функционал доступен только администраторам",
		},
		{
			Name: "delete-undefined-users",
			Description: "Удаляет с базы неизвестных пользователей",
		},
		{
			Name: "re-role",
			Description: "Перепроверяет выданные роли",
		},
		{
			Name:        "get-him",
			Description: "Получить данные игрока",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "Тегните пользователя",
					Required:    true,
				},
			},
		},
		{
			Name:        "re-name",
			Description: "Задать игровое имя пользователю",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "Тегните пользователя",
					Required:    true,
				},
			},
		},
		{
			Name:        "give-boost",
			Description: "Выдать подарок за буст сервера",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "Тегните пользователя",
					Required:    true,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			PrintHiddenMessage(s,i,"Бот для управления и администрирования сервера Rimas Life")
		},
		"get-him": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			GetHim(s,i)
		},
		"re-name-all": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go GetHim(s,i)
		},
		"delete-undefined-users": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !IsDiscordAdmin(s, i.Member.User.ID) {
				PrintHiddenMessage(s,i,"У вас нет доступа")
				return
			}
			PrintHiddenMessage(s,i,"Удаляю всех неизвестных...")
			DeleteUndefinedUsers()
		},
		"re-role": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !IsDiscordAdmin(s, i.Member.User.ID) {
				return
			}
			PrintHiddenMessage(s,i,"Перепроверяю все роли...")
			go GiveRoles()
			logger.PrintLog("reRole called")
		},
		"re-name": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].UserValue(nil)
			pid := user.ID
			sender := i.Interaction.Member.User
			if !IsDiscordAdmin(s, sender.ID) {
				PrintHiddenMessage(s,i,"У вас нет доступа")
				return
			}
			player, err := bd.GetUserByDS(pid)
			if err != nil {
				PrintHiddenMessage(s,i,"Пользователь не найден")
				return
			}
			RenameUser(player.PlayerInfo.SteamId)
			PrintHiddenMessage(s,i,"Запрос отправлен")
		},
		"give-boost": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].UserValue(nil)
			sender := i.Interaction.Member.User
			if !IsDiscordAdmin(s, sender.ID) {
				PrintHiddenMessage(s, i, "У вас нет доступа")
				return
			}
			giveBoostPresent(i.ChannelID,user)
			PrintHiddenMessage(s,i,"Запрос отправлен")
		},
	}
)