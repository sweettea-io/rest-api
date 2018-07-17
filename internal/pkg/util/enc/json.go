package enc

import (
  "bytes"
  "encoding/json"
  "io"
  "github.com/sweettea-io/rest-api/internal/pkg/util/typeconvert"
)

// Type representing a JSON map-like object.
type JSON map[string]interface{}

// Get the JSON data as a byte array.
func (data *JSON) AsByteArray() ([]byte, error) {
  return json.Marshal(data)
}

// Get the JSON data as a string.
func (data *JSON) AsString() (string, error) {
  byteArr, err := data.AsByteArray()

  if err != nil {
    return "", err
  }

  return typeconvert.BytesToStr(byteArr), nil
}

// Get the JSON data as a bytes buffer.
func (data *JSON) AsBuffer() (io.Reader, error) {
  byteArr, err := data.AsByteArray()

  if err != nil {
    return nil, err
  }

  return bytes.NewBuffer(byteArr), nil
}

// Get a new copy of the JSON data that has been marshalled and then immediately unmarshaled.
// This helps ensure types (specifically number types) will match when running something
// like `reflect.DeepEqual` during testing.
func (data *JSON) Cycle() *JSON {
  byteData, _ := data.AsByteArray()

  var cycledData JSON
  json.Unmarshal(byteData, &cycledData)

  return &cycledData
}