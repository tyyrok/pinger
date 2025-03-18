package main

import (
	"fmt"
	"log"
	//"os"
	"time"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func main() {
	//token := os.Getenv("STATUS_BOT_TOKEN")
	token := "7876917523:AAHS_o8wgmRh23thdQhCCoer7c48zwMiTjo"
	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	bot, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("Failed to create new bot: " + err.Error())
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(bot *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)
	dispatcher.AddHandler(handlers.NewCommand("start", start))

	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started....\n", bot.User.Username)
	updater.Idle()
}

func start(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(
		bot, fmt.Sprintf("Hello @%s, I'm a services unavailble bot\nYou are subsribed to notifications", bot.User.Username), &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}