package bot

import (
	"ft-bot/config"
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
)

var (
	Lg *log.Logger
)

// discord session
var s *discordgo.Session

// Starting database
func Start() {
	var err error
	s, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		logger.PrintLog(err.Error())
		return
	}

	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	s.AddHandler(onUserChanged)
	s.AddHandler(onMessageHandle)
	s.AddHandler(onCommandsCall)
	s.AddHandler(onUserConnected)
	s.AddHandler(onUserDisconnected)
	s.AddHandler(onReactMessage)
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.PrintLog("Bot is up!")
	})

	err = s.Open()
	if err != nil {
		logger.PrintLog("Cannot open the session: %v", err)
		return
	}
	defer s.Close()

	// Удаляем и тут же их добавляем, потому что дискорд принимает изменения очень долго
	for _, elem := range s.State.Guilds {
		log.Printf("Guild: %s\n", elem.ID)
		AddRemoveCommands(elem.ID)
	}

	logger.PrintLog("Start goroutines")
	StartRoutine()
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	logger.PrintLog("Gracefully shutdown\n************************************************************************\n\n")
}

// Add commands to bot
func addRemoveCommands(guildId string) {
	logger.PrintLog("Init commands...")

	cmd, err := s.ApplicationCommands(s.State.User.ID, guildId)
	if err != nil {
		logger.PrintLog(err.Error())
	}

	insert := true
	for _, v := range commands {
		for _, elem := range cmd {
			if elem == v {
				insert = false
				break
			}
		}
		if !insert {
			continue
		}
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, v)
		if err != nil {
			logger.PrintLog("Cannot create '%v' command: %v", v.Name, err)
		}
		logger.PrintLog("Command %v created", v.Name)

	}
	s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildId, commands)

	logger.PrintLog("Init commands finished")
}
