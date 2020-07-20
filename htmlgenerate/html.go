package htmlgenerate

import (
	"fmt"
	"html"
)

// HTMLElement is an arbitrary type for an html tag.
type HTMLElement struct {
	tag                 string
	textContent         string
	attributeValuePairs []AttributeValuePair
	childElements       []*HTMLElement
}

// AttributeValuePair is a type representing an html attribute-value pair.
type AttributeValuePair struct {
	attribute string
	value     string
}

var (
	acceptedTags = map[string]struct{}{
		"html":  struct{}{},
		"head":  struct{}{},
		"title": struct{}{},
		"body":  struct{}{},
		"p":     struct{}{},
		"h1":    struct{}{},
		"h2":    struct{}{},
		"a":     struct{}{},
		"li":    struct{}{},
		"ul":    struct{}{},
	}
	selfClosingTags = map[string]struct{}{
		"br":    struct{}{},
		"hr":    struct{}{},
		"input": struct{}{},
	}
)

// NewHTMLElement creates a new HTML element.
func NewHTMLElement(tag, content string, attributeValuePairs []AttributeValuePair, childElements []*HTMLElement) (*HTMLElement, error) {
	if _, ok := acceptedTags[tag]; !ok {
		return nil, fmt.Errorf("invalid tag: %q", tag)
	}
	return &HTMLElement{
		tag:                 tag,
		textContent:         content,
		attributeValuePairs: attributeValuePairs,
		childElements:       childElements,
	}, nil

}

func (e *HTMLElement) String() (htmlString string) {
	// TODO handle doctype better
	if e.tag == "html" {
		htmlString = "<!DOCTYPE html>"
	}

	htmlString += "<" + e.tag

	for _, attributeValue := range e.attributeValuePairs {
		htmlString += " " + html.EscapeString(attributeValue.attribute) + "="
		htmlString += "\"" + html.EscapeString(attributeValue.value) + "\""
	}
	htmlString += ">"
	if _, ok := selfClosingTags[e.tag]; ok {
		return htmlString
	}

	// TODO support tags inside content (e.g. em, strong)
	htmlString += html.EscapeString(e.textContent)

	// TODO prevent recursion
	for _, child := range e.childElements {
		htmlString += child.String()
	}

	htmlString += "</" + e.tag + ">"
	return htmlString
}

// NewLink creates an a tag with href and content set.
func NewLink(content, link string) (*HTMLElement, error) {
	return NewHTMLElement("a", content, []AttributeValuePair{
		AttributeValuePair{attribute: "href", value: link},
	}, nil)
}

// NewList creates an unordered list.
func NewList(items []*HTMLElement) (*HTMLElement, error) {
	liItems := []*HTMLElement{}
	for _, item := range items {
		li, err := NewHTMLElement("li", "", nil, []*HTMLElement{item})
		if err != nil {
			return nil, err
		}
		liItems = append(liItems, li)
	}
	ul, err := NewHTMLElement("ul", "", nil, liItems)
	if err != nil {
		return nil, err
	}
	return ul, nil
}

// NewPage creates an html page.
func NewPage(title string, bodyElements []*HTMLElement) (*HTMLElement, error) {
	titleElement, err := NewHTMLElement("title", title, nil, nil)
	if err != nil {
		return nil, err
	}
	head, err := NewHTMLElement("head", "", nil, []*HTMLElement{titleElement})
	if err != nil {
		return nil, err
	}
	body, err := NewHTMLElement("body", "", nil, bodyElements)

	htmlElements := []*HTMLElement{head, body}
	return NewHTMLElement("html", "", nil, htmlElements)
}
