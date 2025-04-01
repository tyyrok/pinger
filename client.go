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
	Users map[int64]bool
}

type Client struct {
	rwMux sync.RWMutex
	storage Storage
}

func (c *Client) addUserToStorage(user_id int64) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	if c.storage.Users == nil {
		c.storage.Users = make(map[int64]bool)
		c.storage.Users[user_id] = false
		return
	}
	if c.storage.Users[user_id] {
		return
	}
	c.storage.Users[user_id] = false
	log.Printf("Subscribed user with id: %d", user_id)
}

func (c *Client) removeUserFromStorage(user_id int64) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	if c.storage.Users == nil {
		log.Printf("Storage is empty can't remove user with id: %d", user_id)
		return
	}
	delete(c.storage.Users, user_id)
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
	for user_id, _ := range c.storage.Users {
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

func (c *Client) checkUserFlooding(user_id int64) bool {
	c.rwMux.RLock()
	defer c.rwMux.RUnlock()
	return c.storage.Users[user_id]
}

func (c *Client) setUserDelay(user_id int64) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	_, ok := c.storage.Users[user_id]
	if ok {
		c.storage.Users[user_id] = true
	}
}

func (c *Client) removeUserDelay(user_id int64) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	_, ok := c.storage.Users[user_id]
	if ok {
		c.storage.Users[user_id] = false
	}
}

func (c *Client) checkUserSubscribed(user_id int64) bool {
	c.rwMux.RLock()
	defer c.rwMux.RUnlock()
	_, ok := c.storage.Users[user_id]
	return ok
}