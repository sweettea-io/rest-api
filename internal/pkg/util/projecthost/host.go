package projecthost

// Host is the interface for all project host platforms (i.e. GitHub)
type Host interface {
  LatestSha() (string, error)
}

func FromName(name string) Host {
  switch name {
  case GH:
    return &GitHub{}
  default:
    return nil
  }
}