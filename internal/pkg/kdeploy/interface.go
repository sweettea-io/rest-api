package kdeploy

type KDeploy interface {
  Init(args map[string]interface{}) error
  Configure() error
  Perform() error
  GetResultChannel() <-chan Result
  Watch()
}

type Result struct {
  Ok    bool
  Error error
}