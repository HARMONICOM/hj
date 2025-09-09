package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	hj "github.com/HARMONICOM/hj"
)

// showHelp displays the help message
func showHelp() {
	fmt.Println("")
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
  fmt.Println("")
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
	jsonOutput, err := hj.HtmlToJSON(htmlContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output JSON
	fmt.Println(jsonOutput)
}
