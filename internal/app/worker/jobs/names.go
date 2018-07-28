package jobs

type names struct {
  CreateTrainJob string
  TrainDeploy    string
}

// Names for all supported jobs.
var Names = &names{
  CreateTrainJob: "create_train_job",
  TrainDeploy: "train_deploy",
}