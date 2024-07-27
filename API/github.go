package api

import (
	"context"
	"fmt"
	"os"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func SearchGitHubRepos(query string) (*github.RepositoriesSearchResult, error) {
	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GitHub token is not set")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	opts := &github.SearchOptions{
		Sort:       "stars",
		Order:      "desc",
		ListOptions: github.ListOptions{PerPage: 5},  
	}
	results, _, err := client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search repositories: %w", err)
	}
	return results, nil
}
