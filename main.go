package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
     message  "go-chatbot/Message"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// .env dosyasındaki çevresel değişkenleri yükleyin
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Bot token'ınızı çevresel değişkenden alın
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatalf("No token provided. Please set DISCORD_BOT_TOKEN environment variable.")
	}
	weather_token := os.Getenv("OPEN_WEATHER_TOKEN")
	if weather_token == "" {
		log.Fatalf("No token provided. Please set WEATHER_API_TOKEN environment variable.")
	}

	// Bot oturumunu başlatın
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Mesaj oluşturma olayını işleyin
	dg.AddHandler(message.MessageCreate)

	// Botu başlatın
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
