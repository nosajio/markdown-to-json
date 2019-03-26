package parse

import (
	"bytes"
	"github.com/gernest/front"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var filePattern = regexp.MustCompile(`([0-9a-z\-]*)-(20[0-9]{2}-[0-9]{2}-[0-9]{2})\.md$`)

type frontmatter struct {
	title string
}

type mdFile struct {
	filename string
	bytes    []byte
}

// Parsed represents a single parsed file
type Parsed struct {
	Date      string `json:"date"`
	Title     string `json:"title"`
	BodyHTML  string `json:"bodyHTML"`
	BodyPlain string `json:"bodyPlain"`
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
		if filePattern.Match([]byte(n)) {
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

func extractYAMLFrontmatter(body []byte) (map[string]interface{}, string, error) {
	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)
	f, b, err := m.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, "", err
	}
	return f, b, nil
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
		meta, body, err := extractYAMLFrontmatter(f.bytes)
		if err != nil {
			log.Printf("Could not extract frontmatter for %s (%s)", f.filename, err.Error())
			continue
		}
		bodyHTML := blackfriday.Run([]byte(body))
		post := Parsed{
			Title:     meta["title"].(string),
			BodyPlain: body,
			BodyHTML:  string(bodyHTML)}
		posts = append(posts, post)
	}
	return posts, nil
}
