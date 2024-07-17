package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

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

    // Bot oturumunu başlatın
    dg, err := discordgo.New("Bot " + token)
    if err != nil {
        log.Fatalf("Error creating Discord session: %v", err)
    }

    // Mesaj oluşturma olayını işleyin
    dg.AddHandler(messageCreate)

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Botun kendi mesajlarını dikkate almayın
    if m.Author.ID == s.State.User.ID {
        return
    }

    // Mesajı loglayın
    log.Printf("Received message from %s: %s", m.Author.Username, m.Content)

    // Basit bir yanıt örneği
    if m.Content == "!merhaba" {
        s.ChannelMessageSend(m.ChannelID, "Merhaba! Nasılsın?")
    }


}
