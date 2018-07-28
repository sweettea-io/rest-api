package projecthost

const GH = "github"
const GH_DOMAIN = "github.com"

type GitHub struct {}

func (gh *GitHub) LatestSha() (string, error) {
  return "", nil
}