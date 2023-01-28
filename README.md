# jq

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?style=for-the-badge&logo=cookiecutter)](https://www.buymeacoffee.com/mjwhitta)

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/jq)](https://goreportcard.com/report/github.com/mjwhitta/jq)

## What is this?

This Go module allows you to parse JSON blobs in a "sane" manner.

## How to install

Open a terminal and run the following:

```
$ go get --ldflags="-s -w" --trimpath -u github.com/mjwhitta/jq
```

## Usage

```
package main

import (
    "fmt"
    "os"
    "strings"

    "github.com/mjwhitta/jq"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println(r.(error).Error())
        }
    }()

    var e error
    var j *jq.JSON
    var json string
    var prettyJSON string
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
    if j, e = jq.New(json); e != nil {
        fmt.Println(e.Error())
        os.Exit(0)
    }

    // JSON blob
    if json, e = j.GetBlob(); e != nil {
        fmt.Println(e.Error())
    }
    prettyJSON = j.String()

    fmt.Println(json)
    fmt.Println(prettyJSON)

    // Error checking
    if j.HasKey("a") {
        if _, e = j.MustGetBool("b"); e != nil {
            fmt.Println(e.Error())
        }
        if _, e = j.MustGetBool("a"); e != nil {
            fmt.Println(e.Error())
        } else {
            fmt.Println("a is a bool")
        }
    }

    // Top-level keys
    fmt.Printf("a = %v\n", j.GetBool("a"))
    fmt.Printf("b = %v\n", j.GetString("b"))
    fmt.Printf("c = %v\n", j.GetInt("c"))
    fmt.Printf("d = %v\n", j.GetStringArray("d"))
    fmt.Printf("e = %v\n", j.GetMap("e"))

    // Nested keys
    fmt.Printf("f = %v\n", j.GetString("d", 0))
    fmt.Printf("g = %v\n", j.GetString("d", 1))
    fmt.Printf("h = %v\n", j.GetInt64("e", "anInt"))
    fmt.Printf("i = %v\n", j.GetFloat64("e", "aFloat"))

    // Get sub-keys
    if keys, e = j.MustGetKeys("a"); e != nil {
        fmt.Println(e.Error())
    } else {
        fmt.Println(keys)
    }

    keys = j.GetKeys("d")
    fmt.Println(keys)

    keys = j.GetKeys("e")
    fmt.Println(keys)

    keys = j.GetKeys("e", "more")
    fmt.Println(keys)

    // Set keys
    if e = j.Set(false, "asdf", "test"); e != nil {
        fmt.Println(e.Error())
    }
    j.Set(false, "asdf")

    j.Set(false, "a")

    if e = j.Set("asdf", "d", "asdf"); e != nil {
        fmt.Println(e.Error())
    }
    j.Set("asdf", "d", 1)

    if e = j.Set(17, "e", 0); e != nil {
        fmt.Println(e.Error())
    }
    j.Set(17, "e", "anInt")
    j.Set(19, "e", "more", "anInt64")

    fmt.Println(j.String())
}
```

## Links

- [Source](https://github.com/mjwhitta/jq)
