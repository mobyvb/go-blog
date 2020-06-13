package blog

import (
	"fmt"
	"html"
	"net/http"
	"strings"
)

func (blog *Blog) ServeHTTP() error {
	render := &Render{
		Blog: blog,
	}

	for _, page := range blog.Pages {
		http.HandleFunc("/"+string(page.Slug), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", render.Page(page))
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", render.TableOfContents())
	})

	fmt.Println("Listening on port 8080")
	return http.ListenAndServe(":8080", nil)
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

func (render *Render) TableOfContents() string {
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
	s.WriteString("<a href='/'>Table of Contents</a>")
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
