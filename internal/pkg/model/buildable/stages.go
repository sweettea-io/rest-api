package buildable

// Stages of a deploy/trainJob, in order.
const (
  Created = "created"
  BuildScheduled = "build_scheduled"
  Building = "building"
  DeployScheduled = "deploy_scheduled"
  Deploying = "deploying"
  Deployed = "deployed"
)