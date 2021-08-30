package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"moul.io/u"
)

func main() {
	var (
		token       = os.Getenv("DISCORD_BOT_TOKEN")
		guildID     = os.Getenv("DISCORD_GUILD_ID")
		apiEndpoint = os.Getenv("API_ENDPOINT")
	)

	s, err := discordgo.New("Bot " + token)
	u.CheckErr(err)

	// create tls client
	var client *http.Client
	{
		cert, err := tls.LoadX509KeyPair("usercert.pem", "userkey.pem")
		u.CheckErr(err)
		caCertPool := x509.NewCertPool()
		caCert, err := ioutil.ReadFile("cacert.pem")
		u.CheckErr(err)
		caCertPool.AppendCertsFromPEM(caCert)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caCertPool,
			},
		}
		client = &http.Client{Transport: tr}
	}

	apiOpen := func() (string, error) {
		form := url.Values{}
		form.Add("action", "open")
		form.Add("delay", "5")
		req, err := http.NewRequest("POST", apiEndpoint+"/api.php", strings.NewReader(form.Encode()))
		if err != nil {
			return "", fmt.Errorf("creating new HTTP request: %w", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("making HTTP request using client: %w", err)
		}
		respBody, _ := ioutil.ReadAll(resp.Body)
		log.Printf("reply: %q", resp.Body)
		return string(respBody), nil
	}

	apiStatus := func() (string, error) {
		req, err := http.NewRequest("GET", apiEndpoint+"/api.php", nil)
		if err != nil {
			return "", fmt.Errorf("creating new HTTP request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("making HTTP request using client: %w", err)
		}
		respBody, _ := ioutil.ReadAll(resp.Body)
		log.Printf("reply: %q", resp.Body)
		return string(respBody), nil
	}

	{
		status, err := apiStatus()
		u.CheckErr(err)
		log.Printf("status on init: %q", status)
	}

	// commands
	{
		commands := []*discordgo.ApplicationCommand{
			{Name: "sesame", Description: "ouvre toi !"},
			{Name: "status", Description: "kezispass ?"},
		}
		commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
			"sesame": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				log.Println("/sesame called")
				status, err := apiOpen()

				if err != nil {
					log.Println("error: %+v", err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("ERROR: %q", err),
						},
					})
				} else {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("NORMALEMENT C'EST BON: %q", status),
						},
					})
				}
			},
			"status": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				log.Println("/status called")
				status, err := apiStatus()

				if err != nil {
					log.Println("error: %+v", err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("ERROR: %q", err),
						},
					})
				} else {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("NORMALEMENT C'EST BON: %q", status),
						},
					})
				}
			},
		}
		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
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
