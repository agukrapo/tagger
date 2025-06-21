package git

import (
	"errors"
	"fmt"
	"iter"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Version struct {
	major, minor, patch int
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
}

func (v Version) Bump(major, minor, patch bool) Version {
	if major {
		return Version{v.major + 1, 0, 0}
	} else if minor {
		return Version{v.major, v.minor + 1, 0}
	} else if patch {
		return Version{v.major, v.minor, v.patch + 1}
	}

	return v
}

func (v Version) Equals(other Version) bool {
	return v.major == other.major && v.minor == other.minor && v.patch == other.patch
}

type Tag string

func (t Tag) AsVersion() (Version, error) {
	if t == "" {
		return Version{}, nil
	}

	if !strings.HasPrefix(string(t), "v") {
		return Version{}, fmt.Errorf("invalid tag: %s", t)
	}

	chunks := strings.Split(string(t[1:]), ".")
	if len(chunks) != 3 {
		return Version{}, fmt.Errorf("invalid tag: %s", t)
	}

	major, err := strconv.Atoi(chunks[0])
	if err != nil {
		return Version{}, fmt.Errorf("invalid tag: %s", t)
	}

	minor, err := strconv.Atoi(chunks[1])
	if err != nil {
		return Version{}, fmt.Errorf("invalid tag: %s", t)
	}

	patch, err := strconv.Atoi(chunks[2])
	if err != nil {
		return Version{}, fmt.Errorf("invalid tag: %s", t)
	}

	return Version{major, minor, patch}, nil
}

type Change uint8

const (
	None Change = iota
	Breaking
	Feat
	Fix
)

func (c Change) String() string {
	return [...]string{"none", "breaking", "feat", "fix"}[c]
}

type Commit string

func (c Commit) Change() Change {
	re := regexp.MustCompile(`^(\w+)( \(.+\))? (?P<message>.+)$`)

	matches := re.FindStringSubmatch(string(c))
	if len(matches) == 0 {
		return None
	}

	chunks := strings.Split(matches[re.SubexpIndex("message")], ":")
	if len(chunks) == 1 {
		return None
	}

	if strings.HasSuffix(chunks[0], "!") {
		return Breaking
	}

	if strings.HasPrefix(chunks[0], "feat") {
		return Feat
	}

	if strings.HasPrefix(chunks[0], "fix") {
		return Fix
	}

	return None
}

const noTagErr = "fatal: No names found, cannot describe anything."

func LatestTag() (Tag, error) {
	out, err := command("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		if strings.HasPrefix(err.Error(), noTagErr) {
			return "", nil
		}

		return "", fmt.Errorf("git describe: %w", err)
	}

	return Tag(strings.TrimSpace(out)), nil
}

func CommitsSince(tag Tag) (iter.Seq[Commit], error) {
	args := []string{"log", "--oneline"}

	if tag != "" {
		args = slices.Insert(args, 1, fmt.Sprintf("%s..HEAD", tag))
	}

	out, err := command("git", args...)
	if err != nil {
		return nil, err
	}

	return func(yield func(Commit) bool) {
		for _, line := range strings.Split(out, "\n") {
			if !yield(Commit(line)) {
				return
			}
		}
	}, nil
}

func Push(version Version) error {
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
