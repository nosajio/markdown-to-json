package download

import "testing"

func TestDownload(t *testing.T) {
	t.Run("RepoToDisk(<REPO_URL>,<DIR>)", func(t *testing.T) {
		t.Run("download", func(t *testing.T) {
			postsRepo := "https://github.com/nosajio/writing"
			tmpDir := "/tmp/posts"
			if _, err := RepoToDisk(postsRepo, tmpDir); err != nil {
				t.Errorf("RepoToDisk(%s, %s) caused an error: %s", postsRepo, tmpDir, err)
			}
		})

		t.Run("deletes", func(t *testing.T) {
			postsRepo := "https://github.com/nosajio/writing"
			tmpDir := "/tmp/posts"
			if _, err := RepoToDisk(postsRepo, tmpDir); err != nil {
				t.Errorf("RepoToDisk(%s, %s) caused an error: %s", postsRepo, tmpDir, err)
			}
		})

	})
}
