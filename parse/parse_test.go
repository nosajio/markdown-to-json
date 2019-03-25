package parse

import (
	"github.com/nosajio/markdown-to-json/download"
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
			t.Errorf("Files(%s) returned an empty result. Should be slice of Post types", dir)
		}

	})
}
