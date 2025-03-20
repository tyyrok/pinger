package main

import (
	"log"
	"sync"
)

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
