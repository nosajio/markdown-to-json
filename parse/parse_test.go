package parse

import (
	"github.com/nosajio/markdown-to-json/download"
	"strings"
	"testing"
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

	})
}
