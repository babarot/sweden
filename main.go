package main

import (
	"flag"
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

	filepath := flag.Arg(0)
	docs, err := loadDocs(filepath, *category, *version)
	if err != nil {
		panic(err)
	}
	for _, doc := range docs {
		err := doc.Generate(*config)
		if err != nil {
			panic(err)
		}
	}
}
