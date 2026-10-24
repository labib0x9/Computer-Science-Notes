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

func AddCollaborator(ctx context.Context, client *github.Client, org, user, repoName string) error {
	invitation, resp, err := client.Repositories.AddCollaborator(
		ctx,
		org,
		repoName,
		user,
		&github.RepositoryAddCollaboratorOptions{},
	)

	if err != nil {
		return err
	}

	_ = invitation
	_ = resp
	return nil
}

func main() {

	// from github app
	appId := 4569285
	instalationId := 153167517
	org := "ZERO9xz"
	// template := "tcp-echo-server-template"
	user := "labib0x9"
	repoName := "test-" + user + "-00"

	// ---- // Authenticate

	client := NewClient(org, int64(appId), int64(instalationId))

	// ---- //	Create Repo

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// err := CreateRepository(ctx, client, org, template, repoName)
	// if err != nil {
	// 	panic(err)
	// }

	// // ---- // Add User as Collaborator
	// err = AddCollaborator(ctx, client, org, user, repoName)
	// if err != nil {
	// 	panic(err)
	// }

	resp, err := client.Repositories.Delete(ctx, org, repoName)
	if err != nil {
		panic(err)
	}
	_ = resp
}
