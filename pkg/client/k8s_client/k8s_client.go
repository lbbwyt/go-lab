package k8s_client

import (
	"context"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

//kubernetes api 封装
type K8sClient struct {
	clientSet *kubernetes.Clientset
}

func NewK8sClient(configPath string) (*K8sClient, error) {
	clientSet, err := GetClientSet(configPath)
	if err != nil {
		return nil, err
	}
	return &K8sClient{
		clientSet: clientSet,
	}, nil
}

func (c *K8sClient) GetPod(namespace, podName string) (*v1.Pod, error) {
	return c.clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
}

func GetClientSet(configPath string) (*kubernetes.Clientset, error) {
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
