package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
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
}

// NewMarkdown returns Markdown instance
func NewMarkdown(path, categoryID string) (Markdown, error) {
	fp, err := os.Open(path)
	if err != nil {
		return Markdown{}, err
	}
	scanner := bufio.NewScanner(fp)
	scanner.Scan()
	title := scanner.Text()
	if !strings.HasPrefix(title, "#") {
		return Markdown{}, errors.New("invalid title")
	}
	title = strings.TrimPrefix(title, "#")
	title = strings.TrimSpace(title)
	return Markdown{
		path: path,
		metadata: FrontMatter{
			Title:      title,
			CategoryID: categoryID,
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
