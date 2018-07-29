package kdeploy

import (
  "fmt"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
)

func GetRestConfig(ctx string) (*rest.Config, string, error) {
  // Build config from KUBECONFIG file and desired context.
  kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
    clientcmd.NewDefaultClientConfigLoadingRules(),
    &clientcmd.ConfigOverrides{
      CurrentContext: ctx,
    },
  )

  // Get namespace for current context.
  namespace, _, err := kubeconfig.Namespace()

  if err != nil {
    return nil, "", fmt.Errorf("error determining KUBECONFIG namespace: %s", err.Error())
  }

  // Get rest client config from kubeconfig.
  restConfig, err := kubeconfig.ClientConfig()

  if err != nil {
    return nil, "", fmt.Errorf("error creating kubeconfig rest client config: %s", err.Error())
  }

  return restConfig, namespace, nil
}