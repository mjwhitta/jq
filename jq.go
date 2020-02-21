package jq

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// JSON is a struct that holds a JSON blob.
type JSON struct {
	blob   map[string]interface{}
	escape bool
}

// New is a JSON constructor.
func New(blob string) (*JSON, error) {
	var e error
	var j *JSON

	j = &JSON{blob: map[string]interface{}{}, escape: false}

	if e = json.Unmarshal([]byte(blob), &j.blob); e != nil {
		return j, e
	}

	return j, nil
}

// GetBlob will return the JSON blob as a string.
func (j *JSON) GetBlob() (string, error) {
	var blob = &strings.Builder{}
	var e error
	var enc *json.Encoder = json.NewEncoder(blob)

	enc.SetEscapeHTML(j.escape)

	if e = enc.Encode(j.blob); e != nil {
		return "", e
	}

	return strings.TrimSpace(blob.String()), nil
}

// GetBlobIndent will return the JSON blob as a string with the
// specified prefix and indentation.
func (j *JSON) GetBlobIndent(
	pre string,
	indent string,
) (string, error) {
	var blob = &strings.Builder{}
	var e error
	var enc *json.Encoder = json.NewEncoder(blob)

	enc.SetEscapeHTML(j.escape)
	enc.SetIndent(pre, indent)

	if e = enc.Encode(j.blob); e != nil {
		return "", e
	}

	return strings.TrimSpace(blob.String()), nil
}

// GetKeys will return a list of valid keys if the specified key
// returns an array or map.
func (j *JSON) GetKeys(key ...interface{}) ([]string, error) {
	var e error
	var keys []string
	var val interface{}

	if val, e = j.nestedGetKey(key); e != nil {
		return keys, e
	}

	switch val.(type) {
	case []interface{}:
		for i := 0; i < len(val.([]interface{})); i++ {
			keys = append(keys, strconv.Itoa(i))
		}
		return keys, nil
	case map[string]interface{}:
		for k := range val.(map[string]interface{}) {
			keys = append(keys, k)
		}
		return keys, nil
	default:
		return keys, fmt.Errorf("Key %v has no valid keys", key)
	}
}

// HasKey will return true if the JSON blob has the specified key,
// false otherwise.
func (j *JSON) HasKey(key string) bool {
	var hasKey bool
	_, hasKey = j.blob[key]
	return hasKey
}

func (j *JSON) nestedGetKey(keys []interface{}) (interface{}, error) {
	var e error
	var tryInt int
	var tryString string
	var v interface{}
	var val interface{} = j.blob

	for _, key := range keys {
		if tryString, e = asString(key); e == nil {
			v = val.(map[string]interface{})[tryString]
		} else if tryInt, e = asInt(key); e == nil {
			v = val.([]interface{})[tryInt]
		}

		if (e != nil) || (v == nil) {
			return nil, fmt.Errorf("key %v not found", keys)
		}

		val = v
	}

	return val, nil
}

// Set will set the specified value for the specified key in the JSON
// blob.
func (j *JSON) Set(key string, value interface{}) {
	j.blob[key] = value
}

// SetBlob will replace the underlying map[string]interface{} with a
// new JSON blob.
func (j *JSON) SetBlob(blob string) error {
	var e error

	if e = json.Unmarshal([]byte(blob), &j.blob); e != nil {
		return e
	}

	return nil
}

// SetEscapeHTML will set whether or not Marshalling should escape
// HTML special characters.
func (j *JSON) SetEscapeHTML(escape bool) {
	j.escape = escape
}
