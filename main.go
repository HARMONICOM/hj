package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

// showHelp displays the help message
func showHelp() {
	fmt.Println("HJ - HTML to JSON converter")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  hj [HTMLfilePath|URL]     - Read HTML from file or URL and convert to JSON")
	fmt.Println("  cat file.html | hj -      - Read HTML from stdin and convert to JSON")
	fmt.Println("  hj --help                 - Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  hj index.html")
	fmt.Println("  hj https://example.com")
	fmt.Println("  cat test.html | hj -")
}

// getHTML retrieves HTML from file, URL, or stdin
func getHTML(input string) (string, error) {
	if input == "" {
		showHelp()
		os.Exit(0)
	}

	if input == "-" {
		// Read from stdin
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read from stdin: %v", err)
		}
		return string(data), nil
	}

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		// Fetch HTML from URL
		resp, err := http.Get(input)
		if err != nil {
			return "", fmt.Errorf("failed to fetch URL: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response: %v", err)
		}
		return string(data), nil
	}

	// Read from file
	data, err := os.ReadFile(input)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	return string(data), nil
}

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

// htmlToJSON converts HTML to JSON based on new specification
func htmlToJSON(htmlContent string) (string, error) {
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

func main() {
	args := os.Args[1:]

	// Show help when no arguments or help option
	if len(args) == 0 {
		showHelp()
		return
	}

	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		showHelp()
		return
	}

	// Get HTML
	htmlContent, err := getHTML(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Convert HTML to JSON
	jsonOutput, err := htmlToJSON(htmlContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output JSON
	fmt.Println(jsonOutput)
}
