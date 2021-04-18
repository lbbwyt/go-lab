package k8s_client

//k8s client-go 客户端二次封装

//client-go对kubernetes资源对象的调用，需要先获取kubernetes的配置信息，即$HOME/.kube/config。

//整个调用的过程如下：

//kubeconfig→rest.config→clientset→具体的client(CoreV1Client)→具体的资源对象(pod)→RESTClient→http.Client→HTTP请求的发送及响应

//通过clientset中不同的client和client中不同资源对象的方法实现对kubernetes中资源对象的增删改查等操作，常用的client有CoreV1Client、AppsV1beta1Client、ExtensionsV1beta1Client等。

//一般二次开发只需要创建deployment、service、ingress三个资源对象即可，
//pod对象由deployment包含的replicaSet来控制创建和删除。
//函数调用的入参一般只有NAMESPACE和kubernetesObject两个参数，部分操作有Options的参数。
//在创建前，需要对资源对象构造数据，可以理解为编辑一个资源对象的yaml文件，然后通过kubectl create -f xxx.yaml来创建对象。
