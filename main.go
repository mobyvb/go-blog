package main

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	input := flag.String("content", "./content", "content directory")
	output := flag.String("build", "./build", "build of html")

	flag.Parse()

	var pages []*Page

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

		slug := SlugFromPath(relpath)
		page, err := ParseString(slug, string(content))
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

type Slug string

func SlugFromPath(path string) Slug {
	name := path[:len(path)-len(filepath.Ext(path))]
	return Slug(strings.ToLower(name))
}

type Page struct {
	Slug       Slug
	Title      string
	Paragraphs []string
}

func ParseString(slug Slug, s string) (*Page, error) {
	page := &Page{
		Slug: slug,
	}

	for _, block := range strings.Split(s, "\n\n") {
		switch {
		case strings.HasPrefix(block, "#"):
			title := strings.TrimSpace(strings.TrimPrefix(block, "#"))
			if title == "" {
				return nil, fmt.Errorf("invalid #")
			}

			if page.Title != "" {
				return nil, fmt.Errorf("duplicate title")
			}

			page.Title = title
		default:
			paragraph := strings.TrimSpace(block)
			paragraph = strings.ReplaceAll(paragraph, "\n", " ")
			page.Paragraphs = append(page.Paragraphs, paragraph)
		}
	}

	return page, nil
}

func (page *Page) HTML() string {
	var s strings.Builder
	page.WriteHTMLTo(&s)
	return s.String()
}

func (page *Page) WriteHTMLTo(s *strings.Builder) {
	s.WriteString(`<!DOCTYPE html>`)

	s.WriteString(`<html>`)
	defer s.WriteString(`</html>`)

	(func() {
		s.WriteString(`<head>`)
		defer s.WriteString(`</head>`)

		s.WriteString("<title>")
		s.WriteString(html.EscapeString(page.Title))
		s.WriteString("</title>")
	})()

	(func() {
		s.WriteString(`<body>`)
		defer s.WriteString(`</body>`)

		s.WriteString("<h1>")
		s.WriteString(html.EscapeString(page.Title))
		s.WriteString("</h1>")

		for _, paragraph := range page.Paragraphs {
			s.WriteString("<p>")
			s.WriteString(html.EscapeString(paragraph))
			s.WriteString("</p>")
		}
	})()
}
