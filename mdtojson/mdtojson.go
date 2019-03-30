package mdtojson

import (
	"encoding/json"
	"fmt"
	"github.com/nosajio/markdown-to-json/download"
	"github.com/nosajio/markdown-to-json/parse"
)

// ProcessRepo clones and parses a directory of markdown files with YAML
// frontmatter and returns a JSON string
func ProcessRepo(repoURL string, tmpDIR string) (string, error) {
	if repoURL == "" || tmpDIR == "" {
		return "", fmt.Errorf("repoURL and tmpDIR are required")
	}
	_, err := download.RepoToDisk(repoURL, tmpDIR)
	if err != nil {
		return "", err
	}

	md, err := parse.Files(tmpDIR)
	if err != nil {
		return "", err
	}

	j, err := json.Marshal(md)
	if err != nil {
		return "", err
	}

	return string(j), nil
}
