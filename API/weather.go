package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Weather struct {
	Description string  `json:"description"`
	Temperature float64 `json:"temp"`
	Humidity    int     `json:"humidity"`
}

func GetWeather(location string) (*Weather, error) {
	apiKey := os.Getenv("OPEN_WEATHER_TOKEN")
	if apiKey == "" {
		return nil, fmt.Errorf("no weather API token provided")
	}
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", location, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received non-200 response status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Error decoding response: %v", err)
	}

	weather := &Weather{
		Description: result["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
		Temperature: result["main"].(map[string]interface{})["temp"].(float64),
		Humidity:    int(result["main"].(map[string]interface{})["humidity"].(float64)),
	}

	return weather, nil

}
