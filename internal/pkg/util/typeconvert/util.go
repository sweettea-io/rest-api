package typeconvert

func BytesToStr(data []byte) string {
  return string(data[:])
}

func BoolToStr(b bool) string {
  if b == true {
    return "true"
  }

  return "false"
}