package typeconvert

import "fmt"

func BytesToStr(data []byte) string {
  return string(data[:])
}

func BoolToStr(b bool) string {
  if b == true {
    return "true"
  }

  return "false"
}

func InterfaceToStr(val interface{}) (string, error) {
  if val == nil {
    return "", nil
  }

  strVal, ok := val.(string)

  if !ok {
    return "", fmt.Errorf("error converting interface{} to string")
  }

  return strVal, nil
}

func InterfaceToBool(val interface{}) (bool, error) {
  if val == nil {
    return false, nil
  }

  boolVal, ok := val.(bool)

  if !ok {
    return false, fmt.Errorf("error converting interface{} to bool")
  }

  return boolVal, nil
}