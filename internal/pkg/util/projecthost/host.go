package projecthost

// Host is the interface for all project host platforms (i.e. GitHub)
type Host interface {
  Configure(token string)
  LatestSha(owner string, repo string) (string, error)
}