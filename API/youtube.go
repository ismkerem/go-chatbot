package api

import (
    "context"
    "fmt"
    "os"

    "google.golang.org/api/option"
    "google.golang.org/api/youtube/v3"
)
type YouTubeSearchResult struct {
    Items []struct {
        Snippet struct {
            Title       string `json:"title"`
            Description string `json:"description"`
        } `json:"snippet"`
        ID struct {
            VideoID string `json:"videoId"`
        } `json:"id"`
    } `json:"items"`
}

func SearchYouTube(query string) (*youtube.SearchListResponse, error) {
        apiKey := os.Getenv("YOUTUBE_API_KEY")
        ctx := context.Background()
        service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
        if err != nil {
            return nil, fmt.Errorf("Error creating new YouTube client: %v", err)
        }
    
        call := service.Search.List([]string{"snippet"}).Q(query).MaxResults(5)
        response, err := call.Do()
        if err != nil {
            return nil, fmt.Errorf("Error making search API call: %v", err)
        }
    
        return response, nil
    }
