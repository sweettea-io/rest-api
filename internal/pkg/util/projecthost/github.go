package projecthost

import (
  "context"
  "fmt"
  "github.com/google/go-github/github"
  "golang.org/x/oauth2"
  "net/http"
)

const GitHubName = "github"
const GitHubDomain = "github.com"

type GitHub struct {
  Client *github.Client
  Ctx context.Context
}

func (gh *GitHub) Configure(token string) {
  // Initialize http client and background context.
  var httpClient *http.Client
  gh.Ctx = context.Background()

  // Set http client as OAuth2 client if GH access token exists.
  if token != "" {
    tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    httpClient = oauth2.NewClient(gh.Ctx, tokenSource)
  }

  gh.Client = github.NewClient(httpClient)
}

func (gh *GitHub) LatestSha(owner, repo string) (string, error) {
  // We only want the latest commit.
  commitListOpts := github.CommitsListOptions{}
  commitListOpts.Page = 1
  commitListOpts.PerPage = 1

  // Get commits for repo.
  commits, _, err := gh.Client.Repositories.ListCommits(
    gh.Ctx,
    owner,
    repo,
    &commitListOpts,
  )

  if err != nil {
    return "", fmt.Errorf("error listing commits for GitHub project: %s/%s", owner, repo)
  }

  // Ensure repo even has commits...
  if len(commits) == 0 {
    return "", fmt.Errorf("no commits exist yet for GitHub project: %s/%s", owner, repo)
  }

  return *commits[0].SHA, nil
}