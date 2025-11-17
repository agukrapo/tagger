package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/agukrapo/tagger/git"
	"github.com/agukrapo/tagger/github"
	"github.com/agukrapo/tagger/versions"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// TODO ctx + interrupts
	host, err := env("GITHUB_API_URL")
	if err != nil {
		return err
	}

	ownerRepo, err := env("GITHUB_REPOSITORY")
	if err != nil {
		return err
	}

	chunks := strings.Split(ownerRepo, "/")
	if len(chunks) != 2 {
		return fmt.Errorf("invalid owner/repository %q", ownerRepo)
	}

	token, err := env("GITHUB_TOKEN")
	if err != nil {
		return err
	}

	api := github.New(chunks[0], chunks[1], host, token)

	local, err := git.SetupClient()
	if err != nil {
		return err
	}

	return versions.Process(api, local)
}

func env(name string) (string, error) {
	if out, ok := os.LookupEnv(name); ok {
		return out, nil
	}
	return "", fmt.Errorf("environment variable %s not set", name)
}
