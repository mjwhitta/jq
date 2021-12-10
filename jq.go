package jq

import (
	"encoding/json"
	"strings"

	"gitlab.com/mjwhitta/errors"
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

	if e = j.Set(append(parent, value), keys...); e != nil {
		return e
	}

	return nil
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
		e = errors.Newf("failed to encode JSON: %w", e)
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

// MustGetArray will return an array for the specified key(s) as an
// []interface{}.
func (j *JSON) MustGetArray(
	keys ...interface{},
) (ret []interface{}, e error) {
	var val interface{}

	if val, e = j.nestedGetKey(keys); e != nil {
		return
	}

	switch val := val.(type) {
	case []bool, []string:
		ret = mustGetArrayAsInterface(val)
	case []float32, []float64:
		ret = mustGetFloatArrayAsInterface(val)
	case []int, []int8, []int16, []int32, []int64:
		ret = mustGetIntArrayAsInterface(val)
	case []uint, []uint8, []uint16, []uint32, []uint64:
		ret = mustGetUintArrayAsInterface(val)
	case []interface{}:
		ret = val
	default:
		e = errors.Newf("key %v is not of type []interface{}", keys)
	}

	return
}

// MustGetKeys will return a list of valid keys if the specified key
// returns an array or map.
func (j *JSON) MustGetKeys(
	keys ...interface{},
) (ret []string, e error) {
	var val interface{}

	if val, e = j.nestedGetKey(keys); e != nil {
		return
	}

	switch val.(type) {
	case map[string]bool, map[string]string, map[string]interface{}:
		ret = mustGetMapKeys(val)
	case map[string]float32, map[string]float64:
		ret = mustGetFloatMapKeys(val)
	case map[string]int, map[string]int8, map[string]int16,
		map[string]int32, map[string]int64:
		ret = mustGetIntMapKeys(val)
	case map[string]uint, map[string]uint8, map[string]uint16,
		map[string]uint32, map[string]uint64:
		ret = mustGetUintMapKeys(val)
	case []bool, []string, []interface{}:
		ret = mustGetArrayKeys(val)
	case []float32, []float64:
		ret = mustGetFloatArrayKeys(val)
	case []int, []int8, []int16, []int32, []int64:
		ret = mustGetIntArrayKeys(val)
	case []uint, []uint8, []uint16, []uint32, []uint64:
		ret = mustGetUintArrayKeys(val)
	default:
		e = errors.Newf("key %v has no valid sub-keys", keys)
	}

	return
}

// MustGetMap will return a map for the specified key(s) as a
// map[string]interface{}.
func (j *JSON) MustGetMap(
	keys ...interface{},
) (ret map[string]interface{}, e error) {
	var val interface{}

	if val, e = j.nestedGetKey(keys); e != nil {
		return
	}

	switch val := val.(type) {
	case map[string]bool, map[string]string:
		ret = mustGetMapAsInterface(val)
	case map[string]float32, map[string]float64:
		ret = mustGetFloatMapAsInterface(val)
	case map[string]int, map[string]int8, map[string]int16,
		map[string]int32, map[string]int64:
		ret = mustGetIntMapAsInterface(val)
	case map[string]uint, map[string]uint8, map[string]uint16,
		map[string]uint32, map[string]uint64:
		ret = mustGetUintMapAsInterface(val)
	case map[string]interface{}:
		ret = val
	default:
		e = errors.Newf(
			"key %v is not of type map[string]interface{}",
			keys,
		)
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
			return nil, errors.Newf("key %v not found", keys)
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
	var parentMap map[string]interface{}
	var tryInt int
	var tryString string

	if len(keys) == 0 {
		switch value := value.(type) {
		case map[string]bool, map[string]string:
			j.blob = mustGetMapAsInterface(value)
		case map[string]float32, map[string]float64:
			j.blob = mustGetFloatMapAsInterface(value)
		case map[string]int, map[string]int8, map[string]int16,
			map[string]int32, map[string]int64:
			j.blob = mustGetIntMapAsInterface(value)
		case map[string]uint, map[string]uint8, map[string]uint16,
			map[string]uint32, map[string]uint64:
			j.blob = mustGetUintMapAsInterface(value)
		case map[string]interface{}:
			j.blob = value
		default:
			e = errors.Newf(
				"value is not of type map[string]interface{}",
			)
		}

		return e
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

	return errors.Newf("key %v not found", keys)
}

// SetBlob will replace the underlying map[string]interface{} with a
// new JSON blob.
func (j *JSON) SetBlob(blob ...string) (e error) {
	var blobStr = strings.TrimSpace(strings.Join(blob, ""))
	var dec *json.Decoder

	if blobStr == "" {
		blobStr = "{}"
	}

	j.blob = map[string]interface{}{}

	dec = json.NewDecoder(strings.NewReader(blobStr))
	if e = dec.Decode(&j.blob); e != nil {
		e = errors.Newf("failed to decode JSON: %w", e)
	}

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
