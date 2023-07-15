package jq_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mjwhitta/jq"
)

var def = map[string]any{
	"a": true,
	"b": "asdf",
	"c": 1234,
	"d": [2]string{"blah", "test"},
	"e": map[string]any{
		"aFloat": 1.2,
		"anInt":  17,
		"more": map[string]any{
			"aFloat32": 1.2,
			"anInt64":  19,
		},
	},
}

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
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	expected = "[blah test asdf]"
	actual = fmt.Sprintf("%v", j.GetArray("d"))
	if actual != expected {
		t.Errorf("\ngot: %s\nwant: %s", actual, expected)
	}

	if e = j.Append(2, "d"); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	expected = "[blah test asdf 2]"
	actual = fmt.Sprintf("%v", j.GetArray("d"))
	if actual != expected {
		t.Errorf("\ngot: %s\nwant: %s", actual, expected)
	}
}

func TestBadJSON(t *testing.T) {
	var e error
	var expected string
	var j *jq.JSON

	expected = strings.Join(
		[]string{
			"jq: failed to decode JSON",
			"invalid character '}' looking for beginning of value",
		},
		": ",
	)
	if _, e = jq.New("}"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}

	expected = "jq: failed to decode JSON: unexpected EOF"
	if _, e = jq.New("{"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}

	j, _ = jq.New()

	expected = strings.Join(
		[]string{
			"jq: failed to decode JSON",
			"invalid character '}' looking for beginning of value",
		},
		": ",
	)
	if e = j.SetBlob("}"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}

	expected = "jq: failed to decode JSON: unexpected EOF"
	if e = j.SetBlob("{"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}
}

func TestGetArray(t *testing.T) {
	var actual []any
	var expected string
	var e error
	var j *jq.JSON

	j, _ = jq.New()
	j.Set(def)

	expected = fmt.Sprintf("%v", []string{"blah", "test"})
	if actual, e = j.MustGetArray("d"); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	} else if fmt.Sprintf("%v", actual) != expected {
		t.Errorf("\ngot: %v\nwant: %v", actual, expected)
	}

	if actual = j.GetArray("a"); len(actual) > 0 {
		t.Errorf("\ngot: %v\nwant: []", actual)
	}

	expected = "jq: key [a] is not of type []any"
	if _, e = j.MustGetArray("a"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}
}

func TestGetBlob(t *testing.T) {
	var actual string
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if actual, e = j.GetBlob(); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	expected = strings.Join(strings.Fields(json), "")

	if actual != expected {
		t.Errorf("\ngot: %s\nwant: %s", actual, expected)
	}

	if actual, e = j.GetBlob("  "); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	if actual != json {
		t.Errorf("\ngot: %s\nwant: %s", actual, json)
	}

	if j.String() != json {
		t.Errorf("\ngot: %s\nwant: %s", actual, json)
	}
}

func TestGetBool(t *testing.T) {
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if !j.GetBool("a") {
		t.Errorf("\ngot: false\nwant: true")
	}

	if j.GetBool("b") {
		t.Errorf("\ngot: true\nwant: false")
	}

	expected = "jq: key [b] is not of type bool"
	if _, e = j.MustGetBool("b"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}
}

func TestGetFloat(t *testing.T) {
	var actual float64
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if actual = j.GetFloat64("e", "aFloat"); actual != 1.2 {
		t.Errorf("got %0.1f\nwant: 1.2", actual)
	}

	if actual = j.GetFloat64("a"); actual != 0 {
		t.Errorf("\ngot: %0.1f\nwant: 0", actual)
	}

	expected = "jq: key [a] is not of type float64"
	if _, e = j.MustGetFloat64("a"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}
}

func TestGetInt(t *testing.T) {
	var actual int
	var actual8 uint8
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if actual = j.GetInt("c"); actual != 1234 {
		t.Errorf("got %d\nwant: 1234", actual)
	}

	// 210 b/c 1234 is much larger than uint8 (210 + 1024 == 1234)
	if actual8 = j.GetUint8("c"); actual8 != 210 {
		t.Errorf("got %d\nwant: 210", actual8)
	}

	if actual = j.GetInt("b"); actual != 0 {
		t.Errorf("\ngot: %d\nwant: 0", actual)
	}

	expected = "jq: key [b] is not of type int"
	if _, e = j.MustGetInt("b"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}

	if actual = j.GetInt("e", "anInt"); actual != 17 {
		t.Errorf("\ngot: %d\nwant: 17", actual)
	}
}

