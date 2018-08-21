package jobs

type names struct {
  CreateTrainJob  string
  CreateDeploy    string
  BuildDeploy     string
  TrainDeploy     string
  ApiDeploy       string
  UpdateDeploy    string
  ApiUpdate       string
  PublicizeDeploy string
}

// Names for all supported jobs.
var Names = &names{
  CreateTrainJob: "create_train_job",
  CreateDeploy: "create_deploy",
  BuildDeploy: "build_server_deploy",
  TrainDeploy: "train_deploy",
  ApiDeploy: "api_deploy",
  UpdateDeploy: "update_deploy",
  ApiUpdate: "api_update",
  PublicizeDeploy: "publicize_deploy",
}