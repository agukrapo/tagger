package git

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"

	"github.com/agukrapo/tagger/versions"
)

const noTagErr = "fatal: No names found, cannot describe anything."

type Client struct{}

func (Client) LatestTag() (versions.Tag, error) {
	out, err := command("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		if strings.HasPrefix(err.Error(), noTagErr) {
			return "", nil
		}

		return "", fmt.Errorf("git describe: %w", err)
	}

	return versions.Tag(strings.TrimSpace(out)), nil
}

func (Client) CommitsSince(tag versions.Tag) ([]*versions.Commit, error) {
	args := []string{"log", "--oneline"}

	if tag != "" {
		args = slices.Insert(args, 1, fmt.Sprintf("%s..HEAD", tag))
	}

	commits, err := command("git", args...)
	if err != nil {
		return nil, err
	}

	var out []*versions.Commit
	for _, line := range strings.Split(commits, "\n") {
		if commit, ok := parse(line); ok {
			out = append(out, commit)
		}
	}

	return out, nil
}

func parse(line string) (*versions.Commit, bool) {
	re := regexp.MustCompile(`^(?P<sha>\w+)( \(.+\))? (?P<message>.+)$`)

	matches := re.FindStringSubmatch(line)
	if len(matches) == 0 {
		return nil, false
	}

	return versions.NewCommit(matches[re.SubexpIndex("sha")], matches[re.SubexpIndex("message")]), true
}

func (Client) Push(_ *versions.Commit, version versions.Version) error {
	if _, err := command("git", "config", "user.name", "Tagger"); err != nil {
		return fmt.Errorf("git config: %w", err)
	}

	if _, err := command("git", "config", "user.email", "12501907+agukrapo@users.noreply.github.com"); err != nil {
		return fmt.Errorf("git config: %w", err)
	}

	if _, err := command("git", "tag", version.String()); err != nil {
		return fmt.Errorf("git tag: %w", err)
	}

	if _, err := command("git", "push", "--follow-tags"); err != nil {
		return fmt.Errorf("git push: %w", err)
	}

	return nil
}

func command(in string, arg ...string) (string, error) {
	cmd := exec.Command(in, arg...)

	bytes, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("%s", exitErr.Stderr)
		}

		return "", err
	}

	return string(bytes), nil
}
