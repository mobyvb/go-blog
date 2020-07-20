package blog

import (
	"fmt"
	"net/http"

	"github.com/mobyvb/go-blog/htmlgenerate"
)

func (blog *Blog) ServeHTTP(port string) error {
	render := &Render{
		Blog: blog,
	}

	// TODO separate port for editing
	http.HandleFunc(PagePrefix, func(w http.ResponseWriter, r *http.Request) {
		// TODO this is ugly and probably not reliable make it better
		pagePath := r.URL.EscapedPath()

		// TODO make this more efficient
		for _, page := range blog.Pages {
			// TODO compare page name instead of path
			if pagePath == page.Path() {
				pageString, err := render.RenderPage(page)
				if err != nil {
					fmt.Fprintf(w, "error generating page: %w", err)
					return
				}
				fmt.Fprintf(w, "%s", pageString)
				break
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tocString, err := render.RenderTableOfContents()
		if err != nil {
			fmt.Fprintf(w, "error generating page: %w", err)
			return
		}
		fmt.Fprintf(w, "%s", tocString)
	})

	fmt.Println("Listening on port " + port)
	return http.ListenAndServe(port, nil)
}

type Render struct {
	Blog *Blog
}

func (render *Render) RenderPage(page *Page) (string, error) {
	navList, err := render.GetNavElement()
	if err != nil {
		return "", err
	}

	header, err := htmlgenerate.NewHTMLElement("h1", page.Title, nil, nil)
	if err != nil {
		return "", err
	}

	pageElements := []*htmlgenerate.HTMLElement{navList, header}
	for _, paragraph := range page.Paragraphs {
		newP, err := htmlgenerate.NewHTMLElement("p", paragraph, nil, nil)
		if err != nil {
			return "", err
		}
		pageElements = append(pageElements, newP)
	}

	htmlPage, err := htmlgenerate.NewPage(page.Title, pageElements)
	if err != nil {
		return "", err
	}
	return htmlPage.String(), nil
}

/*
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
*/

func (render *Render) RenderTableOfContents() (string, error) {
	header, err := htmlgenerate.NewHTMLElement("h1", render.Blog.Title, nil, nil)
	if err != nil {
		return "", err
	}

	navList, err := render.GetNavElement()
	if err != nil {
		return "", err
	}

	htmlPage, err := htmlgenerate.NewPage(render.Blog.Title, []*htmlgenerate.HTMLElement{header, navList})
	if err != nil {
		return "", err
	}
	return htmlPage.String(), nil
}

func (render *Render) GetNavElement() (*htmlgenerate.HTMLElement, error) {
	homeLink, err := htmlgenerate.NewLink("Table of Contents", "/")
	if err != nil {
		return nil, err
	}

	navLinks := []*htmlgenerate.HTMLElement{homeLink}
	for _, navPage := range render.Blog.Pages {
		newLink, err := htmlgenerate.NewLink(navPage.Title, navPage.Path())
		if err != nil {
			return nil, err
		}
		navLinks = append(navLinks, newLink)
	}
	return htmlgenerate.NewList(navLinks)
}
