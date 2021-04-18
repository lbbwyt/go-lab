package k8s

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//k8s Service 简单二次封装
type KubernetesService struct {
	//metadata
	Name      string            `json:"name"`
	Labels    map[string]string `json:"match_labels"`
	Namespace string            `json:"namespace"`

	//spec:  # 资源规范字段
	Port       int32          `json:"pory"`
	TargetPort int32          `json:"target_port"`
	Protocol   apiv1.Protocol `json:"protocol"`

	Type     apiv1.ServiceType `json:"type"`
	Selector map[string]string `json:"selector"`
}

func (k *KubernetesService) BindPort(port int32) *KubernetesService {
	k.Port = port
	return k
}

func (k *KubernetesService) BindTargetPort(targetPort int32) *KubernetesService {
	k.TargetPort = targetPort
	return k
}

func (k *KubernetesService) BindProtocol(Protocol string) *KubernetesService {
	k.Protocol = apiv1.Protocol(Protocol)
	return k
}

func (k *KubernetesService) BindType(sType string) *KubernetesService {
	k.Type = apiv1.ServiceType(sType)
	return k
}

func (k *KubernetesService) BindName(name string) *KubernetesService {
	k.Name = name
	return k
}

func (k *KubernetesService) BindNamespace(namespace string) *KubernetesService {
	k.Namespace = namespace
	return k
}

func (k *KubernetesService) BindSelector(key, value string) *KubernetesService {
	if k.Selector != nil {
		k.Selector[key] = value
	} else {
		m := make(map[string]string, 0)
		m[key] = value
		k.Selector = m
	}
	return k
}

func (k *KubernetesService) BindLabels(key, value string) *KubernetesService {
	if k.Labels != nil {
		k.Labels[key] = value
	} else {
		m := make(map[string]string, 0)
		m[key] = value
		k.Labels = m
	}
	return k
}

func (k *KubernetesService) Build() *apiv1.Service {
	servicePorts := make([]apiv1.ServicePort, 0)
	if k.Protocol == "" {
		k.Protocol = apiv1.ProtocolTCP
	}
	if k.Type == "" {
		k.Type = apiv1.ServiceTypeLoadBalancer
	}
	servicePorts = append(servicePorts, apiv1.ServicePort{
		Protocol:   k.Protocol,
		Port:       k.Port,
		TargetPort: intstr.FromInt(int(k.TargetPort)),
	})
	if k.Namespace == "" {
		k.Namespace = apiv1.NamespaceDefault
	}

	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      k.Name,
			Namespace: k.Namespace,
			Labels:    k.Labels,
		},
		Spec: apiv1.ServiceSpec{
			Ports:    servicePorts,
			Type:     k.Type,
			Selector: k.Selector,
		},
	}
}
