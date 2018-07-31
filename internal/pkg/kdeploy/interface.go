package kdeploy

type KDeploy interface {
  Init(args map[string]interface{}) error
  Configure() error
  Perform() error
}