package main

import (
	"fmt"
	"os"

	"github.com/agukrapo/tagger/git"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	tag, err := git.LatestTag()
	if err != nil {
		return err
	}

	version, err := tag.AsVersion()
	if err != nil {
		return err
	}

	fmt.Println("current version: ", version)

	commits, err := git.CommitsSince(tag)
	if err != nil {
		return err
	}

	var major, minor, patch bool
	for commit := range commits {
		switch commit.Change() {
		case git.Breaking:
			major = true
		case git.Feat:
			minor = true
		case git.Fix:
			patch = true
		}
	}

	newVersion := version.Bump(major, minor, patch)

	if version.Equals(newVersion) {
		fmt.Println("no version change")
		return nil
	}

	fmt.Println("next version: ", version)

	fmt.Println("continue?")
	if _, err := fmt.Scanln(); err != nil {
		return err
	}

	if err := git.Push(version); err != nil {
		return err
	}

	return nil
}
