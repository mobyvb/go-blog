package blog

import (
	"path/filepath"
	"strings"
)

const PagePrefix = "/posts/"

type Page struct {
	Slug       Slug
	Title      string
	Paragraphs []string
}

type Slug string

func (page *Page) Path() string {
	return filepath.Join(PagePrefix, string(page.Slug))
}

func SlugFromPath(path string) Slug {
	name := path[:len(path)-len(filepath.Ext(path))]
	return Slug(strings.ToLower(name))
}

// TODO writing page to file/db/whatever storage is being used for raw text content
// TODO reloading page from file/db/whatever storage is being used for raw text content
