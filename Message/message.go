package message

import (
	"fmt"
	weather "go-chatbot/API"
	"log"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	if m.Content == "!hava" {
		location := "Istanbul"
		s.ChannelMessageSend(m.ChannelID, "Hava durumu bilgileri alınıyor...")
		weather, err := weather.GetWeather(location)
		if err != nil {
			log.Printf("Error getting weather: %v", err)
			s.ChannelMessageSend(m.ChannelID, "Üzgünüm, hava durumu bilgisi alınamadı.")
			return
		}

		response := fmt.Sprintf("Şu anda %s şehrinde hava durumu: %s, sıcaklık: %.2f°C, nem: %d%%", location, weather.Description, weather.Temperature, weather.Humidity)
		s.ChannelMessageSend(m.ChannelID, response)
	}

}
