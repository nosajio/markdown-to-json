package parse

import (
	"github.com/nosajio/markdown-to-json/download"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	const (
		dir  string = "/tmp/posts"
		repo string = "https://github.com/nosajio/writing"
	)

	// Ensure posts will be available in specified location
	download.RepoToDisk(repo, dir)

	t.Run("Files(<dir>)", func(t *testing.T) {
		posts, err := Files(dir)
		if err != nil {
			t.Errorf("Files(%s) failed with an error: %s", dir, err.Error())
		}
		if posts == nil || len(posts) == 0 {
			t.Errorf("Files(%s) returned an empty result. Should be slice of Parsed types", dir)
		}

		// Test an individual post for evidence of successful parsing
		firstPost := posts[0]
		if firstPost.Title == "" || len(firstPost.Title) == 0 {
			t.Errorf("Files(%s) doesn't parse the post Title", dir)
		}
		if strings.Contains(firstPost.BodyHTML, "<p>") == false {
			t.Errorf("Files(%s) doesn't parse HTML in BodyHTML", dir)
		}

		// Test for specific props like "date" "slug" etc
		if firstPost.Slug == "" {
			t.Errorf("Files(%s) first item has an empty slug", dir)
		}
		if firstPost.Date.IsZero() {
			t.Errorf("Files(%s) first item has an empty date", dir)
		}

		// Test the result order (should be chronological new -> old)
		var prevDate time.Time
		for i := range posts {
			p := posts[i]
			if prevDate.IsZero() {
				prevDate = p.Date
				continue
			}
			if prevDate.Before(p.Date) {
				t.Errorf("Files(%s) doesn't sort posts chronologically", dir)
			}
			prevDate = p.Date
		}

		// Test custom markdown tags for images
		imgTagPattern := regexp.MustCompile(`(?im)\%img\[.*\]\(.*\)`)
		// The parser will misparse by assuming img tags are links with %img before them
		imgTagBadParsingPattern := regexp.MustCompile(`(?im)\%img<a`)
		for i := range posts {
			p := posts[i]
			if imgTagPattern.Match([]byte(p.BodyPlain)) && imgTagBadParsingPattern.Match([]byte(p.BodyHTML)) {
				t.Errorf("Files(%s) doesn't parse custom %%img[]() tags into HTML", dir)
			}
		}
	})
}
