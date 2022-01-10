package bot

import (
	"fmt"
	"ft-bot/config"
	"ft-bot/db"
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
	"log"
)

var vipRoleId = "871508004795723787"
var regRoleId = "864630308242849825"

func renameUser(pid string) {
	var err error
	player, err := db.GetPlayer(pid)
	if err != nil {
		logger.PrintLog(err.Error())
		return
	}
	_, err = s.GuildMember(config.GuildId, pid)
	if err != nil {
		logger.PrintLog("User not found")
		return
	}
	err = s.GuildMemberNickname(config.GuildId, player.PlayerInfo.DSUid, player.PlayerInfo.Name)
	if err != nil {
		logger.PrintLog("Cant rename %v user", player.PlayerInfo.Name)
		logger.PrintLog(err.Error())
	}
}

// Check is discord admin
func isDiscordAdmin(s *discordgo.Session, id string) bool {
	g, err := s.GuildMember(config.GuildId, id)
	if err != nil {
		return false
	}
	roles := g.Roles
	for _, i := range config.GetAdminRoles() {
		for _, y := range roles {
			if i == y {
				log.Println("Admin found")
				return true
			}
		}
	}
	log.Println("Admin not found")
	return false
}

// Check is user have Mute role
func isMuted(s *discordgo.Session, id string) bool {
	g, err := s.GuildMember(config.GuildId, id)
	if err != nil {
		return false
	}
	fmt.Println(g)
	roles := g.Roles

	for _, y := range roles {
		if y == "875120846123962429" {
			log.Println("User muted")
			return true
		}
	}
	return false
}

// Delete leaved or banned users
func deleteUndefinedUsers() {
	uids := db.GetAllDiscordUids()
	for _, elem := range uids {
		_, err := s.GuildMember(config.GuildId, elem)
		if err != nil {
			log.Println(err.Error())
			log.Printf("User: %v will be deleted", elem)
			db.DeleteDiscordUser(elem)
		}
	}
	log.Printf("All inactive users deleted")
}

// Send information from command get-him
func getHim(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].UserValue(nil)
	pid := user.ID
	sender := i.Interaction.Member.User
	logger.PrintLog("User %v#%v want get information about %v", sender.Username, sender.Discriminator, pid)
	if !isDiscordAdmin(s, i.Interaction.Member.User.ID) {
		logger.PrintLog("User %v#%v dont have access to /gethim", sender.Username, sender.Discriminator)
		printHiddenMessage(s, i, "У вас нет доступа")
		return
	}
	data, found := db.GetPlayerStr(pid)
	if found {
		logger.PrintLog("User %v found!", pid)
	} else {
		logger.PrintLog(data)
	}
	printHiddenMessage(s, i, data)
}

func printHiddenMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: msg,
		},
	})
}
func giveBoostPresent(channelId string, user *discordgo.User) {
	player, err := db.GetUserByDS(user.ID)
	if err != nil {
		logger.PrintLog("Cant give vehicle for boost: %v", err.Error())
		s.ChannelMessageSend(channelId, pingUser(user.ID))
		s.ChannelMessageSend(channelId, "Мы не нашли ваш аккаунт на сервере, привяжите ваш аккаунт и напишите администрации за получением бонуса")
		return
	}
	vehicle := config.GetRandomVehicle()
	db.InsertVehicle(vehicle.Classname, player.PlayerInfo.SteamId)
	s.ChannelMessageSend(channelId, pingUser(user.ID))
	s.ChannelMessageSendEmbed(channelId, createEmbedNitroBooster(vehicle))
	logger.PrintLog("%v boosted server and given %v", user.Username, vehicle.DisplayName)
}
func createEmbedNitroBooster(vehicle config.Vehicles) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		URL:         "",
		Type:        discordgo.EmbedTypeImage,
		Title:       "Nitro Booster",
		Description: fmt.Sprintf("Спасибо за буст сервера!\nТвой подарок %v уже доступен на сервере!", vehicle.DisplayName),
		Timestamp:   "",
		Color:       0x9300FF,
		Footer: &discordgo.MessageEmbedFooter{
			Text:         "Nitro Boost",
			IconURL:      "",
			ProxyIconURL: "",
		},
		Image: &discordgo.MessageEmbedImage{
			URL:      vehicle.Image,
			ProxyURL: "",
			Width:    0,
			Height:   0,
		},
	}
	return embed
}
func pingUser(id string) string {
	return fmt.Sprintf("<@%v>", id)
}
