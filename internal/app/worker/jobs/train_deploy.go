package jobs

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
)

/*
  TrainDeploy deploys a TrainJob to the SweetTea Train Cluster for model training.

  Args:
    trainJobID (uint) ID of the TrainJob to deploy
*/
func (c *Context) TrainDeploy(job *work.Job) error {
  // Ensure Train cluster exists first.
  if !app.Config.TrainClusterConfigured() {
    err := fmt.Errorf("train cluster not configured -- leaving CreateTrainJob")
    app.Log.Errorln(err.Error())
    return err
  }

  // Extract args from job.
  trainJobID := uint(job.ArgInt64("trainJobID"))

  if err := job.ArgError(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Find TrainJob by ID.
  trainJob, err := trainjobsvc.FromID(trainJobID)

  if err != nil {
    app.Log.Errorln(err.Error())
    return err
  }



  return nil
}