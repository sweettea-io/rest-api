package modelsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "strings"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
)

func Breakdown(modelName string) (string, string) {
  if modelName == "" {
    return model.DefaultModelSlug, ""
  }

  breakdown := strings.Split(modelName, ":")
  slug := str.Slugify(breakdown[0])
  version := ""

  if len(breakdown) > 1 {
    version = breakdown[1]
  }

  return slug, version
}
