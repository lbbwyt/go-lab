package k8s_client

import (
	"context"
	"go-lab/pkg/utils/retry"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
)

func (c *K8sClient) GetDeploymentClient(namespace string) v1.DeploymentInterface {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	return c.clientSet.AppsV1().Deployments(namespace)
}

//创建deployment
func (c *K8sClient) CreateDeployment(namespace string, deployment *appsv1.Deployment) error {
	deploymentsClient := c.GetDeploymentClient(namespace)
	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

//更新deployment
func (c *K8sClient) UpdateDeployment(namespace string, d *appsv1.Deployment) error {
	deploymentsClient := c.GetDeploymentClient(namespace)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := deploymentsClient.Update(context.TODO(), d, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return retryErr
	}
	return nil
}

//删除deployment
func (c *K8sClient) DeleteDeployment(namespace, name, delPropagation string) error {
	deploymentsClient := c.GetDeploymentClient(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if delPropagation != "" {
		deletePolicy = metav1.DeletionPropagation(delPropagation)
	}
	if err := deploymentsClient.Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
