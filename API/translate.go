package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const libreTranslateEndpoint = "https://libretranslate.com/translate"

type translateRequest struct {
	Q            string `json:"q"`
	Source       string `json:"source"`
	Target       string `json:"target"`
	Format       string `json:"format"`
	Alternatives int    `json:"alternatives"`
	ApiKey       string `json:"api_key"`
}

type translateResponse struct {
	Alternatives   []string `json:"alternatives"`
	TranslatedText string   `json:"translatedText"`
}

func TranslateTextLibre2(text, sourceLang, targetLang string) (string, error) {
	reqBody, err := json.Marshal(translateRequest{
		Q:            text,
		Source:       sourceLang,
		Target:       targetLang,
		Format:       "text",
		Alternatives: 3,
		ApiKey:       "",
	})
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req, err := http.NewRequest("POST", libreTranslateEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var translateResp translateResponse
	if err := json.Unmarshal(bodyBytes, &translateResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	return translateResp.TranslatedText, nil
}
