package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//k8s deployment 简单二次封装
type KubernetesDeployment struct {
	//metadata
	Name      string            `json:"name"`
	Labels    map[string]string `json:"match_labels"`
	Namespace string            `json:"namespace"`

	//spec:  # 资源规范字段
	Replicas             int32             `json:"replicas"`
	RevisionHistoryLimit int32             `json:"revision_history_limit"` // 保留历史版本
	MatchLabels          map[string]string `json:"match_labels"`           // 匹配标签

	//containers
	Image           string              `json:"image"`             // 容器使用的镜像地址
	ImagePullPolicy apiv1.PullPolicy    `json:"image_pull_policy"` // 每次Pod启动拉取镜像策略，三个选择 Always、Never、IfNotPresent
	RestartPolicy   apiv1.RestartPolicy `json:"restart_policy"`
}

func (k *KubernetesDeployment) BindName(name string) *KubernetesDeployment {
	k.Name = name
	return k
}

func (k *KubernetesDeployment) BindNamespace(namespace string) *KubernetesDeployment {
	k.Namespace = namespace
	return k
}

func (k *KubernetesDeployment) BindLabels(key, value string) *KubernetesDeployment {
	if k.Labels != nil {
		k.Labels[key] = value
	} else {
		m := make(map[string]string, 0)
		m[key] = value
		k.Labels = m
	}
	return k
}

func (k *KubernetesDeployment) BindReplicas(replicas int32) *KubernetesDeployment {
	k.Replicas = replicas
	return k
}

func (k *KubernetesDeployment) BindRevisionHistoryLimit(revisionHistoryLimit int32) *KubernetesDeployment {
	k.RevisionHistoryLimit = revisionHistoryLimit
	return k
}

func (k *KubernetesDeployment) BindMatchLabels(key, value string) *KubernetesDeployment {
	if k.MatchLabels != nil {
		k.MatchLabels[key] = value
	} else {
		m := make(map[string]string, 0)
		m[key] = value
		k.MatchLabels = m
	}
	return k
}

func (k *KubernetesDeployment) BindImage(image string) *KubernetesDeployment {
	k.Image = image
	return k
}

func (k *KubernetesDeployment) BindImagePullPolicy(imagePullPolicy string) *KubernetesDeployment {
	if imagePullPolicy == "" {
		k.ImagePullPolicy = apiv1.PullAlways
	}
	k.ImagePullPolicy = apiv1.PullPolicy(imagePullPolicy)
	return k
}

func (k *KubernetesDeployment) BindRestartPolicy(restartPolicy string) *KubernetesDeployment {
	if restartPolicy == "" {
		k.RestartPolicy = apiv1.RestartPolicyAlways
	}
	k.RestartPolicy = apiv1.RestartPolicy(restartPolicy)
	return k
}

func (k *KubernetesDeployment) Build() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      k.Name,
			Namespace: k.Namespace,
			Labels:    k.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(k.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: k.MatchLabels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: k.Labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Image:           k.Image,
							ImagePullPolicy: k.ImagePullPolicy,
						},
					},
					RestartPolicy: k.RestartPolicy,
				},
			},
		},
	}
}

func int32Ptr(i int32) *int32 { return &i }
