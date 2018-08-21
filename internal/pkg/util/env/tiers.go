package env

// Supported environments.
const (
  Test = "test"
  Local = "local"
  Dev = "dev"
  Staging = "staging"
  Prod ="prod"
)

var tiers = map[string]bool {
  Test: true,
  Local: true,
  Dev: true,
  Staging: true,
  Prod: true,
}

// IsValidEnv returns whether the provided environment is a supported environment.
func IsValidEnv(env string) bool {
  return tiers[env] == true
}