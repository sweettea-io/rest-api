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

  serviceCreationArgs := map[string]interface{}{
    "deploy": deploy,
    "port": 443,
    "targetPort": 443,
  }

  // Initialize service creation.
  if err := serviceCreation.Init(serviceCreationArgs); err != nil {
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

  // Get channel to watch for successful service creation.
  resultCh := serviceCreation.GetResultChannel()

  // Watch service until it has a LoadBalancer hostname assigned to it.
  go serviceCreation.Watch()
  serviceResult := <-resultCh

  // Error out if service watching failed.
  if !serviceResult.Ok {
    deploysvc.FailByID(deployID)
    app.Log.Errorf(serviceResult.Error.Error())
    return serviceResult.Error
  }

  // Parse LoadBalancer hostname from result.
  lbHost, ok := serviceResult.Meta["lbHost"].(string)

  if !ok || lbHost == "" {
    deploysvc.FailByID(deployID)
    err := fmt.Errorf("Error parsing lbHost from Service creation result for Deploy(id=%v).", deployID)
    app.Log.Errorln(err.Error())
    return err
  }

  // Assign LoadBalancer hostname to Deploy.
  if err := deploysvc.RegisterLoadBalancerHost(deploy, lbHost); err != nil {
    deploysvc.FailByID(deployID)
    app.Log.Errorln(err.Error())
    return err
  }
  
  // Add a CNAME record to your domain's DNS mapping the publicized deploy url to the ELB url
  // Update Deploy.Public to true

  return nil
}