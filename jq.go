package jq

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// JSON is a struct that holds a JSON blob.
type JSON struct {
	blob   map[string]interface{}
	escape bool
}

// New is a JSON constructor.
func New(blob ...string) (j *JSON, e error) {
	j = &JSON{blob: map[string]interface{}{}, escape: false}
	e = j.SetBlob(strings.Join(blob, ""))
	return
}

// Append will append the specified value to the specified key in the
// JSON blob, if it is an array.
func (j *JSON) Append(value interface{}, keys ...interface{}) error {
	var e error
	var parent []interface{}

	if parent, e = j.MustGetArray(keys...); e != nil {
		return e
	}

	return j.Set(append(parent, value), keys...)
}

// Clear will reset the JSON blob to {}.
func (j *JSON) Clear() {
	j.blob = map[string]interface{}{}
}

// GetBlob will return the JSON blob as a string. An indentation
// string and a prefix string are accepted as optionally parameters.
func (j *JSON) GetBlob(params ...string) (ret string, e error) {
	var blob = &strings.Builder{}
	var enc *json.Encoder = json.NewEncoder(blob)
	var indent string
	var pre string

	if len(params) > 0 {
		indent = params[0]
	}
	if len(params) > 1 {
		pre = params[1]
	}

	enc.SetEscapeHTML(j.escape)
	enc.SetIndent(pre, indent)

	if e = enc.Encode(j.blob); e != nil {
		return
	}

	ret = strings.TrimSpace(blob.String())
	return
}

// GetKeys will return a list of valid keys if the specified key
// returns an array or map.
func (j *JSON) GetKeys(keys ...interface{}) (ret []string) {
	ret, _ = j.MustGetKeys(keys...)
	return
}

// HasKey will return true if the JSON blob has the specified key,
// false otherwise.
func (j *JSON) HasKey(keys ...interface{}) bool {
	var e error
	_, e = j.nestedGetKey(keys)
	return (e == nil)
}

// MustGetKeys will return a list of valid keys if the specified key
// returns an array or map.
func (j *JSON) MustGetKeys(
	keys ...interface{},
) (ret []string, e error) {
	var less = func(i, j int) bool {
		return (strings.ToLower(ret[i]) < strings.ToLower(ret[j]))
	}
	var val interface{}

	if val, e = j.nestedGetKey(keys); e != nil {
		return
	}

	switch val.(type) {
	case []interface{}:
		for i := 0; i < len(val.([]interface{})); i++ {
			ret = append(ret, strconv.Itoa(i))
		}
	case map[string]interface{}:
		for k := range val.(map[string]interface{}) {
			ret = append(ret, k)
		}

		if !sort.SliceIsSorted(ret, less) {
			sort.SliceStable(ret, less)
		}
	default:
		e = fmt.Errorf("Key %v has no valid sub-keys", keys)
	}
	return
}

func (j *JSON) nestedGetKey(keys []interface{}) (interface{}, error) {
	var e error
	var tryInt int
	var tryString string
	var v interface{}
	var val interface{} = j.blob

	for _, key := range keys {
		if tryString, e = asString(keys, key); e == nil {
			v = val.(map[string]interface{})[tryString]
		} else if tryInt, e = asInt(keys, key); e == nil {
			v = val.([]interface{})[tryInt]
		}

		if (e != nil) || (v == nil) {
			return nil, fmt.Errorf("Key %v not found", keys)
		}

		val = v
	}

	return val, nil
}

// Set will set the specified value for the specified key in the JSON
// blob.
func (j *JSON) Set(value interface{}, keys ...interface{}) error {
	var e error
	var parentArr []interface{}
	var parentMap = map[string]interface{}{}
	var tryInt int
	var tryString string

	if len(keys) == 0 {
		switch value.(type) {
		case map[string]interface{}:
			j.blob = value.(map[string]interface{})
			return nil
		default:
			return fmt.Errorf("Value is not a map[string]interface{}")
		}
	} else if len(keys) == 1 {
		if tryString, e = asString(keys, keys[0]); e != nil {
			return e
		}

		j.blob[tryString] = value
		return nil
	}

	if _, e = j.nestedGetKey(keys[0 : len(keys)-1]); e != nil {
		return e
	}

	parentMap, e = j.MustGetMap(keys[0 : len(keys)-1]...)
	if e == nil {
		tryString, e = asString(keys, keys[len(keys)-1])
		if e != nil {
			return e
		}

		parentMap[tryString] = value
		return j.Set(parentMap, keys[0:len(keys)-1]...)
	}

	parentArr, e = j.MustGetArray(keys[0 : len(keys)-1]...)
	if e == nil {
		if tryInt, e = asInt(keys, keys[len(keys)-1]); e != nil {
			return e
		}

		parentArr[tryInt] = value
		return j.Set(parentArr, keys[0:len(keys)-1]...)
	}

	return fmt.Errorf("Key %v not found", keys)
}

// SetBlob will replace the underlying map[string]interface{} with a
// new JSON blob.
func (j *JSON) SetBlob(blob ...string) (e error) {
	var blobStr = strings.TrimSpace(strings.Join(blob, ""))
	var dec *json.Decoder

	if len(blobStr) == 0 {
		blobStr = "{}"
	}

	j.blob = map[string]interface{}{}

	dec = json.NewDecoder(strings.NewReader(blobStr))
	e = dec.Decode(&j.blob)

	return
}

// SetEscapeHTML will set whether or not Marshalling should escape
// HTML special characters.
func (j *JSON) SetEscapeHTML(escape bool) {
	j.escape = escape
}

// String will return a string representation of JSON instance.
func (j *JSON) String() (ret string) {
	ret, _ = j.GetBlob("  ")
	return
}
