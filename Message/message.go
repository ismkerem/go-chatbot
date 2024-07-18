package message

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"go-chatbot/api"
	"github.com/bwmarrin/discordgo"
)

var userStates = struct{
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	log.Printf("Received message from %s: %s", m.Author.Username, m.Content)

	userStates.Lock()
	state, hasState := userStates.m[m.Author.ID]
	userStates.Unlock()

	if hasState && state == "awaiting_search_query" {
		// Kullanıcıdan arama sorgusunu bekliyorsanız
		query := m.Content
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'%s' için arama yapılıyor...", query))
		results, err := api.GoogleSearch(query)
		if err != nil {
			log.Printf("Error performing search: %v", err)
			s.ChannelMessageSend(m.ChannelID, "Üzgünüm, arama yapılırken bir hata oluştu.")
			return
		}

		if len(results.Items) > 0 {
			for _, item := range results.Items {
				response := fmt.Sprintf("Başlık: %s\nLink: %s\nAçıklama: %s\n\n", item.Title, item.Link, item.Description)
				s.ChannelMessageSend(m.ChannelID, response)
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Arama sonucunda hiçbir şey bulunamadı.")
		}

		// Kullanıcı durumunu temizleyin
		userStates.Lock()
		delete(userStates.m, m.Author.ID)
		userStates.Unlock()
		return
	}

	// Basit bir yanıt örneği
	if m.Content == "!merhaba" {
		s.ChannelMessageSend(m.ChannelID, "Merhaba! Nasılsın?")
	}

	if m.Content == "!hava" {
		location := "Istanbul"
		s.ChannelMessageSend(m.ChannelID, "Hava durumu bilgileri alınıyor...")
		weatherInfo, err := api.GetWeather(location)
		if err != nil {
			log.Printf("Error getting weather: %v", err)
			s.ChannelMessageSend(m.ChannelID, "Üzgünüm, hava durumu bilgisi alınamadı.")
			return
		}

		response := fmt.Sprintf("Şu anda %s şehrinde hava durumu: %s, sıcaklık: %.2f°C, nem: %d%%", location, weatherInfo.Description, weatherInfo.Temperature, weatherInfo.Humidity)
		s.ChannelMessageSend(m.ChannelID, response)
	}

	if strings.HasPrefix(m.Content, "!ara") {
		// Kullanıcıya ne aramak istediğini sorun ve durumu güncelleyin
		s.ChannelMessageSend(m.ChannelID, "Ne aramak istersiniz?")

		userStates.Lock()
		userStates.m[m.Author.ID] = "awaiting_search_query"
		userStates.Unlock()
	}
}
