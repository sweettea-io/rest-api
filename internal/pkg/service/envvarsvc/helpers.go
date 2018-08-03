package envvarsvc

import (
  "encoding/json"
  "fmt"
  "strings"
)

func MapFromBytes(data []bytes) (map[string]string, error) {
  var envsMap map[string]string

  // Unmarshal bytes into map.
  if err := json.Unmarshal([]byte(data), &envsMap); err != nil {
    return nil, fmt.Errorf("error unmarshalling env bytes into map: %s", err.Error())
  }

  // Uppercase all the keys.
  upperEnvsMap := map[string]string{}
  for k, v := range envsMap {
    upperEnvsMap[strings.ToUpper(k)] = v
  }

  return upperEnvsMap, nil
}