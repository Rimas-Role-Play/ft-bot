package bot

import (
	"flag"
	"ft-bot/config"
	"ft-bot/logger"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
)

var (
	GuildID        = flag.String("guild", config.GuildId, "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", config.Token, "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
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

	s.AddHandler(OnMessageHandle)
	s.AddHandler(OnCommandsCall)
	s.AddHandler(OnUserConnected)
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {logger.PrintLog("Bot is up!")})

	err = s.Open()
	if err != nil {
		logger.PrintLog("Cannot open the session: %v", err)
		return
	}
	defer s.Close()

	// Удаляем и тут же их добавляем, потому что дискорд принимает изменения очень долго
	AddRemoveCommands()

	logger.PrintLog("Start goroutines")
	StartRoutine()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	logger.PrintLog("Gracefully shutdown\n************************************************************************\n\n")
}

// Add commands to bot
func AddRemoveCommands() {
	logger.PrintLog("Init commands...")

	cmd, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	if err != nil {
		logger.PrintLog(err.Error())
	}

	for _, elem := range cmd {
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, elem.ID)
		if err != nil {
			logger.PrintLog("Cant delete command %v",elem.Name)
			logger.PrintLog(err.Error())
		}else{
			logger.PrintLog("Command %v deleted",elem.Name)
		}
	}

	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			logger.PrintLog("Cannot create '%v' command: %v", v.Name, err)
		}
		logger.PrintLog("Command %v created", v.Name)
	}
	s.ApplicationCommandBulkOverwrite(s.State.User.ID, *GuildID, commands)

	logger.PrintLog("Init commands finished")
}