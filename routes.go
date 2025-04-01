package main

import (
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)


func (c *Client) start(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(
		bot, fmt.Sprintf("Good news everyone! We have a new subscriber %s", ctx.Message.From.Username), &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	c.addUserToStorage(ctx.Message.From.Id)
	c.saveToFile()
	go c.welcomeMessage(ctx.Message.From.Id, bot)
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

func (c *Client) status(bot *gotgbot.Bot, ctx *ext.Context) error {
	if !c.checkUserSubscribed(ctx.Message.From.Id) {
		_, err := ctx.EffectiveMessage.Reply(
			bot, "You're not subscribed", &gotgbot.SendMessageOpts{
				ParseMode: "HTML",
			},
		)
		if err != nil {
			return fmt.Errorf("failed to send start message: %w", err)
		}
	} else if c.checkUserFlooding(ctx.Message.From.Id) {
		_, err := ctx.EffectiveMessage.Reply(
			bot, "Too many request", &gotgbot.SendMessageOpts{
				ParseMode: "HTML",
			},
		)
		if err != nil {
			return fmt.Errorf("failed to send start message: %w", err)
		}
	} else {
		c.setUserDelay(ctx.Message.From.Id)
		c.welcomeMessage(ctx.Message.From.Id, bot)
		time.Sleep(time.Second * 5)
		c.removeUserDelay(ctx.Message.From.Id)
	}
	return nil
}