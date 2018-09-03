package projecthost

// Host is the interface for all project host platforms (i.e. GitHub)
type Host interface {
  Configure()
  LatestSha(owner string, repo string) (string, error)
  EnsureCommitExists(owner string, repo string, sha string) error
  GetToken() string
}

var validHostDomains = map[string]bool{
  GitHubDomain: true,
}

func IsValidHostDomain(domain string) bool {
  return validHostDomains[domain] == true
}