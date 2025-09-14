package hj

import (
	"encoding/json"
	"testing"
)

// TestGenerateElementKey tests the generateElementKey function
func TestGenerateElementKey(t *testing.T) {
	tests := []struct {
		name     string
		tagName  string
		id       string
		expected string
	}{
		{
			name:     "Tag with ID",
			tagName:  "div",
			id:       "main",
			expected: "div#main",
		},
		{
			name:     "Tag without ID",
			tagName:  "p",
			id:       "",
			expected: "p",
		},
		{
			name:     "Empty tag name with ID",
			tagName:  "",
			id:       "test",
			expected: "#test",
		},
		{
			name:     "Empty tag name without ID",
			tagName:  "",
			id:       "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateElementKey(tt.tagName, tt.id)
			if result != tt.expected {
				t.Errorf("generateElementKey(%q, %q) = %q, expected %q",
					tt.tagName, tt.id, result, tt.expected)
			}
		})
	}
}

// TestHTMLtoJSON_SimpleElements tests basic HTML elements
func TestHTMLtoJSON_SimpleElements(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		expectError bool
		description string
	}{
		{
			name:        "Simple div",
			html:        "<div>Hello World</div>",
			expectError: false,
			description: "Basic div element with text content",
		},
		{
			name:        "Div with ID",
			html:        `<div id="main">Content</div>`,
			expectError: false,
			description: "Div element with ID attribute",
		},
		{
			name:        "Multiple attributes",
			html:        `<p class="text" data-value="123">Text</p>`,
			expectError: false,
			description: "Element with multiple attributes",
		},
		{
			name:        "Empty element",
			html:        `<br>`,
			expectError: false,
			description: "Self-closing element",
		},
		{
			name:        "Nested elements",
			html:        `<div><p>Nested</p></div>`,
			expectError: false,
			description: "Nested HTML elements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HTMLtoJSON(tt.html)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !tt.expectError {
				// Validate that the result is valid JSON
				var jsonResult interface{}
				if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
					t.Errorf("Result is not valid JSON: %v\nResult: %s", err, result)
				}

				// Print result for manual inspection
				t.Logf("Test: %s\nHTML: %s\nJSON: %s\n", tt.description, tt.html, result)
			}
		})
	}
}

// TestHTMLtoJSON_ComplexStructures tests complex HTML structures
func TestHTMLtoJSON_ComplexStructures(t *testing.T) {
	tests := []struct {
		name string
		html string
	}{
		{
			name: "Table structure",
			html: `<table>
				<thead>
					<tr>
						<th>Header 1</th>
						<th>Header 2</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>Cell 1</td>
						<td>Cell 2</td>
					</tr>
				</tbody>
			</table>`,
		},
		{
			name: "Form with inputs",
			html: `<form id="myform">
				<input type="text" name="username" placeholder="Username">
				<input type="password" name="password" placeholder="Password">
				<button type="submit">Submit</button>
			</form>`,
		},
		{
			name: "Navigation with list",
			html: `<nav id="navigation">
				<ul>
					<li><a href="#home">Home</a></li>
					<li><a href="#about">About</a></li>
					<li><a href="#contact">Contact</a></li>
				</ul>
			</nav>`,
		},
		{
			name: "Article with multiple sections",
			html: `<article id="article1">
				<header>
					<h1>Article Title</h1>
					<time datetime="2024-01-01">January 1, 2024</time>
				</header>
				<section>
					<p>First paragraph with <strong>bold text</strong> and <em>italic text</em>.</p>
					<p>Second paragraph with a <a href="https://example.com">link</a>.</p>
				</section>
				<footer>
					<p>Article footer</p>
				</footer>
			</article>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HTMLtoJSON(tt.html)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Validate that the result is valid JSON
			var jsonResult interface{}
			if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
				t.Errorf("Result is not valid JSON: %v", err)
				return
			}

			// Print result for inspection
			t.Logf("Complex structure test: %s\nResult: %s\n", tt.name, result)
		})
	}
}

// TestHTMLtoJSON_EdgeCases tests edge cases and potential error conditions
func TestHTMLtoJSON_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		expectError bool
		description string
	}{
		{
			name:        "Empty HTML",
			html:        "",
			expectError: false,
			description: "Empty HTML string",
		},
		{
			name:        "Only whitespace",
			html:        "   \n\t   ",
			expectError: false,
			description: "HTML with only whitespace",
		},
		{
			name:        "Text only",
			html:        "Just plain text",
			expectError: false,
			description: "Plain text without HTML tags",
		},
		{
			name:        "Multiple root elements",
			html:        `<div>First</div><div>Second</div>`,
			expectError: false,
			description: "Multiple root level elements",
		},
		{
			name:        "HTML with comments",
			html:        `<!-- Comment --><div>Content</div>`,
			expectError: false,
			description: "HTML with comments",
		},
		{
			name:        "Special characters",
			html:        `<p>Special chars: &lt; &gt; &amp; &quot; &#39;</p>`,
			expectError: false,
			description: "HTML with encoded special characters",
		},
		{
			name:        "Unicode content",
			html:        `<p>Unicode: こんにちは 🌍 αβγ</p>`,
			expectError: false,
			description: "HTML with Unicode characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HTMLtoJSON(tt.html)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !tt.expectError && result != "" {
				// Validate that the result is valid JSON when not empty
				var jsonResult interface{}
				if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
					t.Errorf("Result is not valid JSON: %v\nResult: %s", err, result)
				}
			}

			t.Logf("Edge case: %s\nHTML: %q\nResult: %s\n", tt.description, tt.html, result)
		})
	}
}

// TestHTMLtoJSON_AttributeHandling tests specific attribute handling
func TestHTMLtoJSON_AttributeHandling(t *testing.T) {
	tests := []struct {
		name string
		html string
	}{
		{
			name: "ID attribute handling",
			html: `<div id="main-content">Content</div>`,
		},
		{
			name: "Class attribute",
			html: `<div class="container main">Content</div>`,
		},
		{
			name: "Data attributes",
			html: `<div data-value="123" data-name="test">Content</div>`,
		},
		{
			name: "Boolean attributes",
			html: `<input type="checkbox" checked disabled>`,
		},
		{
			name: "Mixed attributes with ID",
			html: `<button id="submit-btn" type="submit" class="btn primary" data-action="save">Submit</button>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HTMLtoJSON(tt.html)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Parse the result to verify structure
			var jsonResult map[string]interface{}
			if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
				t.Errorf("Result is not valid JSON: %v", err)
				return
			}

			t.Logf("Attribute test: %s\nHTML: %s\nJSON: %s\n", tt.name, tt.html, result)
		})
	}
}

// BenchmarkHTMLtoJSON benchmarks the HTMLtoJSON function
func BenchmarkHTMLtoJSON(b *testing.B) {
	html := `<article id="article1">
		<header>
			<h1>Benchmark Article</h1>
			<time datetime="2024-01-01">January 1, 2024</time>
		</header>
		<section class="content">
			<p>First paragraph with <strong>bold text</strong>.</p>
			<p>Second paragraph with a <a href="https://example.com">link</a>.</p>
			<ul>
				<li>List item 1</li>
				<li>List item 2</li>
				<li>List item 3</li>
			</ul>
		</section>
	</article>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HTMLtoJSON(html)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
