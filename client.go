package main

import (
	"sync"
	"encoding/json"
	"log"
	"os"
	"time"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Storage struct {
	Hosts map[string]string `json:"hosts"`
	Users []int64 `json:"users"`
}

type Client struct {
	rwMux sync.RWMutex
	storage Storage
}

func (c *Client) addUserToStorage(user_id int64) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	if c.storage.Users == nil {
		c.storage.Users = []int64{user_id}
		return
	}
	for _, elem := range c.storage.Users {
		if elem == user_id {
			return
		}
	}
	c.storage.Users = append(c.storage.Users, user_id)
	log.Printf("Subscribed user with id: %d", user_id)
}

func (c *Client) removeUserFromStorage(user_id int64) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	if c.storage.Users == nil {
		log.Printf("Storage is empty can't remove user with id: %d", user_id)
		return
	}
	index := -1
	for i, elem := range c.storage.Users {
		if elem == user_id {
			index = i
			break
		}
	}
	if index >= 0 {
		c.storage.Users = append(c.storage.Users[:index], c.storage.Users[index+1:]...)
	}
}

func (c *Client) saveToFile() error {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	file, err := os.Create("settings.json")
	if err != nil {
		panic("Failed to open file: " + err.Error())
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(c.storage)
}

func (c *Client) pingHosts(bot *gotgbot.Bot) {
	for key, value := range c.storage.Hosts {
		go c.pinger(key, value, bot)
		time.Sleep(time.Second * 1)
	}
}

func (c *Client) sendMessages(msg string, bot *gotgbot.Bot) {
	for _, user_id := range c.storage.Users {
		c.sendToUser(user_id, msg, bot)
		time.Sleep(time.Second * 1)
	}
}

func (c *Client) sendToUser(user_id int64, msg string, bot *gotgbot.Bot) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	bot.SendMessage(user_id, msg, &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
}

func (c *Client) welcomeMessage(user_id int64, bot *gotgbot.Bot) {
	messages := make(chan string)
	for name, url := range c.storage.Hosts {
		go c.getHostStatus(messages, name, url)
	}
	res := "Status:\n"
	for i := 0; i < len(c.storage.Hosts); i++ {
		res += fmt.Sprintf("%s\n", <- messages)
	}
	c.sendToUser(user_id, res, bot)
}