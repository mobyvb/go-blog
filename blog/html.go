package blog

import (
	"html"
	"strings"
)

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
