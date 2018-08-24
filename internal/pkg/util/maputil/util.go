package maputil

func MergeMaps(base map[string]string, priorities map[string]string) map[string]string {
  for k, v := range priorities {
    base[k] = v
  }

  return base
}

func FromStrSlice(values []string) map[string]string {
  // First split values by even and odd indexes.
  var evens []string
  var odds []string

  for i, val := range values {
    if i % 2 == 0 {
      evens = append(evens, val)
    } else {
      odds = append(odds, val)
    }
  }

  // Then create a map matching pairs based on index.
  j := 0
  m := map[string]string{}
  for j < len(evens) {
    m[evens[j]] = odds[j]
    j++
  }

  return m
}