package projectsvc

import (
  "strings"
  "github.com/sweettea-io/rest-api/internal/pkg/util/projecthost"
  "github.com/gosimple/slug"
)

func IsValidNsp(nsp string) bool {
  comps := strings.Split(nsp, "/")

  return len(comps) == 3 &&
    projecthost.IsValidHostDomain(comps[0]) &&
    slug.IsSlug(comps[1]) &&
    slug.IsSlug(comps[2])
}