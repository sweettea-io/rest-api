package model

type Buildable interface {
  GetCommit() *Commit
  GetUid() string
}

type buildStages struct {
  Created         string
  BuildScheduled  string
  Building        string
  DeployScheduled string  
  Deploying       string
  Deployed        string
}

var BuildStages = buildStages{
  Created: "created",
  BuildScheduled: "build_scheduled",
  Building: "building",
  DeployScheduled: "deploy_scheduled",
  Deploying: "deploying",
  Deployed: "deployed",
}