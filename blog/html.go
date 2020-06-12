package blog

import (
	"fmt"
	"html"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (blog *Blog) WriteToDir(output string) error {
	for _, page := range blog.Pages {
		html := RenderPageHTML(blog, page)
		err := WritePageFile(output, page.Path(), html)
		if err != nil {
			return fmt.Errorf("failed to write page %q: %w", page.Slug, err)
		}
	}

	html := RenderTableOfContentHTML(blog)
	err := WritePageFile(output, "index.html", html)
	if err != nil {
		return fmt.Errorf("failed to table of content: %w", err)
	}

	return nil
}

func WritePageFile(dir, local string, html string) error {
	path := filepath.Join(dir, local)
	err := ioutil.WriteFile(path, []byte(html), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %q: %w", path, err)
	}
	return nil
}

func RenderHTML(title string, body func(*strings.Builder)) string {
	var s strings.Builder
	(func() {
		s.WriteString(`<!DOCTYPE html>`)

		s.WriteString(`<html>`)
		defer s.WriteString(`</html>`)

		(func() {
			s.WriteString(`<head>`)
			defer s.WriteString(`</head>`)

			s.WriteString("<title>")
			s.WriteString(html.EscapeString(title))
			s.WriteString("</title>")
		})()

		s.WriteString(`<body>`)
		defer s.WriteString(`</body>`)
		body(&s)
	})()

	return s.String()
}

func RenderPageHTML(blog *Blog, page *Page) string {
	return RenderHTML(page.Title, func(s *strings.Builder) {
		RenderNav(blog, s)

		s.WriteString("<h1>")
		s.WriteString(html.EscapeString(page.Title))
		s.WriteString("</h1>")

		for _, paragraph := range page.Paragraphs {
			s.WriteString("<p>")
			s.WriteString(html.EscapeString(paragraph))
			s.WriteString("</p>")
		}
	})
}

func RenderTableOfContentHTML(blog *Blog) string {
	return RenderHTML(blog.Title, func(s *strings.Builder) {
		s.WriteString("<h1>")
		s.WriteString(html.EscapeString(blog.Title))
		s.WriteString("</h1>")

		RenderNav(blog, s)
	})
}

func RenderNav(blog *Blog, s *strings.Builder) {
	s.WriteString("<ul>")
	defer s.WriteString("</ul>")

	s.WriteString("<li>")
	s.WriteString("<a href='index.html'>Table of Content</a>")
	s.WriteString("</li>")

	for _, page := range blog.Pages {
		s.WriteString("<li>")
		s.WriteString("<a href='" + page.Path() + "'>")
		s.WriteString(html.EscapeString(page.Title))
		s.WriteString("</a>")
		s.WriteString("</li>")
	}
}
