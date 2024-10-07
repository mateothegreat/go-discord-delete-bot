package messages

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type Cache struct {
	messages map[string]*discordgo.Message
}

func (c *Cache) Add(message *discordgo.Message) {
	c.messages[message.ID] = message
}

func (c *Cache) GetByMessageID(id string) (*discordgo.Message, error) {
	message, ok := c.messages[id]
	if !ok {
		return nil, errors.New("message not found")
	}
	return message, nil
}

func NewCache() *Cache {
	return &Cache{
		messages: make(map[string]*discordgo.Message),
	}
}
