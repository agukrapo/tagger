package github

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

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

type T struct {
	Url          string `json:"url"`
	HtmlUrl      string `json:"html_url"`
	PermalinkUrl string `json:"permalink_url"`
	DiffUrl      string `json:"diff_url"`
	PatchUrl     string `json:"patch_url"`
	BaseCommit   struct {
		Sha    string `json:"sha"`
		NodeId string `json:"node_id"`
		Commit struct {
			Author struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"author"`
			Committer struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"committer"`
			Message string `json:"message"`
			Tree    struct {
				Sha string `json:"sha"`
				Url string `json:"url"`
			} `json:"tree"`
			Url          string `json:"url"`
			CommentCount int    `json:"comment_count"`
			Verification struct {
				Verified   bool        `json:"verified"`
				Reason     string      `json:"reason"`
				Signature  string      `json:"signature"`
				Payload    string      `json:"payload"`
				VerifiedAt interface{} `json:"verified_at"`
			} `json:"verification"`
		} `json:"commit"`
		Url         string `json:"url"`
		HtmlUrl     string `json:"html_url"`
		CommentsUrl string `json:"comments_url"`
		Author      struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			UserViewType      string `json:"user_view_type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"author"`
		Committer struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			UserViewType      string `json:"user_view_type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"committer"`
		Parents []struct {
			Sha     string `json:"sha"`
			Url     string `json:"url"`
			HtmlUrl string `json:"html_url"`
		} `json:"parents"`
	} `json:"base_commit"`
	MergeBaseCommit struct {
		Sha    string `json:"sha"`
		NodeId string `json:"node_id"`
		Commit struct {
			Author struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"author"`
			Committer struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"committer"`
			Message string `json:"message"`
			Tree    struct {
				Sha string `json:"sha"`
				Url string `json:"url"`
			} `json:"tree"`
			Url          string `json:"url"`
			CommentCount int    `json:"comment_count"`
			Verification struct {
				Verified   bool        `json:"verified"`
				Reason     string      `json:"reason"`
				Signature  string      `json:"signature"`
				Payload    string      `json:"payload"`
				VerifiedAt interface{} `json:"verified_at"`
			} `json:"verification"`
		} `json:"commit"`
		Url         string `json:"url"`
		HtmlUrl     string `json:"html_url"`
		CommentsUrl string `json:"comments_url"`
		Author      struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			UserViewType      string `json:"user_view_type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"author"`
		Committer struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			UserViewType      string `json:"user_view_type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"committer"`
		Parents []struct {
			Sha     string `json:"sha"`
			Url     string `json:"url"`
			HtmlUrl string `json:"html_url"`
		} `json:"parents"`
	} `json:"merge_base_commit"`
	Status       string `json:"status"`
	AheadBy      int    `json:"ahead_by"`
	BehindBy     int    `json:"behind_by"`
	TotalCommits int    `json:"total_commits"`
	Commits      []struct {
		Sha    string `json:"sha"`
		NodeId string `json:"node_id"`
		Commit struct {
			Author struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"author"`
			Committer struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"committer"`
			Message string `json:"message"`
			Tree    struct {
				Sha string `json:"sha"`
				Url string `json:"url"`
			} `json:"tree"`
			Url          string `json:"url"`
			CommentCount int    `json:"comment_count"`
			Verification struct {
				Verified   bool       `json:"verified"`
				Reason     string     `json:"reason"`
				Signature  *string    `json:"signature"`
				Payload    *string    `json:"payload"`
				VerifiedAt *time.Time `json:"verified_at"`
			} `json:"verification"`
		} `json:"commit"`
		Url         string `json:"url"`
		HtmlUrl     string `json:"html_url"`
		CommentsUrl string `json:"comments_url"`
		Author      struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			UserViewType      string `json:"user_view_type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"author"`
		Committer struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			UserViewType      string `json:"user_view_type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"committer"`
		Parents []struct {
			Sha     string `json:"sha"`
			Url     string `json:"url"`
			HtmlUrl string `json:"html_url"`
		} `json:"parents"`
	} `json:"commits"`
	Files []struct {
		Sha         string `json:"sha"`
		Filename    string `json:"filename"`
		Status      string `json:"status"`
		Additions   int    `json:"additions"`
		Deletions   int    `json:"deletions"`
		Changes     int    `json:"changes"`
		BlobUrl     string `json:"blob_url"`
		RawUrl      string `json:"raw_url"`
		ContentsUrl string `json:"contents_url"`
		Patch       string `json:"patch,omitempty"`
	} `json:"files"`
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
