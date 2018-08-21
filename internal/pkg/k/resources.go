package k

import (
  "fmt"
  "k8s.io/apimachinery/pkg/watch"
  "k8s.io/api/extensions/v1beta1"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/apimachinery/pkg/util/intstr"
  corev1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
  typedv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
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

func ConfigureV1Beta1(context string) (*typedv1beta1.ExtensionsV1beta1Client, string, error) {
  // Configure k8s rest client for provided context.
  restConfig, nsp, err := GetRestConfig(context)

  if err != nil {
    return nil, "", err
  }

  // Create V1Beta1 client from rest client.
  v1beta1Client, err := typedv1beta1.NewForConfig(restConfig)

  if err != nil {
    return nil, "", fmt.Errorf("error creating New ExtensionsV1beta1Client: %s", err.Error())
  }

  return v1beta1Client, nsp, nil
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

func Ports(ports []int32) []corev1.ContainerPort {
  cPorts := []corev1.ContainerPort{}

  for _, p := range ports {
    cPorts = append(cPorts, corev1.ContainerPort{ContainerPort: p})
  }

  return cPorts
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

func PodTemplateSpec(args map[string]interface{}) *corev1.PodTemplateSpec {
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

  // Create and return PodTemplateSpec.
  return &corev1.PodTemplateSpec{
    ObjectMeta: metav1.ObjectMeta{
      Labels: map[string]string{
        "app": label,
      },
    },
    Spec: podSpec,
  }
}

func DeploymentSpec(podTemplateSpec *corev1.PodTemplateSpec, replicas int32) *v1beta1.DeploymentSpec {
  return &v1beta1.DeploymentSpec{
    Template: *podTemplateSpec,
    Replicas: &replicas,
  }
}

func Deployment(deploymentSpec *v1beta1.DeploymentSpec, name string) *v1beta1.Deployment {
  return &v1beta1.Deployment{
    TypeMeta: metav1.TypeMeta{
      APIVersion: "extensions/v1beta1",
      Kind: "Deployment",
    },
    ObjectMeta: metav1.ObjectMeta{
      Labels: map[string]string{
        "name": name,
      },
    },
    Spec: *deploymentSpec,
  }
}

func ServiceSpec(deploymentName string, port int32, targetPort int32) *corev1.ServiceSpec {
  return &corev1.ServiceSpec{
    Ports: []corev1.ServicePort{{
      Port: port,
      TargetPort: intstr.IntOrString{IntVal: targetPort},
    }},
    Selector: map[string]string{
      "app": deploymentName,
    },
    Type: corev1.ServiceTypeLoadBalancer,
  }
}

func Service(serviceSpec *corev1.ServiceSpec, labels map[string]string) *corev1.Service {
  return &corev1.Service{
    TypeMeta: metav1.TypeMeta{
      APIVersion: "core/v1",
      Kind: "Service",
    },
    ObjectMeta: metav1.ObjectMeta{Labels: labels},
    Spec: *serviceSpec,
  }
}

func CreatePod(client *typedcorev1.CoreV1Client, nsp string, pod *corev1.Pod, targetCluster string) error {
  if _, err := client.Pods(nsp).Create(pod); err != nil {
    return fmt.Errorf("error creating pod during %s deploy: %s", targetCluster, err.Error())
  }

  return nil
}

func CreateDeployment(client *typedv1beta1.ExtensionsV1beta1Client, nsp string, deployment *v1beta1.Deployment) error {
  if _, err := client.Deployments(nsp).Create(deployment); err != nil {
    return fmt.Errorf("error creating deployment during API deploy: %s", err.Error())
  }

  return nil
}

func UpdateDeployment(client *typedv1beta1.ExtensionsV1beta1Client, nsp string, deployment *v1beta1.Deployment) error {
  if _, err := client.Deployments(nsp).Update(deployment); err != nil {
    return fmt.Errorf("error updating deployment during API deploy: %s", err.Error())
  }

  return nil
}

func CreateService(client *typedcorev1.CoreV1Client, nsp string, service *corev1.Service) error {
  if _, err := client.Services(nsp).Create(service); err != nil {
    return fmt.Errorf("error creating k8s service: %s", err.Error())
  }

  return nil
}

func PodWatcherChannel(client *typedcorev1.CoreV1Client, nsp string, deployName string) (<-chan watch.Event, error) {
  // Define selector options to only find the pod we just deployed.
  watchOpts := metav1.ListOptions{
    LabelSelector: fmt.Sprintf("app=%s", deployName),
  }

  // Get a namespaced pod watcher object.
  watcher, err := client.Pods(nsp).Watch(watchOpts)

  if err != nil {
    return nil, fmt.Errorf("error creating pod watcher for deploy(%s): %s", deployName, err.Error())
  }

  // Return the watcher's channel.
  return watcher.ResultChan(), nil
}

func DeploymentWatcherChannel(client *typedv1beta1.ExtensionsV1beta1Client, nsp string, deployName string) (<-chan watch.Event, error) {
  // Define selector options to only find the pod we just deployed.
  watchOpts := metav1.ListOptions{
    LabelSelector: fmt.Sprintf("app=%s", deployName),
  }

  // Get a namespaced deployment watcher object.
  watcher, err := client.Deployments(nsp).Watch(watchOpts)

  if err != nil {
    return nil, fmt.Errorf("error creating deployment watcher for deploy(%s): %s", deployName, err.Error())
  }

  // Return the watcher's channel.
  return watcher.ResultChan(), nil
}