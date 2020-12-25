package htmlgenerate2

import (
	"html"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// TODO can this code be generated?
// maybe a new package that does code generation

var (
	// TODO configure this in main.go
	templateDir = "templates"
)

type TOCItem struct {
	link  string
	title string
}

func NewTOCItem(link, title string) TOCItem {
	return TOCItem{
		link:  link,
		title: title,
	}
}

func GeneratePost(title string, paragraphs []string, tocItems []TOCItem) (string, error) {
	postTemplateBytes, err := ioutil.ReadFile(filepath.Join(templateDir, "post.html"))
	if err != nil {
		return "", err
	}
	postTemplate := string(postTemplateBytes)

	tocElements, err := getTOCElements(tocItems)
	if err != nil {
		return "", err
	}
	postHTML := strings.ReplaceAll(postTemplate, "{TABLEOFCONTENTS}", tocElements)

	postHTML = strings.ReplaceAll(postHTML, "{TITLE}", html.EscapeString(title))

	paragraphElements := ""
	for _, p := range paragraphs {
		paragraphElements += "<p>"
		paragraphElements += html.EscapeString(p)
		paragraphElements += "</p>"
	}

	postHTML = strings.ReplaceAll(postHTML, "{POSTBODY}", paragraphElements)

	return postHTML, nil
}

func GenerateTableOfContents(title string, tocItems []TOCItem) (string, error) {
	tocTemplateBytes, err := ioutil.ReadFile(filepath.Join(templateDir, "toc.html"))
	if err != nil {
		return "", err
	}
	tocTemplate := string(tocTemplateBytes)

	tocElements, err := getTOCElements(tocItems)
	if err != nil {
		return "", err
	}
	postHTML := strings.ReplaceAll(tocTemplate, "{TABLEOFCONTENTS}", tocElements)

	postHTML = strings.ReplaceAll(postHTML, "{TITLE}", html.EscapeString(title))

	return postHTML, nil
}

func getTOCElements(tocItems []TOCItem) (string, error) {
	tocEntryBytes, err := ioutil.ReadFile(filepath.Join(templateDir, "toclink.html"))
	if err != nil {
		return "", err
	}
	tocEntry := string(tocEntryBytes)
	tocElements := ""
	for _, item := range tocItems {
		newTOCElement := strings.ReplaceAll(tocEntry, "{LINK}", html.EscapeString(item.link))
		newTOCElement = strings.ReplaceAll(newTOCElement, "{TITLE}", html.EscapeString(item.title))
		tocElements += newTOCElement
	}
	return tocElements, nil

}
