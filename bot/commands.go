package bot

import (
	"ft-bot/db"
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Бот для администрирования сервера Rimas, функционал доступен только администраторам",
		},
		{
			Name:        "copy-role",
			Description: "Скопирвать роль",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role-option",
					Description: "Выберите роль",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "role-name",
					Description: "Установите имя для новой роли",
					Required:    true,
				},
			},
		},
		{
			Name:        "help-player",
			Description: "Много ответов на много вопросов",
		},
		{
			Name:        "delete-undefined-users",
			Description: "Удаляет с базы неизвестных пользователей",
		},
		{
			Name:        "re-role",
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
		{
			Name:        "clean-channel",
			Description: "Удалить все сообщения в канале",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel-option",
					Description: "Channel option",
					Required:    true,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"copy-role": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !isDiscordAdmin(s, i.Member.User.ID) {
				return
			}
			copyRole(s, i)
		},
		"help-player": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			printHelp(s, i)
		},
		"clean-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !isDiscordAdmin(s, i.Member.User.ID) {
				return
			}
			go cleanMessageInChannel(s, i)
			return
		},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			printHiddenMessage(s, i, "Бот для управления и администрирования сервера Rimas Life")
		},
		"get-him": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			getHim(s, i)
		},
		"delete-undefined-users": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !isDiscordAdmin(s, i.Member.User.ID) {
				printHiddenMessage(s, i, "У вас нет доступа")
				return
			}
			printHiddenMessage(s, i, "Удаляю всех неизвестных...")
			deleteUndefinedUsers()
		},
		"re-role": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !isDiscordAdmin(s, i.Member.User.ID) {
				return
			}
			printHiddenMessage(s, i, "Перепроверяю все роли...")
			go giveRoles()
			logger.PrintLog("reRole called")
		},
		"re-name": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].UserValue(nil)
			pid := user.ID
			sender := i.Interaction.Member.User
			if !isDiscordAdmin(s, sender.ID) {
				printHiddenMessage(s, i, "У вас нет доступа")
				return
			}
			player, err := db.GetUserByDS(pid)
			if err != nil {
				printHiddenMessage(s, i, "Пользователь не найден")
				return
			}
			renameUser(player.PlayerInfo.SteamId)
			printHiddenMessage(s, i, "Запрос отправлен")
		},
		"give-boost": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].UserValue(nil)
			sender := i.Interaction.Member.User
			if !isDiscordAdmin(s, sender.ID) {
				printHiddenMessage(s, i, "У вас нет доступа")
				return
			}
			giveBoostPresent(i.ChannelID, user)
			printHiddenMessage(s, i, "Запрос отправлен")
		},
	}
)
