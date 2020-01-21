package jq

import "encoding/json"

// JSON is a struct that holds a JSON blob.
type JSON struct {
	blob map[string]interface{}
}

// New is a JSON constructor.
func New(blob string) (*JSON, error) {
	var e error
	var j *JSON

	j = &JSON{blob: map[string]interface{}{}}

	if e = json.Unmarshal([]byte(blob), &j.blob); e != nil {
		return j, e
	}

	return j, nil
}

// GetBlob will return the JSON blob as a string.
func (j *JSON) GetBlob() (string, error) {
	var blob []byte
	var e error

	if blob, e = json.Marshal(j.blob); e != nil {
		return "", e
	}

	return string(blob), nil
}

// GetBlobIndent will return the JSON blob as a string with the
// specified prefix and indentation.
func (j *JSON) GetBlobIndent(
	pre string,
	indent string,
) (string, error) {
	var blob []byte
	var e error

	if blob, e = json.MarshalIndent(j.blob, pre, indent); e != nil {
		return "", e
	}

	return string(blob), nil
}

// Has will return true if the JSON blob has the specified key, false
// otherwise.
func (j *JSON) Has(key string) bool {
	var hasKey bool
	_, hasKey = j.blob[key]
	return hasKey
}

// Set will set the specified value for the specified key in the JSON
// blob.
func (j *JSON) Set(key string, value interface{}) {
	j.blob[key] = value
}

func (j *JSON) SetBlob(blob string) error {
	var e error

	if e = json.Unmarshal([]byte(blob), &j.blob); e != nil {
		return e
	}

	return nil
}
