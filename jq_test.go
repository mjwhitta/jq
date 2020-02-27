package jq_test

import (
	"fmt"
	"strings"
	"testing"

	"gitlab.com/mjwhitta/jq"
)

var json string = strings.Join(
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
		"    \"anInt\": 17,",
		"    \"more\": {",
		"      \"aFloat32\": 1.2,",
		"      \"anInt64\": 19",
		"    }",
		"  }",
		"}",
	},
	"\n",
)

func TestAppend(t *testing.T) {
	var actual string
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if e = j.Append("asdf", "d"); e != nil {
		t.Errorf("got: %s; want: nil", e.Error())
	}

	expected = "[blah test asdf]"
	actual = fmt.Sprintf("%v", j.GetArray("d"))
	if actual != expected {
		t.Errorf("got: %s; want: %s", actual, expected)
	}

	if e = j.Append(2, "d"); e != nil {
		t.Errorf("got: %s; want: nil", e.Error())
	}

	expected = "[blah test asdf 2]"
	actual = fmt.Sprintf("%v", j.GetArray("d"))
	if actual != expected {
		t.Errorf("got: %s; want: %s", actual, expected)
	}
}

func TestBadJSON(t *testing.T) {
	var e error
	var expected string

	expected = "EOF"
	if _, e = jq.New(""); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}
}

func TestGetBlob(t *testing.T) {
	var actual string
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if actual, e = j.GetBlob(); e != nil {
		t.Errorf("got: %s; want: nil", e.Error())
	}

	expected = strings.Join(strings.Fields(json), "")

	if actual != expected {
		t.Errorf("got: %s; want: %s", actual, expected)
	}

	if actual, e = j.GetBlob("  "); e != nil {
		t.Errorf("got: %s; want: nil", e.Error())
	}

	if actual != json {
		t.Errorf("got: %s; want: %s", actual, json)
	}

	if j.String() != json {
		t.Errorf("got: %s; want: %s", actual, json)
	}
}

func TestGetBool(t *testing.T) {
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if !j.GetBool("a") {
		t.Errorf("got: false; want: true")
	}

	if j.GetBool("b") {
		t.Errorf("got: true; want: false")
	}

	expected = "Key [b] is not a bool"
	if _, e = j.MustGetBool("b"); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}
}

func TestGetFloat(t *testing.T) {
	var actual float64
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if actual = j.GetFloat64("e", "aFloat"); actual != 1.2 {
		t.Errorf("got %0.1f; want: 1.2", actual)
	}

	if actual = j.GetFloat64("a"); actual != 0 {
		t.Errorf("got: %0.1f; want: 0", actual)
	}

	expected = "Key [a] is not a float64"
	if _, e = j.MustGetFloat64("a"); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}
}

func TestGetInt(t *testing.T) {
	var actual int
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if actual = j.GetInt("c"); actual != 1234 {
		t.Errorf("got %d; want: 1234", actual)
	}

	if actual = j.GetInt("b"); actual != 0 {
		t.Errorf("got: %d; want: 0", actual)
	}

	expected = "Key [b] is not a int"
	if _, e = j.MustGetInt("b"); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}

	if actual = j.GetInt("e", "anInt"); actual != 17 {
		t.Errorf("got: %d; want: 17", actual)
	}
}

func TestGetMap(t *testing.T) {
	var actual = map[string]interface{}{}
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	expected = fmt.Sprintf(
		"%+v",
		map[string]interface{}{
			"aFloat32": 1.2,
			"anInt64":  19,
		},
	)
	actual = j.GetMap("e", "more")
	if fmt.Sprintf("%+v", actual) != expected {
		t.Errorf("got: %+v; want: %+v", actual, expected)
	}

	if actual = j.GetMap("a"); len(actual) > 0 {
		t.Errorf("got: %+v; want: map[]", actual)
	}

	expected = "Key [a] is not a map[string]interface{}"
	if _, e = j.MustGetMap("a"); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}
}

func TestGetString(t *testing.T) {
	var actual string
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if actual = j.GetString("b"); actual != "asdf" {
		t.Errorf("got: %s; want: asdf", actual)
	}

	if actual = j.GetString("a"); actual != "" {
		t.Errorf("got: %s; want: empty string", actual)
	}

	expected = "Key [a] is not a string"
	if _, e = j.MustGetString("a"); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}

	if actual = j.GetString("d", 0); actual != "blah" {
		t.Errorf("got: %s; want: blah", actual)
	}
}

func TestGetStringArray(t *testing.T) {
	var actual []string
	var expected string
	var e error
	var j *jq.JSON

	j, _ = jq.New(json)

	expected = fmt.Sprintf("%v", []string{"blah", "test"})
	actual = j.GetStringArray("d")
	if fmt.Sprintf("%v", actual) != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}

	if actual = j.GetStringArray("a"); len(actual) > 0 {
		t.Errorf("got: %v; want: []", actual)
	}

	expected = "Key [a] is not a []string"
	if _, e = j.MustGetStringArray("a"); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}
}

func TestGoodJSON(t *testing.T) {
	var e error

	if _, e = jq.New(json); e != nil {
		t.Errorf("got: %s; want: nil", e.Error())
	}
}

func TestHasKey(t *testing.T) {
	var j *jq.JSON

	j, _ = jq.New(json)

	if !j.HasKey("a") {
		t.Errorf("got: false; want: true")
	}

	if j.HasKey("asdf") {
		t.Errorf("got: true; want: false")
	}
}

func TestKeys(t *testing.T) {
	var actual []string
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	actual = j.GetKeys("d")
	expected = fmt.Sprintf("%v", []string{"0", "1"})
	if fmt.Sprintf("%v", actual) != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}

	actual = j.GetKeys("e")
	expected = fmt.Sprintf("%v", []string{"aFloat", "anInt", "more"})
	if fmt.Sprintf("%v", actual) != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}

	if actual = j.GetKeys("a"); len(actual) > 0 {
		t.Errorf("got: %v; want: []", actual)
	}

	expected = "Key [a] has no valid sub-keys"
	if _, e = j.MustGetKeys("a"); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}
}

func TestSet(t *testing.T) {
	var actual string
	var e error
	var expected string
	var j *jq.JSON
	var newMap = map[string]interface{}{}

	j, _ = jq.New(json)

	expected = "asdf"
	j.Set("asdf", "d", 0)
	actual = j.GetString("d", 0)
	if actual != expected {
		t.Errorf("got: %s; want: %s", actual, expected)
	}

	expected = "Key [d asdf] is not a int"
	if e = j.Set("asdf", "d", "asdf"); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}

	expected = "Key [e 0] is not a string"
	if e = j.Set("asdf", "e", 0); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}

	expected = "Key [e asdf] not found"
	if e = j.Set("asdf", "e", "asdf", "blah"); e == nil {
		t.Errorf("got: nil; want: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("got: %s; want: %s", e.Error(), expected)
	}

	newMap = map[string]interface{}{"asdf": "blah", "anInt": 7}

	e = j.Set(newMap)
	if e != nil {
		t.Errorf("got: %s; want: nil", e.Error())
	}

	actual = fmt.Sprintf("%+v", j.GetMap())
	expected = fmt.Sprintf("%+v", newMap)
	if actual != expected {
		t.Errorf("got: %s; want: %s", actual, expected)
	}
}
