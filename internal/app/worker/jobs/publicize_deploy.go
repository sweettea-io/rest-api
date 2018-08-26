package jobs

import (
  "fmt"
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/pkg/service/deploysvc"
  "github.com/sweettea-io/rest-api/internal/pkg/k"
  "github.com/sweettea-io/rest-api/internal/pkg/util/dns"
  "github.com/sweettea-io/rest-api/internal/app"
)

/*
  PublicizeDeploy creates a K8S service from a K8S deployment and makes it available at a public url.

  Args:
    deployID (uint) ID of the Deploy model attached to the K8S deployment to publicize
    logKey   (string) log key for buildable
*/
func (c *Context) PublicizeDeploy(job *work.Job) error {
  // Extract args from job.
  deployID := uint(job.ArgInt64("deployID"))
  logKey := job.ArgString("logKey")

  if err := job.ArgError(); err != nil {
    if logKey != "" {
      return failDeploy(deployID, err, logKey, "Arg error occurred inside publicize deploy job.")
    }

    app.Log.Errorln(err.Error())
    return err
  }

  // Get Deploy by ID.
  deploy, err := deploysvc.FromID(deployID)
  if err != nil {
    return failDeploy(deployID, err, logKey, "Deploy not found.")
  }

  // Deploy shouldn't already be public...
  if deploy.Public {
    err := fmt.Errorf("Can't publicize and already public deploy.")
    return failDeploy(deployID, err, logKey, err.Error())
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
    return failDeploy(deployID, err, logKey, "Failed to initialize service creator.")
  }

  // Configure service resources.
  if err := serviceCreation.Configure(); err != nil {
    return failDeploy(deployID, err, logKey, "Failed to configure service creation resources.")
  }

  // Create K8S LoadBalancer Service from Deploy's deployment.
  if err := serviceCreation.Perform(); err != nil {
    return failDeploy(deployID, err, logKey, "Failed to perform service creation.")
  }

  // Get channel to watch for successful service creation.
  resultCh := serviceCreation.GetResultChannel()

  // Watch service until it has a LoadBalancer hostname assigned to it.
  go serviceCreation.Watch()
  serviceResult := <-resultCh

  // Error out if service watching failed.
  if !serviceResult.Ok {
    return failDeploy(deployID, serviceResult.Error, logKey, "Failed to publicize deploy as a service.")
  }

  // Parse LoadBalancer hostname from result.
  lbHostname, ok := serviceResult.Meta["lbHostname"].(string)

  if !ok || lbHostname == "" {
    err := fmt.Errorf("Failed to parse load balancer hostname for service.")
    return failDeploy(deployID, err, logKey, err.Error())
  }

  // Assign LoadBalancer hostname to Deploy.
  if err := deploysvc.RegisterLoadBalancerHost(deploy, lbHostname); err != nil {
    return failDeploy(deployID, err, logKey, "Failed to register load balancer hostname with deploy.")
  }

  // Ensure DNS service is currently configured.
  if dns.CurrentDNS == nil {
    err := fmt.Errorf("SweetTea DNS not currently configured.")
    return failDeploy(deployID, err, logKey, err.Error())
  }

  // Generate new subdomain for this Deploy, using currently configured Domain.
  subdomain := deploysvc.NewHostname(deploy)

  // Add a CNAME RR to your domain's DNS, aliasing the Deploy's hostname to the LoadBalancer hostname.
  if err := dns.CurrentDNS.UpsertRR(dns.CNAME, subdomain, []string{lbHostname}, 60); err != nil {
    return failDeploy(deployID, err, logKey, "Failed to register CNAME record for load balancer.")
  }

  // Mark Deploy as public now at specified subdomain.
  if err := deploysvc.Publicize(deploy, subdomain); err != nil {
    return failDeploy(deployID, err, logKey, "Failed to register deploy as public.")
  }

  // TODO: Stream message back successfully disconnecting client.

  return nil
}