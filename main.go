package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
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
	storage, _ := loadFromFile()
	c := &Client{storage: storage}

	dispatcher.AddHandler(handlers.NewCommand("start", c.start))
	dispatcher.AddHandler(handlers.NewCommand("stop", c.stop))

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

	c.pingHosts(bot)
	updater.Idle()
}

func loadFromFile() (Storage, error) {
	file, err := os.Open("settings.json")
	if err != nil {
		panic("Failed to open file: " + err.Error())
	}
	defer file.Close()
	var store Storage
	err = json.NewDecoder(file).Decode(&store)
	if err != nil {
		panic ("Failed to decode file content")
	}
	log.Printf("Found %d hosts", len(store.Hosts))
	log.Printf("Found %d users", len(store.Users))
	return store, nil
}
