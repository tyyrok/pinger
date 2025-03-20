package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func (c *Client) pinger(name, url string, bot *gotgbot.Bot) {
	log.Printf("Url is collected %s - %s", name, url)
	is_up := true
	for {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error with url %s: %s", url, err.Error())
			continue
		}
		if resp.StatusCode > 299 && is_up {
			is_up = false
			msg := fmt.Sprintf("%s is down", name)
			c.sendMsg(msg, bot)
			log.Printf("Url is down %s - %s", name, url)
		} else if resp.StatusCode <= 299 && !is_up {
			is_up = true
			msg := fmt.Sprintf("%s is up", name)
			c.sendMsg(msg, bot)
		}

		time.Sleep(time.Minute * 1)
	}
}