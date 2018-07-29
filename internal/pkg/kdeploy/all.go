package kdeploy

import (
  "fmt"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
  corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func ConfigureClient(ctx string) (*rest.Config, string, error) {
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

func NewCoreV1(restConfig *rest.Config) (*corev1.CoreV1Client, error) {
  coreV1Client, err := corev1.NewForConfig(restConfig)

  if err != nil {
    return nil, fmt.Errorf("error creating New CoreV1Client: %s", err.Error())
  }

  return coreV1Client, nil
}