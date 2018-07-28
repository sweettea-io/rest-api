package projecthost

import (
  "context"
  "github.com/google/go-github/github"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "golang.org/x/oauth2"
  "fmt"
)

const GitHubName = "github"
const GitHubDomain = "github.com"

type GitHub struct {
  Project *model.Project
  Client *github.Client
  Ctx context.Context
}

func (gh *GitHub) Configure() {
  // Get GitHub access token from config.
  accessToken := app.Config.GitHubAccessToken

  // Initialize background context.
  gh.Ctx = context.Background()

  // Initialize GitHub API client based on whether or not it needs to be authed.
  if accessToken == "" {
    gh.Client = github.NewClient(nil)
  } else {
    gh.Client = github.NewClient(oauth2.NewClient(
      gh.Ctx,
      oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})),
    )
  }
}

func (gh *GitHub) LatestSha() (string, error) {
  // We only want the latest commit.
  commitListOpts := &github.CommitsListOptions{}
  commitListOpts.Page = 1
  commitListOpts.PerPage = 1

  // Get commits for repo.
  commits, _, err := gh.Client.Repositories.ListCommits(
    gh.Ctx,
    gh.Project.Owner(),
    gh.Project.Repo(),
    commitListOpts,
  )

  if err != nil {
    return "", fmt.Errorf("error listing commits for GitHub project: %s", gh.Project.Nsp)
  }

  // Ensure repo even has commits...
  if len(commits) == 0 {
    return "", fmt.Errorf("no commits exist yet for GitHub project: %s", gh.Project.Nsp)
  }

  return *commits[0].SHA, nil
}