package maputil

func MergeMaps(base map[string]string, priorities map[string]string) map[string]string {
  for k, v := range priorities {
    base[k] = v
  }

  return base
}
