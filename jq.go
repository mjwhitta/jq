package jq

import (
	"encoding/json"
	"fmt"
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
	return e == nil
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
