package bot

import (
	"ft-bot/db"
	"ft-bot/logger"
	"time"
)

// Start ticker routines
func StartRoutine() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				ListenQueue()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Listener of discord_queue
func ListenQueue() {
	queue := db.GetQueuePlayers()
	for _, elem := range queue {
		logger.PrintLog("%v in queue right now", elem)
		renameUser(elem)
		giveRole(elem)
	}
}
