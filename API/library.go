package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const openLibraryEndpoint = "http://openlibrary.org/search.json"

type Book struct {
	Title  string   `json:"title"`
	Author []string `json:"author_name"`
}

type openLibraryResponse struct {
	Docs []Book `json:"docs"`
}

func SearchBook(query string) (string, error) {
	req, err := http.NewRequest("GET", openLibraryEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("q", query)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("error response from server: %s", string(bodyBytes))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	log.Printf("Response Body: %s", string(bodyBytes)) // Log the response body

	var olResp openLibraryResponse
	if err := json.Unmarshal(bodyBytes, &olResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(olResp.Docs) == 0 {
		return "No books found for the given query.", nil
	}

	var buffer bytes.Buffer
	for _, book := range olResp.Docs {
		buffer.WriteString(fmt.Sprintf("Title: %s\n", book.Title))
		if len(book.Author) > 0 {
			buffer.WriteString(fmt.Sprintf("Author: %s\n", book.Author[0]))
		} else {
			buffer.WriteString("Author: Unknown\n")
		}
		buffer.WriteString("\n")
	}

	return buffer.String(), nil
}
