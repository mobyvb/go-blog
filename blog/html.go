package blog

import (
	"fmt"
	"html"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (blog *Blog) WriteToDir(output string) error {
	render := &Render{
		Blog: blog,
	}

	for _, page := range blog.Pages {
		html := render.Page(page)
		err := WritePageFile(output, page.Path(), html)
		if err != nil {
			return fmt.Errorf("failed to write page %q: %w", page.Slug, err)
		}
	}

	html := render.TableOfContent()
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

type Render struct {
	Blog *Blog
}

func (render *Render) Page(page *Page) string {
	return render.HTML(page.Title, func(s *strings.Builder) {
		render.Nav(s)

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

func (render *Render) TableOfContent() string {
	return render.HTML(render.Blog.Title, func(s *strings.Builder) {
		s.WriteString("<h1>")
		s.WriteString(html.EscapeString(render.Blog.Title))
		s.WriteString("</h1>")

		render.Nav(s)
	})
}

func (render *Render) Nav(s *strings.Builder) {
	s.WriteString("<ul>")
	defer s.WriteString("</ul>")

	s.WriteString("<li>")
	s.WriteString("<a href='index.html'>Table of Content</a>")
	s.WriteString("</li>")

	for _, page := range render.Blog.Pages {
		s.WriteString("<li>")
		s.WriteString("<a href='" + page.Path() + "'>")
		s.WriteString(html.EscapeString(page.Title))
		s.WriteString("</a>")
		s.WriteString("</li>")
	}
}

func (render *Render) HTML(title string, body func(*strings.Builder)) string {
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
