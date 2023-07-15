package jq

import (
	"encoding/json"
	"strings"

	"github.com/mjwhitta/errors"
)

// JSON is a struct that holds a JSON blob.
type JSON struct {
	blob   map[string]any
	escape bool
}

// New is a JSON constructor.
func New(blob ...string) (j *JSON, e error) {
	j = &JSON{blob: map[string]any{}, escape: false}
	e = j.SetBlob(strings.Join(blob, ""))

	return
}

// Append will append the specified value to the specified key in the
// JSON blob, if it is an array.
func (j *JSON) Append(value any, keys ...any) error {
	var e error
	var parent []any

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
	j.blob = map[string]any{}
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
func (j *JSON) GetKeys(keys ...any) (ret []string) {
	ret, _ = j.MustGetKeys(keys...)
	return
}

// HasKey will return true if the JSON blob has the specified key,
// false otherwise.
func (j *JSON) HasKey(keys ...any) bool {
	var e error

	_, e = j.nestedGetKey(keys)
	return e == nil
}

// MustGetArray will return an array for the specified key(s) as an
// []any.
func (j *JSON) MustGetArray(
	keys ...any,
) (ret []any, e error) {
	var val any

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
	case []any:
		ret = val
	default:
		e = errors.Newf("key %v is not of type []any", keys)
	}

	return
}

// MustGetKeys will return a list of valid keys if the specified key
// returns an array or map.
func (j *JSON) MustGetKeys(
	keys ...any,
) (ret []string, e error) {
	var val any

	if val, e = j.nestedGetKey(keys); e != nil {
		return
	}

	switch val.(type) {
	case map[string]bool, map[string]string, map[string]any:
		ret = mustGetMapKeys(val)
	case map[string]float32, map[string]float64:
		ret = mustGetFloatMapKeys(val)
	case map[string]int, map[string]int8, map[string]int16,
		map[string]int32, map[string]int64:
		ret = mustGetIntMapKeys(val)
	case map[string]uint, map[string]uint8, map[string]uint16,
		map[string]uint32, map[string]uint64:
		ret = mustGetUintMapKeys(val)
	case []bool, []string, []any:
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
// map[string]any.
func (j *JSON) MustGetMap(
	keys ...any,
) (ret map[string]any, e error) {
	var val any

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
	case map[string]any:
		ret = val
	default:
		e = errors.Newf(
			"key %v is not of type map[string]any",
			keys,
		)
	}

	return
}

func (j *JSON) nestedGetKey(keys []any) (any, error) {
	var e error
	var tryInt int
	var tryString string
	var v any
	var val any = j.blob

	for _, key := range keys {
		if tryString, e = asString(keys, key); e == nil {
			v = val.(map[string]any)[tryString]
		} else if tryInt, e = asInt(keys, key); e == nil {
			v = val.([]any)[tryInt]
		}

		if e != nil {
			return nil, e
		}

		if v == nil {
			return nil, errors.Newf("key %v not found", keys)
		}

		val = v
	}

	return val, nil
}

func (j *JSON) replace(value any) error {
	// Replacing whole JSON map
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
	case map[string]any:
		j.blob = value
	default:
		return errors.Newf(
			"value is not of type map[string]any",
		)
	}

	return j.reset()
}

func (j *JSON) reset() error {
	var e error
	var tmp string

	if tmp, e = j.GetBlob(); e != nil {
		return e
	}

	return j.SetBlob(tmp)
}

// Set will set the specified value for the specified key in the JSON
// blob.
func (j *JSON) Set(value any, keys ...any) error {
	var e error
	var parentArr []any
	var parentMap map[string]any
	var tryInt int
	var tryString string

	if len(keys) == 0 {
		return j.replace(value)
	} else if len(keys) == 1 {
		// Replacing top-level key in JSON map
		if tryString, e = asString(keys, keys[0]); e != nil {
			return e
		}

		j.blob[tryString] = value

		return j.reset()
	}

	if _, e = j.nestedGetKey(keys[0 : len(keys)-1]); e != nil {
		return e
	}

	parentMap, e = j.MustGetMap(keys[0 : len(keys)-1]...)
	if e == nil {
		// Replacing key in JSON map
		tryString, e = asString(keys, keys[len(keys)-1])
		if e != nil {
			return e
		}

		parentMap[tryString] = value
		if e = j.Set(parentMap, keys[0:len(keys)-1]...); e != nil {
			return e
		}

		return j.reset()
	}

	parentArr, e = j.MustGetArray(keys[0 : len(keys)-1]...)
	if e == nil {
		// Replacing array in JSON map
		if tryInt, e = asInt(keys, keys[len(keys)-1]); e != nil {
			return e
		}

		parentArr[tryInt] = value
		if e = j.Set(parentArr, keys[0:len(keys)-1]...); e != nil {
			return e
		}

		return j.reset()
	}

	return errors.Newf("key %v not found", keys)
}

// SetBlob will replace the underlying map[string]any with a
// new JSON blob.
func (j *JSON) SetBlob(blob ...string) error {
	var blobStr = strings.TrimSpace(strings.Join(blob, ""))
	var dec *json.Decoder
	var e error

	if blobStr == "" {
		blobStr = "{}"
	}

	j.blob = map[string]any{}

	dec = json.NewDecoder(strings.NewReader(blobStr))
	if e = dec.Decode(&j.blob); e != nil {
		return errors.Newf("failed to decode JSON: %w", e)
	}

	return nil
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
