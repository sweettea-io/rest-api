package jobs

import "github.com/sweettea-io/work"

/*
  UpdateDeploy needs to call the right next job based on what actually needs to be updated. Options:
  (1) sha --> new image needs to be built and then PATCHed on the deployment
  (2) model/modelVersion --> env vars need to be PATCHed on the deployment

  Args:

*/
func (c *Context) UpdateDeploy(job *work.Job) error {

}