func TestGetMap(t *testing.T) {
	var actual map[string]any
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	expected = fmt.Sprintf(
		"%+v",
		map[string]any{
			"aFloat32": 1.2,
			"anInt64":  19,
		},
	)
	actual = j.GetMap("e", "more")
	if fmt.Sprintf("%+v", actual) != expected {
		t.Errorf("\ngot: %+v\nwant: %+v", actual, expected)
	}

	if actual = j.GetMap("a"); len(actual) > 0 {
		t.Errorf("\ngot: %+v\nwant: map[]", actual)
	}

	expected = "jq: key [a] is not of type map[string]any"
	if _, e = j.MustGetMap("a"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}
}

func TestGetString(t *testing.T) {
	var actual string
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New(json)

	if actual = j.GetString("b"); actual != "asdf" {
		t.Errorf("\ngot: %s\nwant: asdf", actual)
	}

	if actual = j.GetString("a"); actual != "" {
		t.Errorf("\ngot: %s\nwant: empty string", actual)
	}

	expected = "jq: key [a] is not of type string"
	if _, e = j.MustGetString("a"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}

	if actual = j.GetString("d", 0); actual != "blah" {
		t.Errorf("\ngot: %s\nwant: blah", actual)
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
		t.Errorf("\ngot: %v\nwant: %v", actual, expected)
	}

	if actual = j.GetStringArray("a"); len(actual) > 0 {
		t.Errorf("\ngot: %v\nwant: []", actual)
	}

	expected = "jq: key [a] is not of type []string"
	if _, e = j.MustGetStringArray("a"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}
}

func TestGoodJSON(t *testing.T) {
	var e error
	var j *jq.JSON

	if _, e = jq.New(json); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	if _, e = jq.New(strings.Fields(json)...); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	if j, e = jq.New(); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	if e = j.SetBlob(); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}
}

func TestHasKey(t *testing.T) {
	var j *jq.JSON

	j, _ = jq.New(json)

	if !j.HasKey("a") {
		t.Errorf("\ngot: false\nwant: true")
	}

	if j.HasKey("asdf") {
		t.Errorf("\ngot: true\nwant: false")
	}
}

func TestKeys(t *testing.T) {
	var actual []string
	var e error
	var expected string
	var j *jq.JSON

	j, _ = jq.New()
	j.Set(def)

	actual = j.GetKeys("d")
	expected = fmt.Sprintf("%v", []string{"0", "1"})
	if fmt.Sprintf("%v", actual) != expected {
		t.Errorf("\ngot: %v\nwant: %v", actual, expected)
	}

	actual = j.GetKeys("e")
	expected = fmt.Sprintf("%v", []string{"aFloat", "anInt", "more"})
	if fmt.Sprintf("%v", actual) != expected {
		t.Errorf("\ngot: %v\nwant: %v", actual, expected)
	}

	if actual = j.GetKeys("a"); len(actual) > 0 {
		t.Errorf("\ngot: %v\nwant: []", actual)
	}

	expected = "jq: key [a] has no valid sub-keys"
	if _, e = j.MustGetKeys("a"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}
}

func TestSet(t *testing.T) {
	var actual string
	var e error
	var expected string
	var j *jq.JSON
	var newMap map[string]any

	j, _ = jq.New(json)

	expected = "asdf"
	j.Set("asdf", "d", 0)
	actual = j.GetString("d", 0)
	if actual != expected {
		t.Errorf("\ngot: %s\nwant: %s", actual, expected)
	}

	expected = "jq: key [d asdf] is not of type int"
	if e = j.Set("asdf", "d", "asdf"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}

	expected = "jq: key [e 0] is not of type string"
	if e = j.Set("asdf", "e", 0); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}

	expected = "jq: key [e asdf] not found"
	if e = j.Set("asdf", "e", "asdf", "blah"); e == nil {
		t.Errorf("\ngot: nil\nwant: %s", expected)
	} else if e.Error() != expected {
		t.Errorf("\ngot: %s\nwant: %s", e.Error(), expected)
	}

	newMap = map[string]any{"asdf": "blah", "anInt": 7}

	e = j.Set(newMap)
	if e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	actual = fmt.Sprintf("%+v", j.GetMap())
	expected = fmt.Sprintf("%+v", newMap)
	if actual != expected {
		t.Errorf("\ngot: %s\nwant: %s", actual, expected)
	}

	if e = j.SetBlob("{\"asdf\": false}"); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	if e = j.SetBlob("{\"blah\": false}"); e != nil {
		t.Errorf("\ngot: %s\nwant: nil", e.Error())
	}

	actual = fmt.Sprintf("%+v", j.GetMap())
	expected = "map[blah:false]"
	if actual != expected {
		t.Errorf("\ngot: %s\nwant: %s", actual, expected)
	}
}
