package github

import (
	"encoding/json"
	"fmt"
	"io"
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

	res, err := c.client.Get(c.url("tags"))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var tags tagResponse
	if err := json.NewDecoder(res.Body).Decode(&tags); err != nil {
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
			Tree    struct {
				Sha string `json:"sha"`
			} `json:"tree"`
		} `json:"commit"`
	} `json:"commits"`
}

func (c *Client) CommitsSince(tag versions.Tag) ([]*versions.Commit, error) {
	res, err := c.client.Get(c.url(fmt.Sprintf("compare/%s...HEAD", tag)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var payload compareResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	out := make([]*versions.Commit, 0, len(payload.Commits))
	for _, commit := range payload.Commits {
		chunks := strings.Split(commit.Data.Message, "\n")
		out = append(out, versions.NewCommit(commit.Data.Tree.Sha, strings.TrimSpace(chunks[0])))
	}

	return out, nil
}

type errorResponse struct {
	Message string `json:"message"`
}

func (c *Client) Push(commit *versions.Commit, version versions.Version) error {
	reqBody := strings.NewReader(fmt.Sprintf(`{"sha":%q,"ref":"refs/heads/%s"}`, commit.SHA(), version))
	req, err := http.NewRequest(http.MethodPost, c.url("git/refs"), reqBody)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	fmt.Printf("Create tag response: %s, %s\n", res.Status, resBody)

	if res.StatusCode != http.StatusCreated {
		var errRes errorResponse
		if err := json.Unmarshal(resBody, &errRes); err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}
		return fmt.Errorf("create tag failed: %s", errRes.Message)
	}

	return nil
}
