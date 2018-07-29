package kdeploy

import (
  "fmt"
  corev1 "k8s.io/api/core/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func ConfigureCoreV1(ctx string) (*typedcorev1.CoreV1Client, string, error) {
  // Configure k8s rest client for provided context.
  restConfig, nsp, err := GetRestConfig(ctx)

  if err != nil {
    return nil, "", err
  }

  // Create CoreV1 client from rest client.
  coreV1Client, err := corev1.NewForConfig(restConfig)

  if err != nil {
    return nil, "", fmt.Errorf("error creating New CoreV1Client: %s", err.Error())
  }

  return coreV1Client, nsp, nil
}

func CoreV1Volume(name string, path string) corev1.Volume {
  return corev1.Volume{
    Name: name,
    VolumeSource: corev1.VolumeSource{
      HostPath: &corev1.HostPathVolumeSource{
        Path: path,
      },
    },
  }
}

func CoreV1EnvVars(envVars map[string]string) []corev1.EnvVar {
  envs := []corev1.EnvVar{}

  for k, v := range envVars {
    envs = append(envs, corev1.EnvVar{Name: k, Value: v})
  }

  return envs
}