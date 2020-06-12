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

	var pages []*blog.Page

	err := filepath.Walk(*input, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".txt" {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file: %w", err)
		}

		relpath, err := filepath.Rel(*input, path)
		if err != nil {
			return fmt.Errorf("unable to get relative path %q : %q: %w", *input, path, err)
		}

		slug := blog.SlugFromPath(relpath)
		page, err := blog.ParsePageString(slug, string(content))
		if err != nil {
			return fmt.Errorf("unable to parse content %q: %w", path, err)
		}

		pages = append(pages, page)

		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse input: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(*output, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create dir %q: %v\n", *output, err)
	}

	for _, page := range pages {
		html := page.HTML()

		path := filepath.Join(*output, string(page.Slug)+".html")
		err := ioutil.WriteFile(path, []byte(html), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write file %q: %v\n", path, err)
			os.Exit(1)
		}
	}
}
