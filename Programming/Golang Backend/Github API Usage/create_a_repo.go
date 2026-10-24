// app must be installed on the organization
// github/organizations/<org-name>/settings

package main

import (
	"context"
	"net/http"
	"time"

	gh "github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v62/github"
)

func main() {

	// from github app, replace with real id
	appId := 0773456
	instalationId := 478367040
	org := "ZERO9xz"
	repoName := "hikma-from-github-api-test"

	transport, err := gh.NewKeyFromFile(
		http.DefaultTransport,
		int64(appId),
		int64(instalationId),
		"hikma-test-service.2026-08-12.private-key.pem",
	)
	if err != nil {
		panic(err)
	}

	client := github.NewClient(&http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	})

	// ---- //

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo, resp, err := client.Repositories.Create(
		ctx,
		org,
		&github.Repository{
			Name:    github.String(repoName),
			Private: github.Bool(false),
		})
	if err != nil {
		panic(err)
	}

	_ = repo
	_ = resp
}
