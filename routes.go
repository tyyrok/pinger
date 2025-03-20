package main

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)


func (c *Client) start(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(
		bot, fmt.Sprintf("Good news everyone! We have a new subscriber %s", ctx.Message.From.Username), &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		},
	)
	c.addUserToStorage(ctx.Message.From.Id)
	c.saveToFile()
	
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

func (c *Client) stop(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(
		bot, fmt.Sprintf("Goodbye %s", ctx.Message.From.Username), &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		},
	)
	c.removeUserFromStorage(ctx.Message.From.Id)
	c.saveToFile()
	
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}