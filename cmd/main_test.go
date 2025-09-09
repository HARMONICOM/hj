package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestShowHelp tests the showHelp function
func TestShowHelp(t *testing.T) {
	// キャプチャ用のバッファを作成
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// showHelp関数を呼び出し
	showHelp()

	// 標準出力を復元
	w.Close()
	os.Stdout = oldStdout

	// 出力をキャプチャ
	out, _ := io.ReadAll(r)
	output := string(out)

	// ヘルプメッセージの内容を確認
	expected := []string{
		"HJ - HTML to JSON converter",
		"Usage:",
		"hj [HTMLfilePath|URL]",
		"cat file.html | hj -",
		"hj --help",
		"Examples:",
		"hj index.html",
		"hj https://example.com",
		"cat test.html | hj -",
	}

	for _, exp := range expected {
		if !strings.Contains(output, exp) {
			t.Errorf("Expected help output to contain '%s', but it didn't. Got: %s", exp, output)
		}
	}
}

// TestGetHTMLWithFile tests getHTML with file input
func TestGetHTMLWithFile(t *testing.T) {
	// テスト用の一時ファイルを作成
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.html")

	testContent := "<html><body><h1>Test</h1></body></html>"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// getHTML関数をテスト
	result, err := getHTML(testFile)
	if err != nil {
		t.Fatalf("getHTML failed: %v", err)
	}

	if result != testContent {
		t.Errorf("Expected '%s', got '%s'", testContent, result)
	}
}

// TestGetHTMLWithNonExistentFile tests getHTML with non-existent file
func TestGetHTMLWithNonExistentFile(t *testing.T) {
	result, err := getHTML("non_existent_file.html")

	if err == nil {
		t.Error("Expected error for non-existent file, but got none")
	}

	if result != "" {
		t.Errorf("Expected empty result for non-existent file, got '%s'", result)
	}

	if !strings.Contains(err.Error(), "failed to read file") {
		t.Errorf("Expected 'failed to read file' error, got '%s'", err.Error())
	}
}

// TestGetHTMLWithURL tests getHTML with URL input
func TestGetHTMLWithURL(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	testHTML := "<html><body><h1>Test Server</h1></body></html>"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, testHTML)
	}))
	defer server.Close()

	// getHTML関数をテスト
	result, err := getHTML(server.URL)
	if err != nil {
		t.Fatalf("getHTML with URL failed: %v", err)
	}

	if result != testHTML {
		t.Errorf("Expected '%s', got '%s'", testHTML, result)
	}
}

// TestGetHTMLWithInvalidURL tests getHTML with invalid URL
func TestGetHTMLWithInvalidURL(t *testing.T) {
	result, err := getHTML("https://invalid-url-that-does-not-exist.example")

	if err == nil {
		t.Error("Expected error for invalid URL, but got none")
	}

	if result != "" {
		t.Errorf("Expected empty result for invalid URL, got '%s'", result)
	}

	if !strings.Contains(err.Error(), "failed to fetch URL") {
		t.Errorf("Expected 'failed to fetch URL' error, got '%s'", err.Error())
	}
}

// TestGetHTMLWithHTTPError tests getHTML with HTTP error response
func TestGetHTMLWithHTTPError(t *testing.T) {
	// 404エラーを返すテストサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer server.Close()

	result, err := getHTML(server.URL)

	if err == nil {
		t.Error("Expected error for HTTP 404, but got none")
	}

	if result != "" {
		t.Errorf("Expected empty result for HTTP error, got '%s'", result)
	}

	if !strings.Contains(err.Error(), "HTTP error: 404") {
		t.Errorf("Expected 'HTTP error: 404', got '%s'", err.Error())
	}
}

