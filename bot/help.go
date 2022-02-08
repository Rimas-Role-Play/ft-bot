package bot

import (
	"fmt"
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
)

func printHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response *discordgo.InteractionResponse
	response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Выберите вопрос из списка ниже 👇",
			Flags:   1 << 6,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							// Select menu, as other components, must have a customID, so we set it to this value.
							CustomID:    "questions",
							Placeholder: "Выберите вопрос",
							Options: []discordgo.SelectMenuOption{
								{
									Label: "Рация",
									// As with components, this things must have their own unique "id" to identify which is which.
									// In this case such id is Value field.
									Value: "tfar",
									Emoji: discordgo.ComponentEmoji{
										Name: "🦦",
									},
									// You can also make it a default option, but in this case we won't.
									Default:     false,
									Description: "Ссылка на плагин",
								},
								{
									Label: "Как начать играть?",
									Value: "how2play",
									Emoji: discordgo.ComponentEmoji{
										Name: "🟨",
									},
									Description: "Инструкция",
								},
								{
									Label: "Python",
									Value: "py",
									Emoji: discordgo.ComponentEmoji{
										Name: "🐍",
									},
									Description: "Python programming language",
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
		logger.PrintLog("[ERROR] printHelp: %s\n", err.Error())
	}
}

var (
	helpCommands = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

func selectHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response *discordgo.InteractionResponse
	fmt.Printf("select help\n")
	data := i.MessageComponentData()
	switch data.Values[0] {
	case "tfar":
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					generateEmbed("Скачать последнюю версию плагина можно по этой ссылке\nhttps://rimasrp.life/task_force_radio.ts3_plugin", "Плагин рации", "https://rimasrp.life/task_force_radio.ts3_plugin"),
				},
			},
		}
	case "how2play":
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					generateEmbed("Как начать играть?\nСледуйте инструкции ниже\nШаг 1\nКупите и скачайте ArmA 3 в Steam.\nhttps://store.steampowered.com/app/107410/Arma_3/?l=russian\nШаг 2\nПодпишитесь на мод Rimas Role Play в мастерской Steam.\nhttps://steamcommunity.com/sharedfiles/filedetails/?id=1368860933\nШаг 3\nСкачайте клиент TeamSpeak и установите его.\nhttps://files.teamspeak-services.com/releases/client/3.5.6/TeamSpeak3-Client-win64-3.5.6.exe\nШаг 4\nСкачайте плагин Task Force Radio и установите его.\nhttps://rimasrp.life/task_force_radio.ts3_plugin\nЗапуск\nЗапустите ArmA 3 в Steam, кликнув на кнопку играть.\n\nВ пункте \"Моды\" проверьте, включен ли мод Rimas Role Play, если отключен — включите его.\n\nНажмите на оранжевую кнопку играть в лаунчере ArmA 3.\n\nВ правом верхнем углы игры зайдите в свой профиль и укажите имя и фамилию вашего персонажа.\n\nЗайдите в браузер серверов и нажмите прямое подключение.\n\nS1.RIMASRP.LIFE\n", "Как начать играть?", ""),
				},
			},
		}
	default:
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "It is not the way to go.",
			},
		}
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		panic(err)
	}
}
