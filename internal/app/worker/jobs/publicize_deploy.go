package jobs

import (
  "fmt"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
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
    failDeploy(deployID, err)
  }

  // Get Deploy by ID.
  deploy, err := deploysvc.FromID(deployID)

  if err != nil {
    return failDeploy(deployID, err)
  }

  // Deploy shouldn't already be public...
  if deploy.Public {
    return failDeploy(deployID, fmt.Errorf("Deploy(id=%v) can't be publicized -- it is already public.", deployID))
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
    return failDeploy(deployID, err)
  }

  // Configure service resources.
  if err := serviceCreation.Configure(); err != nil {
    return failDeploy(deployID, err)
  }

  // Create K8S LoadBalancer Service from Deploy's deployment.
  if err := serviceCreation.Perform(); err != nil {
    return failDeploy(deployID, err)
  }

  // Get channel to watch for successful service creation.
  resultCh := serviceCreation.GetResultChannel()

  // Watch service until it has a LoadBalancer hostname assigned to it.
  go serviceCreation.Watch()
  serviceResult := <-resultCh

  // Error out if service watching failed.
  if !serviceResult.Ok {
    return failDeploy(deployID, serviceResult.Error)
  }

  // Parse LoadBalancer hostname from result.
  lbHost, ok := serviceResult.Meta["lbHost"].(string)

  if !ok || lbHost == "" {
    return failDeploy(deployID, fmt.Errorf("Error parsing lbHost from Service creation result for Deploy(id=%v).", deployID))
  }

  // Assign LoadBalancer hostname to Deploy.
  if err := deploysvc.RegisterLoadBalancerHost(deploy, lbHost); err != nil {
    return failDeploy(deployID, err)
  }

  // Add a CNAME record to your domain's DNS mapping the publicized deploy url to the ELB url
  // Update Deploy.Public to true

  return nil
}