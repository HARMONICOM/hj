# HTMLをJSONにするコマンド `HJ` 仕様書

## 動作仕様

### 入力
Linuxコマンドラインで以下を入力する

1. HTMLファイルを渡す場合
```sh
hj HTMLファイルのパス
```
この場合、HTMLファイルを読み取る。

2. URLを渡す場合
```sh
hj URL
```
この場合、URLにアクセスしてHTMLを取得して利用する。

3. 標準入力を受け取る場合
catコマンドの標準出力を受け取る場合
```sh
cat test.html | hj -
```
この場合、標準入力から入力された文字列を利用する。

4. 何も入力しない場合
```sh
hj
```
この場合、ヘルプを表示して終了する

### オプション
- --help：ヘルプを表示する

### 出力
JSONを標準出力として出力する。

1. まず、HTMLを取得する。
2. 取得したHTMLをパースする
3. パースしたHTMLをJSON化する
4. JSONを標準出力する

HTMLに対してどのようなJSONにするかは以下の通り。

```html
<html>
    <head>
        <title>タイトル</title>
    </head>
    <body>
        <h1>大見出し</h1>
        <div id="content" class="flex">
            <div id="left">
                <h2>左見出し</h2>
                <p>テキスト左</p>
            </div>
            <div id="right">
                <h2>右見出し</h2>
                <p>テキスト右</p>
                <img src="test.png" />
            </div>
        </div>
    </body>
</html>
```
↓
```json
{
    "html": {
        "child": [
            {
                "head": {
                    "child": [
                        {
                            "title": {
                                "child": "タイトル"
                            }
                        }
                    ]
                }
            },
            {
                "body": {
                    "child": [
                        {
                            "h1": {
                                "child": "大見出し"
                            }
                        },
                        {
                            "div#content": {
                                "attributes": {
                                    "class": "flex"
                                },
                                "child": [
                                    {
                                        "div#left": {
                                            "child": [
                                                {
                                                    "h2": {
                                                        "child": "左見出し"
                                                    }
                                                },
                                                {
                                                    "p": {
                                                        "child": "テキスト左"
                                                    }
                                                }

                                            ]
                                        }
                                    },
                                    {
                                        "div#right": {
                                            "child": [
                                                {
                                                    "h2": {
                                                        "child": "右見出し"
                                                    }
                                                },
                                                {
                                                    "p": {
                                                        "child": "テキスト右"
                                                    }
                                                },
                                                {
                                                    "img": {
                                                        "attributes": {
                                                            "src": "test.png"
                                                        }
                                                    }
                                                }
                                            ]
                                        }
                                    }
                                ]
                            }
                        }
                    ]
                }
            }
        ]
    }
}
```
