// app must be installed on the organization
// github/organizations/<org-name>/settings

// does it clones private repos ? not sure, maybe

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	gh "github.com/bradleyfalzon/ghinstallation/v2"
)

func NewToken(org string, appId, instalationId int64) string {
	transport, err := gh.NewKeyFromFile(
		http.DefaultTransport,
		appId,
		instalationId,
		"hikma-test-service.2026-08-12.private-key.pem",
	)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token, err := transport.Token(ctx)
	if err != nil {
		panic(err)
	}
	return token
}

func main() {

	appId := 5454545
	instalationId := 121112121
	org := "ZERO9xz"
	repo := "Submission-labib0x9-tcp-echo-server-template"
	token := NewToken(org, int64(appId), int64(instalationId))
	cloneURL := fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s.git", token, org, repo)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	distPath := filepath.Join(os.TempDir(), repo)
	fmt.Println("DIST:", distPath)

	cmd := exec.CommandContext(ctx, "git", "clone", cloneURL, distPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		return
	}

	fmt.Println(string(out))

}
