# jq

## What is this?

This Go module allows you to parse JSON blobs in a "sane" manner.

## How to install

Open a terminal and run the following:

```
$ go get -u gitlab.com/mjwhitta/jq
```

## Usage

```
package main

import (
    "fmt"
    "strings"

    "gitlab.com/mjwhitta/jq"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println(r.(error).Error())
        }
    }()

    var a bool
    var b string
    var c int
    var d []string
    var e map[string]interface{}
    var f string
    var g string
    var h int64
    var i float64

    var err error
    var j *jq.JSON
    var json string

    json = strings.Join(
        []string{
            "{",
            "  \"a\": true,",
            "  \"b\": \"asdf\",",
            "  \"c\": 1234,",
            "  \"d\": [",
            "    \"blah\",",
            "    \"test\"",
            "  ],",
            "  \"e\": {",
            "    \"int\": 0,",
            "    \"float\": 1.2",
            "  }",
            "}",
        },
        "\n",
    )

    // Initialize JSON object
    if j, err = jq.New(json); err != nil {
        panic(err)
    }

    // JSON blob
    if json, err = j.GetBlob(); err != nil {
        panic(err)
    }

    // Top-level keys
    a = j.GetBool("a")
    b = j.GetString("b")
    c = j.GetInt("c")
    d = j.GetStringArray("d")
    e = j.GetMap("e")

    // Nested keys
    f = j.GetString("d", 0)
    g = j.GetString("d", 1)
    h = j.GetInt64("e", "int")
    i = j.GetFloat64("e", "float")

    // Print
    fmt.Println(json)
    fmt.Println(a)
    fmt.Println(b)
    fmt.Println(c)
    fmt.Println(d)
    fmt.Println(e)
    fmt.Println(f)
    fmt.Println(g)
    fmt.Println(h)
    fmt.Println(i)
}
```

## Links

- [Source](https://gitlab.com/mjwhitta/jq)
