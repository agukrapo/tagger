package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/agukrapo/tagger/versions"
)

type Client struct {
	client *http.Client

	owner, repo, host, token string

	assets []string

	debugInfo []string
}

func New(owner, repo, host, token string, assets []string) *Client {
	return &Client{
		client: http.DefaultClient,
		owner:  owner,
		repo:   repo,
		host:   host,
		token:  token,
		assets: assets,
	}
}

func (c *Client) url(path string) string {
	return fmt.Sprintf("%s/repos/%s/%s/%s", c.host, c.owner, c.repo, path)
}

type request struct {
	method     string
	reader     io.Reader
	size       int64
	name, body string
	headers    map[string]string
	url        string
}

type tagsResponse []struct {
	Name string `json:"name"`
}

func (c *Client) LatestTag() (versions.Tag, error) {
	req := &request{
		method: http.MethodGet,
		name:   "tags",
		url:    c.url("tags"),
	}

	var tags tagsResponse
	if err := c.send(req, &tags); err != nil {
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
		SHA  string `json:"sha"`
		Data struct {
			Message string `json:"message"`
		} `json:"commit"`
	} `json:"commits"`
}

func (c *Client) CommitsSince(tag versions.Tag) ([]*versions.Commit, error) {
	req := &request{
		method: http.MethodGet,
		name:   "compare",
		url:    c.url(fmt.Sprintf("compare/%s...HEAD", tag)),
	}

	var payload compareResponse
	if err := c.send(req, &payload); err != nil {
		return nil, err
	}

	out := make([]*versions.Commit, 0, len(payload.Commits))
	for _, commit := range payload.Commits {
		chunks := strings.Split(commit.Data.Message, "\n")
		out = append(out, versions.NewCommit(commit.SHA, strings.TrimSpace(chunks[0])))
	}

	return out, nil
}

type asset struct {
	name string
	data io.ReadCloser
	size int64
}

func (a asset) close() {
	_ = a.data.Close()
}

func (c *Client) Release(version versions.Version, commits []*versions.Commit) error {
	var files []asset
	defer func() {
		for _, f := range files {
			f.close()
		}
	}()

	for _, name := range c.assets {
		file, err := os.Open(filepath.Clean(name))
		if err != nil {
			return err
		}

		stat, err := file.Stat()
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return fmt.Errorf("asset %q is a directory", name)
		}

		files = append(files, asset{stat.Name(), file, stat.Size()})
	}

	uploadURL, err := c.createRelease(version, commits)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := c.uploadAsset(uploadURL, file); err != nil {
			return err
		}
	}

	return nil
}

type releaseResponse struct {
	UploadURL string `json:"upload_url"`
}

func (c *Client) createRelease(version versions.Version, commits []*versions.Commit) (string, error) {
	body := fmt.Sprintf(`{"tag_name":%q,"name":%q,"body":%q}`, version, version, c.changeLog(commits))

	req := &request{
		method: http.MethodPost,
		reader: strings.NewReader(body),
		name:   "releases",
		body:   body,
		url:    c.url("releases"),
	}

	var out releaseResponse
	return out.UploadURL, c.send(req, &out)
}

func (c *Client) changeLog(commits []*versions.Commit) string {
	var (
		breaking string
		feat     string
		fix      string
		other    string
	)

	appendTo := func(section, title, msg, sha string) string {
		if section == "" {
			section = fmt.Sprintf("#### %s:\n", title)
		}
		url := fmt.Sprintf("https://github.com/%s/%s/commit/%s", c.owner, c.repo, sha)
		return section + fmt.Sprintf("- [%s](%s)\n", msg, url)
	}

	for _, commit := range commits {
		change, msg := commit.Change()
		switch change {
		case versions.Breaking:
			breaking = appendTo(breaking, "Breaking changes", msg, commit.SHA())
		case versions.Feat:
			feat = appendTo(feat, "New features", msg, commit.SHA())
		case versions.Fix:
			fix = appendTo(fix, "Bug fixes", msg, commit.SHA())
		case versions.None:
			other = appendTo(other, "Other", msg, commit.SHA())
		}
	}

	return breaking + feat + fix + other
}

func (c *Client) uploadAsset(url string, file asset) error {
	url = strings.Replace(url, "{?name,label}", "?name="+file.name, 1)

	req := &request{
		method: http.MethodPost,
		reader: file.data,
		size:   file.size,
		name:   "upload",
		body:   "<binary>",
		url:    url,
		headers: map[string]string{
			"Content-Type": "application/octet-stream",
		},
	}

	return c.send(req, nil)
}

type errorResponse struct {
	Message string `json:"message"`
}

func (c *Client) send(in *request, out any) (err error) {
	defer func() {
		if err != nil {
			fmt.Println("DEBUG info:")
			for _, msg := range c.debugInfo {
				fmt.Println(msg)
			}
		}
	}()

	req, err := http.NewRequest(in.method, in.url, in.reader)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	for k, v := range in.headers {
		req.Header.Set(k, v)
	}

	if in.size > 0 {
		req.ContentLength = in.size
	}

	c.debugInfo = append(c.debugInfo, fmt.Sprintf("%s request: %s %s, %s", in.name, in.method, in.url, in.body))

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	c.debugInfo = append(c.debugInfo, fmt.Sprintf("%s response: %s, %s\n", in.name, res.Status, raw))

	if !strings.HasPrefix(res.Status, "2") {
		var errRes errorResponse
		if err := json.Unmarshal(raw, &errRes); err != nil && len(raw) != 0 {
			return fmt.Errorf("error json.Unmarshal: %w", err)
		}
		return fmt.Errorf("%s failed: %s", in.name, errRes.Message)
	}

	if out != nil {
		if err := json.Unmarshal(raw, &out); err != nil {
			return fmt.Errorf("out json.Unmarshal: %w", err)
		}
	}

	return nil
}
