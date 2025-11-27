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

	debugInfo []string
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

type tagsResponse []struct {
	Name string `json:"name"`
}

func (c *Client) LatestTag() (versions.Tag, error) {
	var tags tagsResponse
	if err := c.send(http.MethodGet, "tags", "", &tags); err != nil {
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
	var payload compareResponse
	if err := c.send(http.MethodGet, fmt.Sprintf("compare/%s...HEAD", tag), "", &payload); err != nil {
		return nil, err
	}

	out := make([]*versions.Commit, 0, len(payload.Commits))
	for _, commit := range payload.Commits {
		chunks := strings.Split(commit.Data.Message, "\n")
		out = append(out, versions.NewCommit(commit.Data.Tree.Sha, strings.TrimSpace(chunks[0])))
	}

	return out, nil
}

func (c *Client) Push(commit *versions.Commit, version versions.Version) error {
	if err := c.createTagRef(commit, version); err != nil {
		return err
	}

	return c.createRelease(version)
}

func (c *Client) createTagRef(commit *versions.Commit, version versions.Version) error {
	return c.send(http.MethodPost, "git/refs", fmt.Sprintf(`{"sha":%q,"ref":"refs/tags/%s"}`, commit.SHA(), version), nil)
}

func (c *Client) createRelease(version versions.Version) error {
	body := fmt.Sprintf(`{"tag_name":%q}`, version)

	return c.send(http.MethodPost, "releases", body, nil)
}

type errorResponse struct {
	Message string `json:"message"`
}

func (c *Client) send(method, path, body string, out any) (err error) {
	defer func() {
		if err != nil {
			fmt.Println("DEBUG info:")
			for _, msg := range c.debugInfo {
				fmt.Println(msg)
			}
		}
	}()

	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	url := c.url(path)
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	c.debugInfo = append(c.debugInfo, fmt.Sprintf("%s request: %s %s, %s", method, path, url, body))

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	c.debugInfo = append(c.debugInfo, fmt.Sprintf("%s response: %s, %s\n", path, res.Status, raw))

	if !strings.HasPrefix(res.Status, "2") {
		var errRes errorResponse
		if err := json.Unmarshal(raw, &errRes); err != nil {
			return fmt.Errorf("error json.Unmarshal: %w", err)
		}
		return fmt.Errorf("%s failed: %s", path, errRes.Message)
	}

	if out != nil {
		if err := json.Unmarshal(raw, &out); err != nil {
			return fmt.Errorf("out json.Unmarshal: %w", err)
		}
	}

	return nil
}
