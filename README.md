# `HJ` - HTML to JSON Converter

## Overview
A module and executable that reads HTML and outputs JSON.<br>
Written in Go.

## Usage

### Module
Import the module and call `hj.HTMLtoJSON(string)` to output JSON.
```go
import (
  hj "github.com/HARMONICOM/hj"
)

func main() {
  ...
  json, err := hj.HTMLtoJSON(htmlstring)
  ...
}
```
See `cmd/hj.go` for details.

### Command
The `cmd` directory contains code for execution as a command.<br>
When built and executed, it outputs the input HTML as JSON.
```sh
hj [HTMLFilePath|URL]
```

Read HTML from stdin and convert to JSON
```sh
cat sample.html | hj -
```

Show help
```sh
hj --help
```

Examples:
```sh
hj sample.html
hj https://example.com | jq
cat sample.html | hj -
```

## Output Format
The output format is custom designed to make elements easy to reference.
```html
<html>
    <head>
        <title>Title</title>
    </head>
    <body>
        <h1>Main Heading</h1>
        <div id="content" class="flex">
            <div id="left">
                <h2>Left Heading</h2>
                <p>Left Text</p>
            </div>
            <div id="right">
                <h2>Right Heading</h2>
                <p>Right Text</p>
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
                                "child": "Title"
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
                                "child": "Main Heading"
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
                                                        "child": "Left Heading"
                                                    }
                                                },
                                                {
                                                    "p": {
                                                        "child": "Left Text"
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
                                                        "child": "Right Heading"
                                                    }
                                                },
                                                {
                                                    "p": {
                                                        "child": "Right Text"
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

You can retrieve data using JQ as follows:
```sh
hj sample.html | jq .html.child[0].head.child[0].title.child
"Title"
```

## Build Command
```sh
go build cmd/hj.go
```

You can also use Docker to build Linux commands without environment dependencies. This requires Docker runtime and docker compose.

```sh
make build
make go "build cmd/hj.go"
```
* Make errors can be ignored.

## Other
This program base was created by AI.<br>
The specifications used during creation are stored in the specifications directory.<br>

## License
MIT License

## Copyright
©2025 HARMONICOM.
