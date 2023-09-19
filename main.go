package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
	"path/filepath"
)

func main() {
	tests := []struct {
		name      string
		userAgent string
		expect    string
	}{
		{
			name:      "custom",
			userAgent: "test-agent",
			expect:    "test-agent",
		},
	}

	for _, tc := range tests {
		gv := apiv1.SchemeGroupVersion
		config := &rest.Config{
			Host: "485d-18-177-9-140.ngrok-free.app",
		}
		config.GroupVersion = &gv
		//config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
		config.UserAgent = tc.userAgent
		config.ContentType = "application/json"

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println("failed to create REST client: %v", err)
		}

		deploymentsClient := client.AppsV1().Deployments(apiv1.NamespaceDefault)

		fmt.Println(deploymentsClient)

		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "nginx-deployment",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: pointer.Int32Ptr(2),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "nginx",
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  "web",
								Image: "nginx:1.12",
								Ports: []apiv1.ContainerPort{
									//Name:          "http",
									//Protocol:      apiv1.ProtocolTcp,
									//ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		}
		fmt.Println(deployment)
		// create deployment
		fmt.Println("creating or stop deployment nginx...")
		//停止流程
		//err = deploymentsClient.Delete(context.TODO(), "nginx-deployment", metav1.DeleteOptions{})
		//if err != nil {
		//	fmt.Println(err)
		//}

		//create secrets
		var s *apiv1.Secret
		yamlFile, err := ioutil.ReadFile("secret.yaml")
		if err != nil {
			fmt.Println(err.Error())
		} // 将读取的yaml文件解析为响应的 struct

		err = json.Unmarshal(yamlFile, &s)
		if err != nil {
			fmt.Printf("unmarshall error:", err.Error())
			return
		}
		secret, err := client.CoreV1().Secrets("testSecret").Create(context.TODO(), s, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Created Secret %q.\n", secret.GetObjectMeta().GetName())
		}

		//创建流程
		//result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
		//if err != nil {
		//	fmt.Println(err)
		//}
		//fmt.Printf("created deployment %q.\n", result.GetObjectMeta().GetName())
		//
		//_, err = client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
		//if err != nil {
		//	fmt.Println(err)
		//}
	}

	// 加载 Kubernetes 配置
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Dir("config"))
	if err != nil {
		fmt.Println(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(clientset)
	// 获取所有的 Pod
	//pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//// 遍历并输出每个 Pod 的名称和 IP 地址
	//for _, pod := range pods.Items {
	//	fmt.Printf("Pod Name: %s, IP Address: %s\n", pod.Name, pod.Status.PodIP)
	//}
}
