package projecthost

const GitHubName = "github"
const GitHubDomain = "github.com"

type GitHub struct {}

func (gh *GitHub) LatestSha() (string, error) {
  return "woo", nil
}