package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	assets, closeAll, err := parseAssets()
	if err != nil {
		return err
	}
	defer closeAll()

	api := github.New(chunks[0], chunks[1], host, token, assets)

	local, err := git.SetupClient()
	if err != nil {
		return err
	}

	return versions.Process(api, local, api)
}

func env(name string) (string, error) {
	if out, ok := os.LookupEnv(name); ok {
		return out, nil
	}
	return "", fmt.Errorf("environment variable %s not set", name)
}

func parseAssets() ([]github.Asset, func(), error) {
	assets, err := env("RELEASE_ASSETS")
	if err != nil {
		return nil, func() {}, nil
	}

	var (
		out     []github.Asset
		closers []func() error
	)
	for _, pattern := range strings.Split(assets, "\n") {
		pattern := strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		local, err := filepath.Localize(pattern)
		if err != nil {
			return nil, nil, fmt.Errorf("%q: %w", pattern, err)
		}

		matches, err := filepath.Glob(local)
		if err != nil {
			return nil, nil, err
		}

		count := 0
		for _, path := range matches {
			file, err := os.Open(path) // #nosec G304
			if err != nil {
				return nil, nil, err
			}

			stat, err := file.Stat()
			if err != nil {
				return nil, nil, err
			}

			if stat.IsDir() {
				continue
			}

			count++
			out = append(out, github.NewAsset(stat.Name(), file, stat.Size()))
			closers = append(closers, file.Close)
		}

		if count == 0 {
			fmt.Printf("No assets found in %s\n", pattern)
		}
	}

	return out, func() {
		for _, fn := range closers {
			_ = fn()
		}
	}, nil
}
