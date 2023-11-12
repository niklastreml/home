package api

import (
	"context"

	"github.com/google/go-github/v56/github"
)

func GetProjects(ctx context.Context, client *github.Client) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.List(ctx, "NiklasTreml", &github.RepositoryListOptions{
		Type: "all",
	})
	if err != nil {
		return nil, err
	}

	return repos, nil
}
