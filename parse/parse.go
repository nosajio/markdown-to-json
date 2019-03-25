package parse

import (
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type mdFile struct {
	filename string
	bytes    []byte
}

// Parsed represents a single parsed file
type Parsed struct {
	date, title, bodyHTML, bodyPlain string
}

func readMDFiles(dir string) ([]mdFile, error) {
	mdFiles := []mdFile{}
	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	files, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, n := range files {
		if strings.HasSuffix(n, ".md") {
			f, err := ioutil.ReadFile(filepath.Join(dir, n))
			if err != nil {
				log.Printf("Cannot read file %s/%s", dir, n)
				continue
			}
			mdFiles = append(mdFiles, mdFile{filename: n, bytes: f})
		}
	}
	return mdFiles, nil
}

// Files parses a directory of markdown files and converts them into Post
// types
func Files(dir string) ([]Parsed, error) {
	posts := []Parsed{}
	postFiles, err := readMDFiles(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range postFiles {
		posts = append(posts, Parsed{bodyHTML: string(blackfriday.Run(f.bytes))})
	}
	return posts, nil
}
