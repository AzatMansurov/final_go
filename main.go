package main

import (
	"upgrade/cmd/bot"
)

func main() {

	upgradeBot := bot.CreateBot()

	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)

	go func() {
		upgradeBot.Bot.Start()
	}()

	bot.Listen()
}
