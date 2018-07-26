package projecthost

import "strings"

// Supported hosts
const Github = "github"

// FromNsp returns the host for a given project namespace (if a supported one exists).
func FromNsp(nsp string) string {
  if strings.HasPrefix(nsp, "gitub.com") {
    return Github
  }

  return ""
}