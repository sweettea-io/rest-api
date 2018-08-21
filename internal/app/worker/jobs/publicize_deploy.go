package jobs

import (
  "fmt"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/k"
)

/*
  PublicizeDeploy creates a K8S service from a K8S deployment and makes it available at a public url.

  Args:
    deployID (uint) ID of the Deploy model attached to the K8S deployment to publicize
*/
func (c *Context) PublicizeDeploy(job *work.Job) error {
  // Extract args from job.
  deployID := uint(job.ArgInt64("deployID"))

  if err := job.ArgError(); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Get Deploy by ID.
  deploy, err := deploysvc.FromID(deployID)

  if err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Deploy shouldn't already be public...
  if deploy.Public {
    deploysvc.FailByID(deployID)
    err := fmt.Errorf("Deploy(id=%v) can't be publicized -- it is already public.", deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Used to create a K8S Service from a Deploy.
  serviceCreation := k.Expose{}

  // Initialize service creation.
  if err := serviceCreation.Init(map[string]interface{}{"deploy": deploy}); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Configure service resources.
  if err := serviceCreation.Configure(); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Create K8S LoadBalancer Service from Deploy's deployment.
  if err := serviceCreation.Perform(); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Get ELB from service
  // Assign ELB to Deploy.LBHost
  // Add a CNAME record to your domain's DNS mapping the publicized deploy url to the ELB url
  // Update Deploy.Public to true

  return nil
}