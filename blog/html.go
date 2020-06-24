package blog

import (
	"fmt"
	"html"
	"net/http"
	"path/filepath"
	"strings"
)

func (blog *Blog) ServeHTTP() error {
	render := &Render{
		Blog: blog,
	}

	for _, page := range blog.Pages {
		page := page
		http.HandleFunc(page.Path(), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", render.Page(page))
		})

		http.HandleFunc(filepath.Join(page.Path(), "edit"), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", render.EditPage(page))
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

		s.WriteString("<a href='" + filepath.Join(page.Path(), "edit") + "'>")
		s.WriteString("Edit post")
		s.WriteString("</a>")
	})
}

func (render *Render) EditPage(page *Page) string {
	return render.HTML(page.Title, func(s *strings.Builder) {
		render.Nav(s)

		s.WriteString("<form action='todo'>")
		s.WriteString("<input type='text' value='")
		s.WriteString(html.EscapeString(page.Title))
		s.WriteString("'><br>")

		s.WriteString("<textarea cols=100 rows=20>")
		for i, paragraph := range page.Paragraphs {
			s.WriteString(html.EscapeString(paragraph))
			if i+1 < len(page.Paragraphs) {
				s.WriteString("\n\n")
			}
		}
		s.WriteString("</textarea><br>")

		s.WriteString("<input type='submit' value='update'>")
		s.WriteString("</form>")
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
