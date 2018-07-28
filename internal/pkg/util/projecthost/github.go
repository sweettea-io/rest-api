package projecthost

import "github.com/sweettea-io/rest-api/internal/pkg/model"

const GitHubName = "github"
const GitHubDomain = "github.com"

type GitHub struct {
  Project *model.Project
}

func (gh *GitHub) LatestSha() (string, error) {
  return "woo", nil
}