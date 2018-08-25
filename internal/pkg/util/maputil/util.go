package maputil

import "fmt"

func MergeMaps(base map[string]string, priorities map[string]string) map[string]string {
  for k, v := range priorities {
    base[k] = v
  }

  return base
}

func FromStrSlice(values []interface{}) (map[string]interface{}, error) {
  // First split values by even and odd indexes.
  var evens []string
  var odds []interface{}

  for i, val := range values {
    if i % 2 == 0 {
      // Convert all evens to strings.
      strVal, ok := val.(string)
      if !ok {
        err := fmt.Errorf("error creating map[string]interface{} from array: all even-indexed keys need to be string")
        return nil, err
      }

      evens = append(evens, strVal)
    } else {
      odds = append(odds, val)
    }
  }

  // Then create a map matching pairs based on index.
  j := 0
  m := map[string]interface{}{}
  for j < len(evens) {
    m[evens[j]] = odds[j]
    j++
  }

  return m, nil
}