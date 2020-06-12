package blog

import (
	"fmt"
	"strings"
)

func ParsePageString(slug Slug, s string) (*Page, error) {
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
