package jobs

type names struct {
  CreateTrainJob string
  BuildDeploy    string
  TrainDeploy    string
  ApiDeploy      string
}

// Names for all supported jobs.
var Names = &names{
  CreateTrainJob: "create_train_job",
  BuildDeploy: "build_server_deploy",
  TrainDeploy: "train_deploy",
  ApiDeploy: "api_deploy",
}