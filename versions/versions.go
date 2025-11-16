package versions

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	major, minor, patch int
}

func (v Version) String() string {
	var patch string
	if v.patch != 0 {
		patch = fmt.Sprintf(".%d", v.patch)
	}

	var minor string
	if v.minor != 0 || patch != "" {
		minor = fmt.Sprintf(".%d%s", v.minor, patch)
	}

	return fmt.Sprintf("v%d%s", v.major, minor)
}

func (v Version) bump(major, minor, patch bool) Version {
	if major {
		return Version{v.major + 1, 0, 0}
	} else if minor {
		return Version{v.major, v.minor + 1, 0}
	} else if patch {
		return Version{v.major, v.minor, v.patch + 1}
	}

	return v
}

func (v Version) equals(other Version) bool {
	return v.major == other.major && v.minor == other.minor && v.patch == other.patch
}

type Tag string

func (t Tag) Valid() bool {
	_, err := t.asVersion()
	return err == nil
}

func (t Tag) asVersion() (Version, error) {
	if t == "" {
		return Version{}, nil
	}

	if !strings.HasPrefix(string(t), "v") {
		return Version{}, fmt.Errorf("invalid tag %q", t)
	}

	chunks := strings.Split(string(t[1:]), ".")
	if len(chunks) == 0 || len(chunks) > 3 {
		return Version{}, fmt.Errorf("invalid tag %q", t)
	}

	var major, minor, patch int

	for i, chunk := range chunks {
		v, err := strconv.Atoi(chunk)
		if err != nil {
			return Version{}, fmt.Errorf("invalid tag %q", t)
		}

		switch i {
		case 0:
			major = v
		case 1:
			minor = v
		case 2:
			patch = v
		default:
			return Version{}, fmt.Errorf("invalid tag %q", t)
		}
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

type Commit struct {
	sha, message string
}

func NewCommit(sha, message string) *Commit {
	return &Commit{sha, message}
}

func (c *Commit) SHA() string {
	return c.sha
}

func (c *Commit) change() Change {
	chunks := strings.Split(c.message, ":")
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

type provider interface {
	LatestTag() (Tag, error)
	CommitsSince(tag Tag) ([]*Commit, error)
	Push(*Commit, Version) error
}

func Process(local provider, api provider) error {
	tag, err := api.LatestTag()
	if err != nil {
		return err
	}

	version, err := tag.asVersion()
	if err != nil {
		return err
	}

	fmt.Println("Current version: ", version)

	commits, err := api.CommitsSince(tag)
	if err != nil {
		return err
	}

	var major, minor, patch bool
	for _, commit := range commits {
		fmt.Printf("Commit %s %q\n", commit.sha, commit.message)

		switch commit.change() {
		case Breaking:
			major = true
		case Feat:
			minor = true
		case Fix:
			patch = true
		}
	}

	newVersion := version.bump(major, minor, patch)

	if version.equals(newVersion) {
		fmt.Println("No version change")
		return nil
	}

	fmt.Println("New version: ", newVersion)

	lastCommit := commits[len(commits)-1]

	if err := local.Push(lastCommit, newVersion); err != nil {
		return err
	}

	return nil
}