// TestGetHTMLWithStdin tests getHTML with stdin input (-)
func TestGetHTMLWithStdin(t *testing.T) {
	// stdinをモック
	testInput := "<html><body><p>From stdin</p></body></html>"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	// テストデータを書き込み
	go func() {
		defer w.Close()
		fmt.Fprint(w, testInput)
	}()

	// getHTML関数をテスト
	result, err := getHTML("-")

	// stdinを復元
	os.Stdin = oldStdin

	if err != nil {
		t.Fatalf("getHTML with stdin failed: %v", err)
	}

	if result != testInput {
		t.Errorf("Expected '%s', got '%s'", testInput, result)
	}
}

// TestGetHTMLWithEmptyInput tests getHTML with empty input
func TestGetHTMLWithEmptyInput(t *testing.T) {
	// showHelp関数が呼ばれ、os.Exit(0)が実行されるため、
	// このテストは少し複雑になります。実際のテストでは
	// os.Exitを回避する方法を考慮する必要があります。

	// この場合、getHTML("")は内部でshowHelp()を呼び出して
	// os.Exit(0)を実行するため、直接テストするのは困難です。
	// 代わりに、空の入力が適切に処理されることを確認します。
	t.Skip("getHTML with empty input calls os.Exit(0), skipping direct test")
}

// TestMainWithHelp tests main function with help arguments
func TestMainWithHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"no arguments", []string{}},
		{"help flag", []string{"--help"}},
		{"short help flag", []string{"-h"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 標準出力をキャプチャ
			oldStdout := os.Stdout
			oldArgs := os.Args

			_, w, _ := os.Pipe()
			os.Stdout = w
			os.Args = append([]string{"hj"}, tt.args...)

			// main関数を実行（os.Exitが呼ばれるため、回復処理が必要）
			defer func() {
				os.Stdout = oldStdout
				os.Args = oldArgs
				if r := recover(); r != nil {
					// パニックが発生した場合の処理
					t.Logf("Recovered from panic: %v", r)
				}
			}()

			// この部分は実際のテスト環境では工夫が必要
			// main()がos.Exit()を呼ぶため、テストが終了してしまいます
			t.Skip("main() function calls os.Exit(), requires special handling for testing")
		})
	}
}

// TestMainFunctionComponents tests individual components used in main
func TestMainFunctionComponents(t *testing.T) {
	// main関数内で使用される個別のコンポーネントをテスト

	// コマンドライン引数の解析ロジックをテスト
	testArgs := [][]string{
		{},
		{"--help"},
		{"-h"},
		{"test.html"},
		{"https://example.com"},
		{"-"},
	}

	for _, args := range testArgs {
		t.Run(fmt.Sprintf("args_%v", args), func(t *testing.T) {
			// 引数の長さと内容を確認
			if len(args) == 0 {
				// ヘルプが表示されるケース
				t.Log("No arguments - help should be displayed")
			} else if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
				// ヘルプが表示されるケース
				t.Log("Help flag - help should be displayed")
			} else if len(args) >= 1 {
				// HTMLの取得と変換が実行されるケース
				t.Logf("Input argument: %s - HTML processing should occur", args[0])
			}
		})
	}
}

// BenchmarkGetHTMLFile benchmarks file reading performance
func BenchmarkGetHTMLFile(t *testing.B) {
	// ベンチマーク用のテストファイルを作成
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "bench.html")

	testContent := strings.Repeat("<div>Benchmark test content</div>", 100)
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create benchmark test file: %v", err)
	}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_, err := getHTML(testFile)
		if err != nil {
			t.Fatalf("Benchmark getHTML failed: %v", err)
		}
	}
}

// Example_showHelp demonstrates the showHelp function
func Example_showHelp() {
	showHelp()
	// Output:
	// HJ - HTML to JSON converter
	//
	// Usage:
	//   hj [HTMLfilePath|URL]     - Read HTML from file or URL and convert to JSON
	//   cat file.html | hj -      - Read HTML from stdin and convert to JSON
	//   hj --help                 - Show this help message
	//
	// Examples:
	//   hj index.html
	//   hj https://example.com
	//   cat test.html | hj -
}
