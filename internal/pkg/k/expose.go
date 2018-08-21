package k

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  corev1 "k8s.io/api/core/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
  "k8s.io/apimachinery/pkg/watch"
  "fmt"
  "github.com/sweettea-io/rest-api/internal/pkg/util/maputil"
)

type Expose struct {
  // Establish on Init
  Deploy        *model.Deploy
  ApiCluster    *model.ApiCluster
  ServiceName   string
  Port          int32
  TargetPort    int32
  Labels        map[string]string
  ResultChannel <-chan Result

  // K8S resources
  Namespace     string
  Client        *typedcorev1.CoreV1Client
  ServiceSpec   *corev1.ServiceSpec
  Service       *corev1.Service
}

func (expose *Expose) Init(args map[string]interface{}) error {
  // Parse args.
  expose.Deploy = args["deploy"].(*model.Deploy)
  expose.Port = args["port"].(int32)
  expose.TargetPort = args["targetPort"].(int32)

  // Initialize the result channel.
  expose.ResultChannel = make(chan Result)

  // Set further models through associations.
  expose.ApiCluster = &expose.Deploy.ApiCluster

  // Service name is simply the deployment name.
  expose.ServiceName = expose.Deploy.DeploymentName

  // Create labels map for Service metadata.
  expose.Labels = map[string]string{
    "name": expose.ServiceName,
  }

  // Check if need to expand service labels.
  cloud := app.Config.Cloud()

  // Add SSL labels if port is 443
  if expose.Port == 443 && cloud != nil {
    expose.Labels = maputil.MergeMaps(expose.Labels, cloud.SSLServiceLabels())
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

func (expose *Expose) GetResultChannel() <-chan Result {
  return expose.ResultChannel
}

func (expose *Expose) Watch() {
  // Get a namespaced service watcher channel.
  ch, err := ServiceWatcherChannel(expose.Client, expose.Namespace, expose.ServiceName)

  if err != nil {
    expose.ResultChannel <- Result{Ok: false, Error: err}
    return
  }

  // Start watching for events.
  for event := range ch {
    if result := expose.checkEventForResult(event); result != nil {
      expose.ResultChannel <- *result
      return
    }
  }
}

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

func (expose *Expose) checkEventForResult(event watch.Event) *Result {
  // Error out early if ever encounter a watch error.
  if event.Type == watch.Error {
    err := fmt.Errorf("Service %s encountered error.", expose.ServiceName)
    app.Log.Errorf(err.Error())
    return &Result{Ok: false, Error: err}
  }

  // Parse service resource from event.
  service, ok := event.Object.(*corev1.Service)

  if !ok {
    err := fmt.Errorf("Job %s encountered unexpected event object type.", expose.ServiceName)
    app.Log.Errorf(err.Error())
    return &Result{Ok: false, Error: err}
  }

  // Check if service has a LoadBalancer Ingress hostname yet.
  // Only return successful event once that has happened.
  ingressList := service.Status.LoadBalancer.Ingress

  if len(ingressList) == 0 {
    return nil
  }

  // We only care about first one (there should only ever be one).
  lbHost := ingressList[0].Hostname

  if lbHost == "" {
    return nil
  }

  return &Result{
    Ok: true,
    Meta: map[string]interface{}{
      "lbHost": lbHost,
    },
  }
}