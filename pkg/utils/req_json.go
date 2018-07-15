package utils

import (
  "encoding/json"
  "io"
  "bytes"
)

type JSON map[string]interface{}

func (data JSON) AsByteArray() ([]byte, error) {
  return json.Marshal(data)
}

func (data JSON) AsString() (string, error) {
  byteArr, err := data.AsByteArray()

  if err != nil {
    return "", err
  }

  return BytesToStr(byteArr), nil
}

func (data JSON) AsBuffer() (io.Reader, error) {
  byteArr, err := data.AsByteArray()

  if err != nil {
    return nil, err
  }

  return bytes.NewBuffer(byteArr), nil
}