package bot

import (
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "help",
			Description: "Помогите мне пользоваться мной",
		},
		{
			Name:        "gethim",
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
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Используй команду gethim чтобы получить данные о пользователе!",
				},
			})
		},
		"gethim": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			GetHim(s,i)
		},
	}
)