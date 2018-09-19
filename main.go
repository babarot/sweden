package main

import (
	"flag"
	"os"
)

func main() {
	var (
		category = flag.String("category", "", "Specify category name")
		version  = flag.String("version", "", "Specify version name")
		config   = flag.String("config", "sweden.yaml", "Specify config path")
	)
	flag.Parse()

	if *category == "" {
		panic("no category")
	}

	cfg, err := LoadConfig(*config)
	if err != nil {
		panic(err)
	}

	filepath := flag.Arg(0)
	categoryID, err := cfg.CategoryID(*version, *category)
	if err != nil {
		panic(err)
	}

	md, err := NewMarkdown(filepath, categoryID)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	md.Convert(file)
}
