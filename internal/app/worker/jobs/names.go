package jobs

type names struct {
  CreateTrainJob string
  CreateDeploy   string
  BuildDeploy    string
  TrainDeploy    string
  ApiDeploy      string
}

// Names for all supported jobs.
var Names = &names{
  CreateTrainJob: "create_train_job",
  CreateDeploy: "create_deploy",
  BuildDeploy: "build_server_deploy",
  TrainDeploy: "train_deploy",
  ApiDeploy: "api_deploy",
}