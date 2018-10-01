package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// FrontMatter represents the essential elements for posting the pages
type FrontMatter struct {
	Title      string `yaml:"title"`
	CategoryID string `yaml:"category"`
	ParentDoc  string `yaml:"parentDoc,omitempty"`
}

type doc struct {
	baseDir       string
	category      string
	version       string
	title         string
	parentDir     string
	parentDirName string
	path          string
	name          string
}

func (d *doc) getCategoryID(config string) string {
	cfg, _ := LoadConfig(config)
	for _, version := range cfg.Versions {
		if version.Name == d.version {
			for _, category := range version.Categories {
				if category.Name == d.category {
					return category.ID
				}
			}
		}
	}
	return ""
}

func (d *doc) getParentDoc(config string) string {
	cfg, _ := LoadConfig(config)
	for _, version := range cfg.Versions {
		if version.Name == d.version {
			for _, category := range version.Categories {
				if category.Name == d.category {
					for _, parent := range category.Parents {
						if parent.Name == d.parentDirName {
							return parent.ID
						}
					}
				}
			}
		}
	}
	return ""
}

func loadDocs(path, category, version string) ([]*doc, error) {
	var docs []*doc
	file, err := os.Open(path)
	if err != nil {
		return docs, err
	}
	fi, err := file.Stat()
	if err != nil {
		return docs, err
	}
	if fi.IsDir() {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			docs = append(docs, &doc{
				baseDir:       filepath.Dir(filepath.Dir(path)),
				category:      category,
				version:       version,
				parentDir:     filepath.Dir(path),
				parentDirName: filepath.Base(filepath.Dir(path)),
				path:          path,
				name:          filepath.Base(path),
			})
			return nil
		})
		if err != nil {
			return docs, err
		}
	} else {
		docs = append(docs, &doc{
			baseDir:  filepath.Dir(path),
			category: category,
			version:  version,
			path:     path,
			name:     filepath.Base(path),
		})
	}
	return docs, nil
}

func (d *doc) Generate(config string) error {
	file, err := os.Open(d.path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	title := scanner.Text()
	if !strings.HasPrefix(title, "#") {
		log.Println("Skip")
		return nil
	}
	title = strings.TrimPrefix(title, "#")
	title = strings.TrimSpace(title)
	matter, err := yaml.Marshal(FrontMatter{
		Title:      title,
		CategoryID: d.getCategoryID(config),
		ParentDoc:  d.getParentDoc(config),
	})
	if err != nil {
		return nil
	}
	var body []byte
	body = append(body, []byte("---\n")...)
	body = append(body, matter...)
	body = append(body, []byte("---\n")...)
	var contents []byte
	for scanner.Scan() {
		contents = append(contents, scanner.Bytes()...)
		contents = append(contents, []byte("\n")...)
	}
	body = append(body, Render(contents)...)
	newfile, err := os.OpenFile(filepath.Join(d.baseDir, d.name), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer newfile.Close()
	newfile.Write(body)
	return nil
}
