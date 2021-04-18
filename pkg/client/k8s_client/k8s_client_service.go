package k8s_client

import (
	"context"
	"go-lab/pkg/utils/retry"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func (c *K8sClient) GetKubeServiceClient(namespace string) v12.ServiceInterface {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	return c.clientSet.CoreV1().Services(namespace)
}

//创建kubernetes service
func (c *K8sClient) CreateKubeService(namespace string, service *apiv1.Service) error {
	kubeServiceClient := c.GetKubeServiceClient(namespace)
	_, err := kubeServiceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

//更新deployment
func (c *K8sClient) UpdateKuberService(namespace string, service *apiv1.Service) error {
	kubeServiceClient := c.GetKubeServiceClient(namespace)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := kubeServiceClient.Update(context.TODO(), service, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return retryErr
	}
	return nil
}

//删除Service
func (c *K8sClient) DeleteKubeService(namespace, name, delPropagation string) error {
	deletePolicy := metav1.DeletePropagationForeground
	if delPropagation != "" {
		deletePolicy = metav1.DeletionPropagation(delPropagation)
	}
	kubeServiceClient := c.GetKubeServiceClient(namespace)

	if err := kubeServiceClient.Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}
