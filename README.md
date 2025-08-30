# `HJ` - HTML to JSON Converter

## Overview
A command that takes HTML as standard input and outputs JSON as standard output.<br>
Written in Go.

## Usage
Read HTML from file or URL and convert to JSON
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

## Build
```sh
go build -o hj
```

You can also use Docker to build Linux commands without environment dependencies. This requires Docker runtime and docker compose.

```sh
make build
make go "build -o hj"
```
* Make errors can be ignored.

## Other
This program was created by AI.<br>
The specifications used during creation are stored in the specifications directory.<br>

## License
MIT License

## Copyright
©2025 HARMONICOM.
