package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/agukrapo/tagger/versions"
)

type Client struct {
	client *http.Client

	owner, repo, host, token string
}

func New(owner, repo, host, token string) *Client {
	return &Client{
		client: http.DefaultClient,
		owner:  owner,
		repo:   repo,
		host:   host,
		token:  token,
	}
}

func (c *Client) url(path string) string {
	return fmt.Sprintf("%s/repos/%s/%s/%s", c.host, c.owner, c.repo, path)
}

type tagResponse []struct {
	Name string `json:"name"`
}

func (c *Client) LatestTag() (versions.Tag, error) {
	// FIXME use token for private repos

	resp, err := c.client.Get(c.url("tags"))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tags tagResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}

	for _, t := range tags {
		if tag := versions.Tag(t.Name); tag.Valid() {
			return tag, nil
		}
	}

	return "", nil
}

type compareResponse struct {
	Commits []struct {
		Data struct {
			Message string `json:"message"`
		} `json:"commit"`
	} `json:"commits"`
}

func (c *Client) CommitsSince(tag versions.Tag) ([]versions.Commit, error) {
	resp, err := c.client.Get(c.url(fmt.Sprintf("compare/%s...HEAD", tag)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload compareResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	out := make([]versions.Commit, 0, len(payload.Commits))
	for _, commit := range payload.Commits {
		chunks := strings.Split(commit.Data.Message, "\n")
		out = append(out, versions.Commit(strings.TrimSpace(chunks[0])))
	}

	return out, nil
}

func (c *Client) Push(version versions.Version) error {
	// TODO call api

	return nil
}
