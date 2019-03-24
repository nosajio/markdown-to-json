package main

import "testing"

func TestDownload(t *testing.T) {
	t.Run("DownloadPostsToDisk(<REPO_URL>)", func(t *testing.T) {
		postsRepo := "https://github.com/nosajio/writing"
		tmpDir := "/tmp/posts/"
		if _, err := DownloadPostsToDisk(postsRepo, tmpDir); err != nil {
			t.Errorf("DownloadPostsToDisk(%s, %s) caused an error %s", postsRepo, tmpDir, err)
		}
	})
}
