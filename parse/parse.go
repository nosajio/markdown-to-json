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
	"sort"
	"time"
)

var filePattern = regexp.MustCompile(`([0-9a-z\-]*)-(20[0-9]{2}-[0-9]{2}-[0-9]{2})\.md$`)

type frontmatter struct {
	title string
}

type mdFile struct {
	filename string
	bytes    []byte
	date     time.Time
	slug     string
}

// Parsed represents a single parsed file
type Parsed struct {
	Date      time.Time `json:"date"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	BodyHTML  string    `json:"bodyHTML"`
	BodyPlain string    `json:"bodyPlain"`
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
			// Extract slug and date from filename
			filenameParts := filePattern.FindAllStringSubmatch(n, -1)
			dateStr := filenameParts[0][2]
			slugStr := filenameParts[0][1]
			d, err := time.Parse("2006-01-02", dateStr)
			newFile := mdFile{
				filename: n,
				bytes:    f,
				date:     d,
				slug:     slugStr}
			mdFiles = append(mdFiles, newFile)
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

func sortFilesChronological(f []mdFile) ([]mdFile, error) {
	fSorted := make([]mdFile, len(f))
	copy(fSorted, f)
	sort.Slice(fSorted, func(i, j int) bool { return fSorted[i].date.After(fSorted[j].date) })
	return fSorted, nil
}

// Files parses a directory of markdown files and converts them into Post
// types
func Files(dir string) ([]Parsed, error) {
	posts := []Parsed{}
	// Find post files in specified dir
	postFiles, err := readMDFiles(dir)
	if err != nil {
		return nil, err
	}
	// Sort the files by the date in the title
	postFiles, err = sortFilesChronological(postFiles)
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
			Date:      f.date,
			Slug:      f.slug,
			BodyPlain: body,
			BodyHTML:  string(bodyHTML)}
		posts = append(posts, post)
	}
	return posts, nil
}
