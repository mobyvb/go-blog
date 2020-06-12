package blog

import (
	"path/filepath"
	"strings"
)

type Page struct {
	Slug       Slug
	Title      string
	Paragraphs []string
}

type Slug string

func SlugFromPath(path string) Slug {
	name := path[:len(path)-len(filepath.Ext(path))]
	return Slug(strings.ToLower(name))
}
