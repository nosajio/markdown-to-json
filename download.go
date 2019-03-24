package main

import (
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path/filepath"
)

// DeletePreviousRepo will check if there's already a repo by the same name
// and delete it if there is. This is to avoid conflicts when the repo is being
// cloned
func DeletePreviousRepo(tmpDir string) error {
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		return nil
	}
	d, err := os.Open(tmpDir)
	if err != nil {
		return err
	}
	defer d.Close()
	files, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range files {
		err := os.RemoveAll(filepath.Join(tmpDir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// DownloadPostsToDisk takes a url string pointing to a git repo and it
// checks out the repo, then saves the files to $TMP_DIR
func DownloadPostsToDisk(fromURL string, tmpDir string) (*git.Repository, error) {
	err := DeletePreviousRepo(tmpDir)
	if err != nil {
		return nil, err
	}
	rp, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: fromURL,
	})
	return rp, err
}
