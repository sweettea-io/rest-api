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