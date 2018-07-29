package kdeploy

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
)

type Api struct {
  Deploy *model.Deploy
}

func (api *Api) Init(args map[string]interface{}) error {
  // Find Deploy by ID.
  deploy, err := deploysvc.FromID(args["resourceID"].(uint))

  if err != nil {
    return err
  }

  api.Deploy = deploy
  return nil
}

func (api *Api) Configure() error {
  return nil
}

func (api *Api) Perform() error {
  return nil
}

func (api *Api) Watch() {

}