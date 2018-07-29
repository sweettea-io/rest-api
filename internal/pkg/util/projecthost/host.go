package projecthost

// Host is the interface for all project host platforms (i.e. GitHub)
type Host interface {
  Init()
  Configure()
  LatestSha(owner string, repo string) (string, error)
  GetToken() string
}