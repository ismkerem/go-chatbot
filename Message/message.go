package message

import (
	"fmt"
	"log"
	"strings"
	"sync"

	api "go-chatbot/api"
	"github.com/bwmarrin/discordgo"
)

var userStates = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var translationRequests = struct {
	sync.RWMutex
	m map[string]TranslationRequest
}{m: make(map[string]TranslationRequest)}

type TranslationRequest struct {
	SourceLang string
	TargetLang string
	Text       string
	State      string
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	log.Printf("Received message from %s: %s", m.Author.Username, m.Content)

	userStates.Lock()
	state, hasState := userStates.m[m.Author.ID]
	userStates.Unlock()

	if hasState {
		handleState(s, m, state)
		return
	}

	switch {
	case m.Content == "!merhaba":
		s.ChannelMessageSend(m.ChannelID, "Merhaba! Sana nasıl yardımcı olabilirim?")
	case m.Content == "!komut":
		s.ChannelMessageSend(m.ChannelID, "Komutlar: \n!hava [şehir] - Hava durumu bilgisi alın\n!ara [arama terimi] - Google'da arama yapın\n!youtube [arama terimi] - YouTube'da arama yapın\n!github [arama terimi] - GitHub'da arama yapın\n!çevir - Metin çevirisi yapın\n!book [kitap adı] - Kitap araması yapın")
	case strings.HasPrefix(m.Content, "!hava"):
		handleWeather(s, m)
	case strings.HasPrefix(m.Content, "!ara"):
		setUserState(m.Author.ID, "awaiting_search_query")
		s.ChannelMessageSend(m.ChannelID, "Ne aramak istersiniz?")
	case strings.HasPrefix(m.Content, "!youtube"):
		setUserState(m.Author.ID, "awaiting_youtube_query")
		s.ChannelMessageSend(m.ChannelID, "YouTube'da ne aramak istersiniz?")
	case strings.HasPrefix(m.Content, "!github"):
		setUserState(m.Author.ID, "awaiting_github_query")
		s.ChannelMessageSend(m.ChannelID, "GitHub'da ne aramak istersiniz?")
	case strings.HasPrefix(m.Content, "!çevir"):
		setUserState(m.Author.ID, "awaiting_source_language")
		s.ChannelMessageSend(m.ChannelID, "Hangi dilden çevirmek istiyorsunuz? (örnek: en, tr)")
	case strings.HasPrefix(m.Content, "!book"):
		handleBookSearch(s, m)
	}
}

func setUserState(userID, state string) {
	userStates.Lock()
	defer userStates.Unlock()
	userStates.m[userID] = state
}

func handleWeather(s *discordgo.Session, m *discordgo.MessageCreate) {
	location := strings.TrimSpace(strings.TrimPrefix(m.Content, "!hava"))

	if location == "" {
		s.ChannelMessageSend(m.ChannelID, "Lütfen hava durumu için bir konum girin. Örneğin, '!hava Istanbul'.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hava durumu bilgileri %s için alınıyor...", location))
	weatherInfo, err := api.GetWeather(location)
	if err != nil {
		log.Printf("Error getting weather for %s: %v", location, err)
		s.ChannelMessageSend(m.ChannelID, "Üzgünüm, hava durumu bilgisi alınamadı.")
		return
	}

	response := fmt.Sprintf("Şu anda %s şehrinde hava durumu: %s, sıcaklık: %.2f°C, nem: %d%%", location, weatherInfo.Description, weatherInfo.Temperature, weatherInfo.Humidity)
	s.ChannelMessageSend(m.ChannelID, response)
}

func handleState(s *discordgo.Session, m *discordgo.MessageCreate, state string) {
	switch state {
	case "awaiting_search_query":
		handleSearchQuery(s, m)
	case "awaiting_youtube_query":
		handleYouTubeQuery(s, m)
	case "awaiting_github_query":
		handleGitHubQuery(s, m)
	case "awaiting_source_language":
		handleSourceLanguage(s, m)
	case "awaiting_target_language":
		handleTargetLanguage(s, m)
	case "awaiting_translation_text":
		handleTranslationText(s, m)
	}
}

func handleSearchQuery(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	setUserState(m.Author.ID, "")
}

