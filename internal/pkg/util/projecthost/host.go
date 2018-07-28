package projecthost

// Host is the interface for all project host platforms (i.e. GitHub)
type Host interface {
  LatestSha() (string, error)
}