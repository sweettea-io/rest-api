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

func InterfaceToBytes(val interface{}) ([]byte, error) {
  if val == nil {
    return nil, nil
  }

  bytes, ok := val.([]byte)

  if !ok {
    return nil, fmt.Errorf("error converting interface{} to []byte")
  }

  return bytes, nil
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