package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Markdown represents the markdown document
type Markdown struct {
	w        io.Writer
	path     string
	metadata FrontMatter
}

// FrontMatter represents the essential elements for posting the pages
type FrontMatter struct {
	Title      string `yaml:"title"`
	CategoryID string `yaml:"category"`
	ParentDoc  string `yaml:"parentDoc"`
}

// NewMarkdown returns Markdown instance
func NewMarkdown(path, categoryID string) (Markdown, error) {
	fp, err := os.Open(path)
	if err != nil {
		return Markdown{}, err
	}
	fi, err := fp.Stat()
	if err != nil {
		return Markdown{}, err
	}
	if fi.IsDir() {
		log.Printf("%q is dir", path)
		return Markdown{}, nil
	}
	scanner := bufio.NewScanner(fp)
	scanner.Scan()
	title := scanner.Text()
	if !strings.HasPrefix(title, "#") {
		// return Markdown{}, errors.New("invalid title")
		log.Printf("path %q, category ID %q is skipped", path, categoryID)
		return Markdown{}, nil
	}
	title = strings.TrimPrefix(title, "#")
	title = strings.TrimSpace(title)
	return Markdown{
		path: path,
		metadata: FrontMatter{
			Title:      title,
			CategoryID: categoryID,
			// ParentDoc:  parentDoc,
		},
	}, nil
}

// Convert converts from the common markdown to README.io flavored markdown
func (m Markdown) Convert(w io.Writer) error {
	out, err := yaml.Marshal(m.metadata)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(m.path)
	if err != nil {
		return err
	}
	var body []byte
	body = append(body, []byte("---\n")...)
	body = append(body, out...)
	body = append(body, []byte("---\n")...)
	body = append(body, Render(buf)...)
	m.w = w
	m.w.Write(body)
	return nil
}
