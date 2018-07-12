package utils

type envTiers struct {
  Test    string
  Local   string
  Dev     string
  Staging string
  Prod    string
}

var Envs = envTiers{
  Test:    "test",
  Local:   "local",
  Dev:     "dev",
  Staging: "staging",
  Prod:    "prod",
}

var validEnvs = map[string]bool {
  Envs.Test: true,
  Envs.Local: true,
  Envs.Dev: true,
  Envs.Staging: true,
  Envs.Prod: true,
}

// Check if an env is supported.
func IsValidEnv(env string) bool {
  return validEnvs[env] == true
}