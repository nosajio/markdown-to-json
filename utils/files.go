package utils

import (
	"os"
)

// DirectoryExists returns true when a directory in the specified location exists
// on disk
func DirectoryExists(dir string) bool {
	_, err := os.Stat(dir)
	return os.IsNotExist(err)
}
