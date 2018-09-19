package main

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config represents the configuration file of this app
type Config struct {
	Versions []Version `yaml:"version"`
}

// Version represents the documentation branches of README.io
type Version struct {
	Name       string     `yaml:"name"`
	Categories []Category `yaml:"categories"`
}

// Category represents the documentation category of README.io
type Category struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
}

// LoadConfig reads config file and set it into Config structure
func LoadConfig(path string) (Config, error) {
	var cfg Config
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(buf, &cfg)
	return cfg, err
}

// CategoryID returns the category ID from the configuration file
func (c Config) CategoryID(versionName, categoryName string) (string, error) {
	for _, version := range c.Versions {
		if version.Name == versionName {
			for _, category := range version.Categories {
				if category.Name == categoryName {
					return category.ID, nil
				}
			}
		}
	}
	return "", fmt.Errorf("version %s or category %s is invalid", versionName, categoryName)
}
