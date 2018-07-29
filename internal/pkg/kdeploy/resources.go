package kdeploy

import (
  "fmt"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
  corev1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
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

func ConfigureCoreV1(context string) (*typedcorev1.CoreV1Client, string, error) {
  // Configure k8s rest client for provided context.
  restConfig, nsp, err := GetRestConfig(context)

  if err != nil {
    return nil, "", err
  }

  // Create CoreV1 client from rest client.
  coreV1Client, err := typedcorev1.NewForConfig(restConfig)

  if err != nil {
    return nil, "", fmt.Errorf("error creating New CoreV1Client: %s", err.Error())
  }

  return coreV1Client, nsp, nil
}

func VolumeMounts(vms []map[string]string) []corev1.VolumeMount {
  volumeMounts := []corev1.VolumeMount{}

  for _, vm := range vms {
    volumeMounts = append(volumeMounts, VolumeMount(vm["name"], vm["path"]))
  }

  return volumeMounts
}

func VolumeMount(name string, path string) corev1.VolumeMount {
  return corev1.VolumeMount{Name: name, MountPath: path}
}

func Volumes(vols []map[string]string) []corev1.Volume {
  volumes := []corev1.Volume{}

  for _, vol := range vols {
    volumes = append(volumes, Volume(vol["name"], vol["path"]))
  }

  return volumes
}

func Volume(name string, path string) corev1.Volume {
  return corev1.Volume{
    Name: name,
    VolumeSource: corev1.VolumeSource{
      HostPath: &corev1.HostPathVolumeSource{
        Path: path,
      },
    },
  }
}

func EnvVars(envVars map[string]string) []corev1.EnvVar {
  envs := []corev1.EnvVar{}

  for k, v := range envVars {
    envs = append(envs, corev1.EnvVar{Name: k, Value: v})
  }

  return envs
}

func Containers(conts []map[string]interface{}) []corev1.Container {
  containers := []corev1.Container{}

  for _, cont := range conts {
    containers = append(containers, Container(cont))
  }

  return containers
}

func Container(args map[string]interface{}) corev1.Container {
  // Parse container args.
  name := args["name"].(string)
  image := args["image"].(string)
  envs, envsProvided := args["envs"]
  ports, portsProvided := args["ports"]
  volMounts, volMountsProvided := args["volumeMounts"]

  // Create container will Name and Image (the required params)
  container := corev1.Container{
    Name: name,
    Image: image,
  }

  // Add in the optional params:

  // Envs (optional)
  if envsProvided {
    container.Env = envs.([]corev1.EnvVar)
  }

  // Ports (optional)
  if portsProvided {
    container.Ports = ports.([]corev1.ContainerPort)
  }

  // VolumeMounts (optional)
  if volMountsProvided {
    container.VolumeMounts = volMounts.([]corev1.VolumeMount)
  }

  return container
}

func Pod(args map[string]interface{}) *corev1.Pod {
  // Parse args.
  label := args["label"].(string)
  containers := args["containers"].([]corev1.Container)
  restart := args["restart"].(corev1.RestartPolicy)
  vols, volsProvided := args["volumes"]

  // Create PodSpec.
  podSpec := corev1.PodSpec{
    Containers: containers,
    RestartPolicy: restart,
  }

  // Add volumes to PodSpec if provided.
  if volsProvided {
    podSpec.Volumes = vols.([]corev1.Volume)
  }

  // Create and return Pod.
  return &corev1.Pod{
    ObjectMeta: metav1.ObjectMeta{
      Labels: map[string]string{
        "app": label,
      },
    },
    Spec: podSpec,
  }
}