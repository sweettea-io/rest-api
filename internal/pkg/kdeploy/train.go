package kdeploy

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/service/trainjobsvc"
)

type Train struct {
  TrainJob *model.TrainJob
}

func (t *Train) Init(args map[string]interface{}) error {
  // Find TrainJob by ID.
  trainJob, err := trainjobsvc.FromID(args["resourceID"].(uint))

  if err != nil {
    return err
  }

  t.TrainJob = trainJob
  return nil
}

func (t *Train) Configure() error {
  return nil
}

func (t *Train) Perform() error {
  return nil
}

func (t *Train) Watch() {

}