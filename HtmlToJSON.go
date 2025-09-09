package hj

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// HTMLElement represents an HTML element
type HTMLElement struct {
	TagName    string                 `json:"-"`
	ID         string                 `json:"-"`
	Attributes map[string]string      `json:"attributes,omitempty"`
	Child      interface{}            `json:"child,omitempty"`
}

// JSONOutput represents the final JSON output format
type JSONOutput map[string]*HTMLElement

// generateElementKey generates a key according to specification from tag name and ID
func generateElementKey(tagName, id string) string {
	if id != "" {
		return tagName + "#" + id
	}
	return tagName
}

// parseHTMLToJSON converts HTML to JSON based on new specification
func parseHTMLToJSON(n *html.Node) interface{} {
	switch n.Type {
	case html.DocumentNode:
		// For document node, process child nodes (usually html element)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				return parseHTMLToJSON(c)
			}
		}
		return nil

	case html.ElementNode:
		element := &HTMLElement{
			TagName: n.Data,
		}

		// Process attributes
		var id string
		if len(n.Attr) > 0 {
			for _, attr := range n.Attr {
				if attr.Key == "id" {
					id = attr.Val
					element.ID = id
				} else {
					if element.Attributes == nil {
						element.Attributes = make(map[string]string)
					}
					element.Attributes[attr.Key] = attr.Val
				}
			}
		}

		// Process child nodes
		var children []interface{}
		var textContent strings.Builder

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				childJSON := parseHTMLToJSON(c)
				if childJSON != nil {
					children = append(children, childJSON)
				}
			} else if c.Type == html.TextNode {
				text := strings.TrimSpace(c.Data)
				if text != "" {
					textContent.WriteString(text)
				}
			}
		}

		// Determine child content
		if len(children) > 0 {
			element.Child = children
		} else if textContent.Len() > 0 {
			element.Child = textContent.String()
		}

		// 結果をマップ形式で返す
		key := generateElementKey(n.Data, id)
		result := make(map[string]*HTMLElement)
		result[key] = element

		return result

	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			return text
		}
		return nil

	default:
		return nil
	}
}

// HtmlToJSON converts HTML to JSON based on new specification
func HtmlToJSON(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Create JSON structure based on new specification
	jsonStructure := parseHTMLToJSON(doc)

	jsonData, err := json.MarshalIndent(jsonStructure, "", "    ")
	if err != nil {
		return "", fmt.Errorf("failed to convert to JSON: %v", err)
	}

	return string(jsonData), nil
}
