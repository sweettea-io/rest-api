package k

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  corev1 "k8s.io/api/core/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Expose struct {
  // Establish on Init
  Deploy      *model.Deploy
  ApiCluster  *model.ApiCluster
  ServiceName string
  Port        int32
  TargetPort  int32
  Labels      map[string]string

  // K8S resources
  Namespace   string
  Client      *typedcorev1.CoreV1Client
  ServiceSpec *corev1.ServiceSpec
  Service     *corev1.Service
}

func (expose *Expose) Init(args map[string]interface{}) error {
  // Parse args.
  expose.Deploy = args["deploy"].(*model.Deploy)
  expose.Port = args["port"].(int32)
  expose.TargetPort = args["targetPort"].(int32)

  // Set further models through associations.
  expose.ApiCluster = &expose.Deploy.ApiCluster

  // Service name is simply the deployment name.
  expose.ServiceName = expose.Deploy.DeploymentName

  // Create labels map for Service metadata.
  expose.Labels = map[string]string{
    "name": expose.ServiceName,
  }

  // Add SSL labels if port is 443
  if expose.Port == 443 {
    expose.Labels["service.beta.kubernetes.io/aws-load-balancer-ssl-cert"] = app.Config.WildcardSSLCertArn
    expose.Labels["service.beta.kubernetes.io/aws-load-balancer-ssl-ports"] = "443"
  }

  return nil
}

func (expose *Expose) Configure() error {
  // Configure k8s client.
  if err := expose.makeClient(); err != nil {
    return err
  }

  // Build up K8S resources into the Service we wish to create.
  expose.makeServiceSpec()
  expose.makeService()

  return nil
}

// Perform creates a service from the Deploy's deployment.
func (expose *Expose) Perform() error {
  return CreateService(expose.Client, expose.Namespace, expose.Service)
}

// TODO: add watch functionality and check for loadbalancer url that way

func (expose *Expose) makeClient() error {
  // Configure V1Beta1 client.
  client, nsp, err := ConfigureCoreV1(expose.ApiCluster.Slug)

  if err != nil {
    return err
  }

  // Store refs to client and namespace.
  expose.Client = client
  expose.Namespace = nsp

  return nil
}

func (expose *Expose) makeServiceSpec() {
  expose.ServiceSpec = ServiceSpec(expose.ServiceName, expose.Port, expose.TargetPort)
}

func (expose *Expose) makeService() {
  expose.Service = Service(expose.ServiceSpec, expose.Labels)
}