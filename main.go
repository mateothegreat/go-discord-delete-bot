package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/mateothegreat/go-discord-delete-bot/messages"
	"github.com/mateothegreat/go-multilog/multilog"
)

var (
	s       *discordgo.Session
	cache   = messages.NewCache()
	token   string
	channel string
)

func init() {
	godotenv.Load()

	var err error

	flag.StringVar(&token, "token", os.Getenv("DISCORD_TOKEN"), "Discord bot token")
	flag.StringVar(&channel, "channel", os.Getenv("DISCORD_CHANNEL_ID"), "Discord channel ID")
	flag.Parse()

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		multilog.Fatal("main", "error creating discord session", map[string]interface{}{
			"error": err,
		})
	}
}

func main() {
	multilog.RegisterLogger(multilog.LogMethod("console"), multilog.NewConsoleLogger(&multilog.NewConsoleLoggerArgs{
		Level:  multilog.DEBUG,
		Format: multilog.FormatText,
	}))

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		cache.Add(m.Message)
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		message, err := cache.GetByMessageID(m.ID)
		if err != nil {
			multilog.Error("main", "error fetching message", map[string]interface{}{
				"error": err,
			})
			return
		}

		_, err = s.ChannelMessageSendComplex(channel, &discordgo.MessageSend{
			Content: "A message was deleted.",
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:  discordgo.EmbedTypeRich,
					Title: "Message Deleted",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Description",
							Value:  fmt.Sprintf("Message deleted in channel <#%s> by %s", m.ChannelID, message.Author.Username),
							Inline: false,
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("Reported at %s", time.Now().Format(time.RFC1123)),
					},
				},
			},
		})
		if err != nil {
			multilog.Error("main", "error sending message", map[string]interface{}{
				"error": err,
			})
		}
	})

	// Open the discord session.
	err := s.Open()
	if err != nil {
		multilog.Fatal("main", "open discord session", map[string]interface{}{
			"error": err,
		})
	}
	defer s.Close()

	// Wait for the user to close the program.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	multilog.Info("main", "gracefully shutting down", map[string]interface{}{})
}
