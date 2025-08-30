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

// HTMLElementは、新しい仕様に基づいたHTML要素を表現する構造体
type HTMLElement struct {
	TagName    string                 `json:"-"`
	ID         string                 `json:"-"`
	Attributes map[string]string      `json:"attributes,omitempty"`
	Child      interface{}            `json:"child,omitempty"`
}

// JSONOutputは最終的なJSON出力の形式
type JSONOutput map[string]*HTMLElement

// showHelpはヘルプメッセージを表示する
func showHelp() {
	fmt.Println("HTMLをJSONに変換するコマンド hj")
	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  hj [HTMLファイルパス|URL]  - HTMLファイルまたはURLからHTMLを読み取り、JSONに変換")
	fmt.Println("  cat file.html | hj -      - 標準入力からHTMLを読み取り、JSONに変換")
	fmt.Println("  hj --help                 - このヘルプを表示")
	fmt.Println("")
	fmt.Println("例:")
	fmt.Println("  hj index.html")
	fmt.Println("  hj https://example.com")
	fmt.Println("  cat test.html | hj -")
}

// getHTMLはファイル、URL、または標準入力からHTMLを取得する
func getHTML(input string) (string, error) {
	if input == "" {
		showHelp()
		os.Exit(0)
	}

	if input == "-" {
		// 標準入力から読み取り
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("標準入力の読み取りに失敗しました: %v", err)
		}
		return string(data), nil
	}

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		// URLからHTMLを取得
		resp, err := http.Get(input)
		if err != nil {
			return "", fmt.Errorf("URLの取得に失敗しました: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("HTTPエラー: %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("レスポンスの読み取りに失敗しました: %v", err)
		}
		return string(data), nil
	}

	// ファイルから読み取り
	data, err := os.ReadFile(input)
	if err != nil {
		return "", fmt.Errorf("ファイルの読み取りに失敗しました: %v", err)
	}
	return string(data), nil
}

// generateElementKeyは要素名とIDから仕様に従ったキーを生成する
func generateElementKey(tagName, id string) string {
	if id != "" {
		return tagName + "#" + id
	}
	return tagName
}

// parseHTMLToJSONはHTMLを新しい仕様に基づいてJSONに変換する
func parseHTMLToJSON(n *html.Node) interface{} {
	switch n.Type {
	case html.DocumentNode:
		// ドキュメントノードの場合、子ノード（通常はhtml要素）を処理
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

		// 属性を処理
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

		// 子ノードを処理
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

		// childの内容を決定
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

// htmlToJSONはHTMLを新しい仕様に基づいてJSONに変換する
func htmlToJSON(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("HTMLのパースに失敗しました: %v", err)
	}

	// 新しい仕様に基づいたJSON構造を作成
	jsonStructure := parseHTMLToJSON(doc)

	jsonData, err := json.MarshalIndent(jsonStructure, "", "    ")
	if err != nil {
		return "", fmt.Errorf("JSONへの変換に失敗しました: %v", err)
	}

	return string(jsonData), nil
}

func main() {
	args := os.Args[1:]

	// 引数がない場合またはhelpオプションの場合
	if len(args) == 0 {
		showHelp()
		return
	}

	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		showHelp()
		return
	}

	// HTMLを取得
	htmlContent, err := getHTML(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// HTMLをJSONに変換
	jsonOutput, err := htmlToJSON(htmlContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// JSON出力
	fmt.Println(jsonOutput)
}
