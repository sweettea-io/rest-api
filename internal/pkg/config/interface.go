package config

type ConfigItf interface {
  BaseRoute() string
  OnTest() bool
  OnLocal() bool
  OnDev() bool
  OnStaging() bool
  OnProd() bool
  TrainClusterConfigured() bool
  BuildpackEnvs() map[string]string
}