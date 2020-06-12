package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"blog.test/blog"
)

func main() {
	input := flag.String("content", "./content", "content directory")
	output := flag.String("build", "./build", "build of html")

	flag.Parse()

	b, err := blog.BlogFromDir(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse blog: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(*output, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create dir %q: %v\n", *output, err)
	}

	for _, page := range b.Pages {
		html := page.HTML()

		path := filepath.Join(*output, string(page.Slug)+".html")
		err := ioutil.WriteFile(path, []byte(html), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write file %q: %v\n", path, err)
			os.Exit(1)
		}
	}
}
