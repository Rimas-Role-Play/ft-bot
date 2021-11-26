package bot

import (
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
)

func helpFaq(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Выберите вопрос из списка ниже",
			Flags:   1 << 6,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							// Select menu, as other components, must have a customID, so we set it to this value.
							CustomID:    "select",
							Placeholder: "Выберите вопрос 👇",
							Options: []discordgo.SelectMenuOption{
								{
									Label: "Рация",
									// As with components, this things must have their own unique "id" to identify which is which.
									// In this case such id is Value field.
									Value: "go",
									Emoji: discordgo.ComponentEmoji{
										Name: "🦦",
									},
									// You can also make it a default option, but in this case we won't.
									Default:     false,
									Description: "Где скачать и как пользоваться рацией",
								},
								{
									Label: "Начинаем игру",
									Value: "js",
									Emoji: discordgo.ComponentEmoji{
										Name: "🟨",
									},
									Description: "Как начать играть, где скачать игру и как установить",
								},
								{
									Label: "Выгоняет...",
									Value: "py",
									Emoji: discordgo.ComponentEmoji{
										Name: "🐍",
									},
									Description: "Что можно сделать, чтобы справиться с проблемой",
								},
							},
						},
					},
				},
			},
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		logger.PrintLog(err.Error())
	}
}