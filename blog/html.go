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
		html := page.HTML()

		path := filepath.Join(output, string(page.Slug)+".html")
		err := ioutil.WriteFile(path, []byte(html), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %q: %w", path, err)
		}
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

func (page *Page) HTML() string {
	return RenderHTML(page.Title, func(s *strings.Builder) {
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
