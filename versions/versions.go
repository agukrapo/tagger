package versions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	major, minor, patch int
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
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

func (c Commit) change() Change {
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

type provider interface {
	LatestTag() (Tag, error)
	CommitsSince(tag Tag) ([]Commit, error)
	Push(Commit, Version) error
}

func Process(p provider) error {
	tag, err := p.LatestTag()
	if err != nil {
		return err
	}

	version, err := tag.asVersion()
	if err != nil {
		return err
	}

	fmt.Println("current version: ", version)

	commits, err := p.CommitsSince(tag)
	if err != nil {
		return err
	}

	var major, minor, patch bool
	for _, commit := range commits {
		fmt.Printf("commit %q\n", commit)

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
		fmt.Println("no version change")
		return nil
	}

	fmt.Println("new version: ", newVersion)

	lastCommit := commits[len(commits)-1]

	if err := p.Push(lastCommit, newVersion); err != nil {
		return err
	}

	return nil
}
