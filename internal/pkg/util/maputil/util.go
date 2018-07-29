package maputil

func MergeMaps(base map[string]string, additions map[string]string) map[string]string {
  for k, v := range additions {
    base[k] = v
  }

  return base
}
