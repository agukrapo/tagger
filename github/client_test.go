package github

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/agukrapo/tagger/versions"
)

func TestClient_LatestTag(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write(readFile(t, "test-data/tag-response.json"))
	}))
	defer svr.Close()

	c := Client{
		client: svr.Client(),
		host:   svr.URL,
	}

	got, err := c.LatestTag()
	if err != nil {
		t.Fatalf("LatestTag() error = %v", err)
	}

	want := versions.Tag("v4.1.1")
	if got != want {
		t.Errorf("LatestTag() got = %v, want %v", got, want)
	}
}

func TestClient_CommitsSince(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write(readFile(t, "test-data/compare-response.json"))
	}))
	defer svr.Close()

	c := Client{
		client: svr.Client(),
		host:   svr.URL,
	}

	got, err := c.CommitsSince("v4.1.1")
	if err != nil {
		t.Fatalf("CommitsSince() error = %v", err)
	}

	if len(got) != 15 {
		t.Errorf("LatestTag() len(got) = %v, want 15", got)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()

	out, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("readFile: %v", err)
	}

	return out
}
