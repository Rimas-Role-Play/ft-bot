package bot

import "github.com/bwmarrin/discordgo"

var (
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"questions": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			selectHelp(s, i)
		},
	}
)
