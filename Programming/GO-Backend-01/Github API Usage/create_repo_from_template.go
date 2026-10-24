// app must be installed on the organization
// template repository te `Template repository` select krte hobe

package main

import (
	"context"
	"net/http"
	"time"

	gh "github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v62/github"
)

func main() {

	// from github app
	appId := 0773456
	instalationId := 478367040
	org := "ZERO9xz"
	repoName := "test-02"

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

	repo, resp, err := client.Repositories.CreateFromTemplate(
		ctx,
		org,
		"tcp-echo-server-template",
		&github.TemplateRepoRequest{
			Name:    github.String(repoName),
			Private: github.Bool(false),
			Owner:   github.String(org),
		})
	if err != nil {
		panic(err)
	}

	_ = repo
	_ = resp
}
