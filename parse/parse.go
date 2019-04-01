package parse

import (
	"gopkg.in/russross/blackfriday.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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

func extractYAMLFrontmatter(body []byte) (map[string]string, string, error) {
	frontmatterPattern := regexp.MustCompile(`---\n(.*: .*\n)+---`)
	bodyString := frontmatterPattern.ReplaceAllString(string(body), "")
	frontmatterString := frontmatterPattern.Find(body)
	plainYAMLString := strings.Replace(string(frontmatterString), "---", "", 2)
	parsedYAML := make(map[string]string)
	err := yaml.Unmarshal([]byte(plainYAMLString), &parsedYAML)
	if err != nil {
		return nil, bodyString, err
	}
	return parsedYAML, bodyString, nil
}

func sortFilesChronological(f []mdFile) ([]mdFile, error) {
	fSorted := make([]mdFile, len(f))
	copy(fSorted, f)
	sort.Slice(fSorted, func(i, j int) bool { return fSorted[i].date.After(fSorted[j].date) })
	return fSorted, nil
}

func parseBodyHTML(b []byte) []byte {
	// Custom img tag
	imgTagPattern := regexp.MustCompile(`(?im)\%img(\[.*\])(\(.*\))`)
	b = imgTagPattern.ReplaceAll(b, []byte("<img src=\"$2\" alt=\"$1\" />"))
	// Render standard markdown
	bodyHTML := blackfriday.Run(b)
	return bodyHTML
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
		bodyHTML := parseBodyHTML([]byte(body))
		post := Parsed{
			Title:     meta["title"],
			Date:      f.date,
			Slug:      f.slug,
			BodyPlain: body,
			BodyHTML:  string(bodyHTML)}
		posts = append(posts, post)
	}
	return posts, nil
}
