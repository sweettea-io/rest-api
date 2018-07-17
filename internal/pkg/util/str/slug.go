package str

import "github.com/gosimple/slug"

func Slugify(text string) string {
  return slug.MakeLang(text, "en")
}