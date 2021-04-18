package k8s_client

import (
	"context"
	"flag"
	"fmt"
	"go-lab/model/k8s"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	"testing"
	"time"
)

func TestGetClientSet(t *testing.T) {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	clientSet, err := GetClientSet(*kubeconfig)
	if err != nil {
		panic(err)
	}

	for {
		pods, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientSet.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod example-xxxxx not found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found example-xxxxx pod in default namespace\n")
		}
		time.Sleep(10 * time.Second)
	}
}

func TestK8sClient_GetPod(t *testing.T) {
	c, err := NewK8sClient(getLocalConfig())
	if err != nil {
		panic(err)
	}
	pod, err := c.GetPod("default", "goweb")
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%v", pod))
}

func TestK8sClient_CreateDeployment(t *testing.T) {
	c, err := NewK8sClient("")
	if err != nil {
		panic(err)
	}
	labels := map[string]string{
		"app": "gin-web-test",
	}
	deployment := k8s.KubernetesDeployment{
		Name:                 "gin-web-test",
		Labels:               labels,
		Namespace:            "",
		Replicas:             1,
		RevisionHistoryLimit: 3,
		MatchLabels:          labels,
		ImageName:            "gin-web-test",
		Image:                "localhost:5000/gin_swagger",
		ImagePullPolicy:      "",
		RestartPolicy:        "",
	}
	err = c.CreateDeployment("", deployment.Build())
	if err != nil {
		panic(err)
	}
}

func TestK8sClient_CreateKubeService(t *testing.T) {
	c, err := NewK8sClient("")
	if err != nil {
		panic(err)
	}
	labels := map[string]string{
		"app": "gin-web-test",
	}

	Selector := map[string]string{
		"app": "gin-web-test",
	}

	service := k8s.KubernetesService{
		Name:       "gin-web-test",
		Labels:     labels,
		Namespace:  "",
		Port:       8001,
		TargetPort: 8888,
		Protocol:   "",
		Type:       "",
		Selector:   Selector,
	}
	err = c.CreateKubeService("", service.Build())
	if err != nil {
		panic(err)
	}

}

func TestK8sClient_DeleteKubeService(t *testing.T) {
	c, err := NewK8sClient("")
	if err != nil {
		panic(err)
	}
	err = c.DeleteKubeService("", "nginx", "")
	if err != nil {
		panic(err)
	}
	err = c.DeleteDeployment("", "nignx-deployment", "")
	if err != nil {
		panic(err)
	}
}
