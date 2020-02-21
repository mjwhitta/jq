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
    "os"
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
    var keys []string

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
            "    \"aFloat\": 1.2,",
            "    \"anInt\": 0,",
            "    \"more\": {",
            "      \"aFloat32\": 1.2,",
            "      \"anInt64\": 0",
            "    }",
            "  }",
            "}",
        },
        "\n",
    )

    // Initialize JSON object
    if j, err = jq.New(json); err != nil {
        fmt.Println(err.Error())
        os.Exit(0)
    }

    // JSON blob
    if json, err = j.GetBlobIndent("", "  "); err != nil {
        fmt.Println(err.Error())
    }

    // Top-level keys
    if j.HasKey("a") {
        if a, err = j.GetBool("a"); e != nil {
            fmt.Println(err.Error())
        }
    }
    b, _ = j.GetString("b")
    c, _ = j.GetInt("c")
    d, _ = j.GetStringArray("d")
    e, _ = j.GetMap("e")

    // Nested keys
    f, _ = j.GetString("d", 0)
    g, _ = j.GetString("d", 1)
    h, _ = j.GetInt64("e", "anInt")
    i, _ = j.GetFloat64("e", "aFloat")

    // Print
    fmt.Println(json)
    fmt.Printf("a = %v\n", a)
    fmt.Printf("b = %v\n", b)
    fmt.Printf("c = %v\n", c)
    fmt.Printf("d = %v\n", d)
    fmt.Printf("e = %v\n", e)
    fmt.Printf("f = %v\n", f)
    fmt.Printf("g = %v\n", g)
    fmt.Printf("h = %v\n", h)
    fmt.Printf("i = %v\n", i)

    // Get sub-keys
    if keys, err = j.GetKeys("a"); err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println(keys)
    }

    if keys, err = j.GetKeys("d"); err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println(keys)
    }

    if keys, err = j.GetKeys("e"); err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println(keys)
    }

    if keys, err = j.GetKeys("e", "more"); err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println(keys)
    }
}
```

## Links

- [Source](https://gitlab.com/mjwhitta/jq)
