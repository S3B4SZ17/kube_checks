package healthcheck

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

type K8s struct {
	Client     *kubernetes.Clientset
	Metrics    *metrics.Clientset
	RestConfig *rest.Config
}

func Connection(kubeconfig *string) *K8s {
	if home := homedir.HomeDir(); home != "" && *kubeconfig == "" {
		temp_kubeconfig := filepath.Join(home, ".kube", "config")
		kubeconfig = &temp_kubeconfig
	}
	// use the current context in kubeconfig
	restConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("[Error] Argument passed does not correspond to a valid kubeconfig file")
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
	}

	//create metrics
	metrics, err := metrics.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &K8s{
		Client:     clientset,
		Metrics:    metrics,
		RestConfig: restConfig,
	}
}

func (k *K8s) GetFailingPods() (result string) {
	// namespace := "namespace"
	pods, err := k.Client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range pods.Items {
		if v.Status.Phase != "Running" && v.Status.Phase != "Completed" && v.Status.Phase != "Succeeded" {
			result += fmt.Sprintf("Ns: %v; Pn: %v; St: %v\n", v.Namespace, v.Name, v.Status.Phase)
		}
	}

	return

}

func (k *K8s) GetPodName(partial_name string) (result *corev1.Pod) {
	// namespace := "namespace"
	pods, err := k.Client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range pods.Items {
		if strings.Contains(v.Name, partial_name) {
			result = &v
			break
		}
	}
	if result == nil {
		log.Println("[WARN] Couldnt find pod with name specified = " + partial_name)
		return nil
	}
	return

}

func (k *K8s) Exec(namespace string, pod *corev1.Pod, command []string) string {
	attachOptions := &corev1.PodExecOptions{
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
		Container: pod.Spec.Containers[0].Name,
		Command:   command,
	}

	request := k.Client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(attachOptions, scheme.ParameterCodec)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	streamOptions := remotecommand.StreamOptions{
		Stdout: stdout,
		Stderr: stderr,
	}

	exec, err := remotecommand.NewSPDYExecutor(k.RestConfig, "POST", request.URL())
	if err != nil {
		log.Fatal(err)
	}

	err = exec.Stream(streamOptions)
	if err != nil {
		log.Fatal(err)
	}

	result := strings.TrimSpace(stdout.String()) + "\n" + strings.TrimSpace(stderr.String())
	result = strings.TrimSpace(result)
	return result
}

func (k *K8s) NodeMetrics() (result string) {
	// namespace := "namespace"
	node, err := k.Metrics.MetricsV1alpha1().NodeMetricses().List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range node.Items {
		result = fmt.Sprintf("Memory usage is %v" ,v.Usage)
	}
	if result == "" {
		log.Println("[WARN] An error happened")
		return ""
	}
	return

}