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

func NewClient(org string, appId, instalationId int64) *github.Client {
	transport, err := gh.NewKeyFromFile(
		http.DefaultTransport,
		appId,
		instalationId,
		"hikma-test-service.2026-08-12.private-key.pem",
	)
	if err != nil {
		panic(err)
	}

	return github.NewClient(&http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	})
}

func CreateRepository(ctx context.Context, client *github.Client, org, template, repo string) error {
	trepo, resp, err := client.Repositories.CreateFromTemplate(
		ctx,
		org,
		template, &github.TemplateRepoRequest{
			Name:    github.String(repo),
			Private: github.Bool(false),
			Owner:   github.String(org),
		},
	)

	if err != nil {
		return err
	}

	_ = trepo
	_ = resp
	return nil
}

func main() {

	// from github app
	appId := 0123456
	instalationId := 123456780
	org := "ZERO9xz"
	template := "tcp-echo-server-template"
	user := "labib0x9"
	repoName := "test-" + user + "-00"

	// ---- // Authenticate

	client := NewClient(org, int64(appId), int64(instalationId))

	// ---- //	Create Repo

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := CreateRepository(ctx, client, org, template, repoName)
	if err != nil {
		panic(err)
	}

	// ---- //

	invitation, resp, err := client.Repositories.AddCollaborator(
		ctx,
		org,
		repoName,
		user,
		&github.RepositoryAddCollaboratorOptions{},
	)

	if err != nil {
		panic(err)
	}

	_ = invitation
	_ = resp
}