func handleYouTubeQuery(s *discordgo.Session, m *discordgo.MessageCreate) {
	query := m.Content
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'%s' için YouTube araması yapılıyor...", query))
	results, err := api.SearchYouTube(query)
	if err != nil {
		log.Printf("Error performing YouTube search: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Üzgünüm, YouTube araması yapılırken bir hata oluştu.")
		return
	}

	if len(results.Items) > 0 {
		for _, item := range results.Items {
			response := fmt.Sprintf("Başlık: %s\nLink: https://www.youtube.com/watch?v=%s\nAçıklama: %s\n\n", item.Snippet.Title, item.Id.VideoId, item.Snippet.Description)
			s.ChannelMessageSend(m.ChannelID, response)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Arama sonucunda hiçbir şey bulunamadı.")
	}

	setUserState(m.Author.ID, "")
}

func handleGitHubQuery(s *discordgo.Session, m *discordgo.MessageCreate) {
	query := m.Content
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'%s' için GitHub araması yapılıyor...", query))
	results, err := api.SearchGitHubRepos(query)
	if err != nil {
		log.Printf("Error performing GitHub search: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Üzgünüm, GitHub araması yapılırken bir hata oluştu: %v", err))
		return
	}

	if len(results.Repositories) > 0 {
		for i, item := range results.Repositories {
			if i >= 5 {
				break
			}
			response := fmt.Sprintf("Repo Adı: %s\nLink: %s\nAçıklama: %s\n\n", item.GetFullName(), item.GetHTMLURL(), item.GetDescription())
			s.ChannelMessageSend(m.ChannelID, response)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Arama sonucunda hiçbir şey bulunamadı.")
	}

	setUserState(m.Author.ID, "")
}

func handleSourceLanguage(s *discordgo.Session, m *discordgo.MessageCreate) {
	sourceLang := m.Content
	translationRequests.Lock()
	translationRequests.m[m.Author.ID] = TranslationRequest{
		SourceLang: sourceLang,
		State:      "awaiting_target_language",
	}
	translationRequests.Unlock()

	s.ChannelMessageSend(m.ChannelID, "Hangi dile çevirmek istiyorsunuz? (örnek: en, tr)")

	setUserState(m.Author.ID, "awaiting_target_language")
}

func handleTargetLanguage(s *discordgo.Session, m *discordgo.MessageCreate) {
	targetLang := m.Content
	translationRequests.Lock()
	req := translationRequests.m[m.Author.ID]
	req.TargetLang = targetLang
	req.State = "awaiting_translation_text"
	translationRequests.m[m.Author.ID] = req
	translationRequests.Unlock()

	s.ChannelMessageSend(m.ChannelID, "Çevirmek istediğiniz metni girin:")

	setUserState(m.Author.ID, "awaiting_translation_text")
}

func handleTranslationText(s *discordgo.Session, m *discordgo.MessageCreate) {
	text := m.Content
	translationRequests.Lock()
	req := translationRequests.m[m.Author.ID]
	req.Text = text
	translationRequests.m[m.Author.ID] = req
	translationRequests.Unlock()

	translatedText, err := api.TranslateTextLibre2(req.Text, req.SourceLang, req.TargetLang)
	if err != nil {
		log.Printf("Error performing translation: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Üzgünüm, çeviri yapılırken bir hata oluştu.")
		return
	}

	log.Printf("Translated Text: %s", translatedText)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Çeviri:\n%s", translatedText))

	setUserState(m.Author.ID, "")
	clearTranslationRequest(m.Author.ID)
}

func clearTranslationRequest(userID string) {
	translationRequests.Lock()
	defer translationRequests.Unlock()
	delete(translationRequests.m, userID)
}

func handleBookSearch(s *discordgo.Session, m *discordgo.MessageCreate) {
	query := strings.TrimSpace(strings.TrimPrefix(m.Content, "!book"))
	if query == "" {
		s.ChannelMessageSend(m.ChannelID, "Lütfen aramak istediğiniz kitabın adını girin. Örneğin, '!book Harry Potter'.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'%s' kitabı için arama yapılıyor...", query))
	bookInfo, err := api.SearchBook(query)
	if err != nil {
		log.Printf("Error searching for book: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Üzgünüm, kitap araması yapılırken bir hata oluştu.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, bookInfo)
}
