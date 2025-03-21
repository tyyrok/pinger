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
	var start time.Time
	var msg string
	cl := &http.Client{
		Timeout: 8 * time.Second,
	}
	for {
		resp, err := cl.Get(url)
		if err != nil {
			if is_up {
				log.Printf("Error with url %s: %s", url, err.Error())
				is_up = false
				start = time.Now()
				msg := fmt.Sprintf("%s - is down", name)
				c.sendMessages(msg, bot)
			} else {
				time.Sleep(time.Minute * 1)
				continue
			}

		} else if resp.StatusCode > 299 && is_up {
			is_up = false
			start = time.Now()
			msg := fmt.Sprintf("%s - is down", name)
			c.sendMessages(msg, bot)
			log.Printf("Url is down %s - %s", name, url)

		} else if resp.StatusCode <= 299 && !is_up {
			is_up = true
			if !start.IsZero() {
				msg = fmt.Sprintf("%s - is up\n (downtime: %v)", name, time.Since(start).Round(time.Second))
				start = time.Time{}
			} else {
				msg = fmt.Sprintf("%s - is up", name)
			}
			c.sendMessages(msg, bot)
		}

		time.Sleep(time.Minute * 1)
	}
}

func (c *Client) getHostStatus(msgs chan string, name, url string) {
	cl := &http.Client{
		Timeout: 8 * time.Second,
	}
	resp, err := cl.Get(url)
	if err != nil {
		log.Printf("Error with url %s: %s", url, err.Error())
		msgs <- fmt.Sprintf("%s - is down", name)
		return
	}
	if resp.StatusCode > 299 {
		msgs <- fmt.Sprintf("%s - is down", name)
	} else {
		msgs <- fmt.Sprintf("%s - is ok", name)
	}
}