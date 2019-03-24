package main

import (
	"gopkg.in/src-d/go-git.v4"
)

// DownloadPostsToDisk takes a url string pointing to a git repo and it
// checks out the repo, then saves the files to $TMP_DIR
func DownloadPostsToDisk(fromURL string, tmpDir string) (*git.Repository, error) {
	rp, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: fromURL,
	})
	return rp, err
}
