package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const searchURL = "https://www.googleapis.com/customsearch/v1"

type SearchResult struct {
	Items []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		Description string `json:"snippet"`
	} `json:"items"`
}

func GoogleSearch(query string) (*SearchResult, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	cx := os.Getenv("SEARCH_ENGINE_ID")
	if apiKey == "" || cx == "" {
		return nil, fmt.Errorf("API key or Search Engine ID not set")
	}

	encodedQuery := url.QueryEscape(query)
	url := fmt.Sprintf("%s?q=%s&key=%s&cx=%s", searchURL, encodedQuery, apiKey, cx)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 response status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
