package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"moul.io/u"
)

func main() {
	var (
		token   = os.Getenv("DISCORD_BOT_TOKEN")
		guildID = os.Getenv("DISCORD_GUILD_ID")
	)

	s, err := discordgo.New("Bot " + token)
	u.CheckErr(err)

	// commands
	{
		commands := []*discordgo.ApplicationCommand{
			{Name: "sesame", Description: "ouvre toi !"},
			{Name: "status", Description: "kezispass ?"},
		}
		commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
			"sesame": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				fmt.Println("TEST3")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "FAUT PAS REVER, CA MARCHE PAS.",
					},
				})
			},
			"status": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "FAUT PAS REVER, CA MARCHE PAS !",
					},
				})
			},
		}
		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fmt.Println("TEST1")
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				fmt.Println("TEST2")
				h(s, i)
			}
		})
		s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			log.Println("Bot is up!")
		})
		u.CheckErr(s.Open())
		defer s.Close()
		for _, v := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
			u.CheckErr(err)
		}
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")
}
