package main

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
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
			TypeMeta: metav1.TypeMeta{
				Kind:       "Development",
				APIVersion: "apps/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test1", //什么名字？
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: pointer.Int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{ //这个标记什么？
						"app": "broker",
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{ //这个标记什么？
							"app": "broker",
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  "broker",
								Image: "docker.io/library/broker:v1", //策略image？
								//Command: ["/app/broker",],
								Env: []apiv1.EnvVar{
									{Name: "DualSidePosition", Value: "false"},
									{Name: "GapWithBestOfferPct", Value: "0.0001"},
								},
								EnvFrom: []apiv1.EnvFromSource{
									{SecretRef: &apiv1.SecretEnvSource{
										LocalObjectReference: apiv1.LocalObjectReference{
											Name: "maxpy",
										}}, ConfigMapRef: &apiv1.ConfigMapEnvSource{
										LocalObjectReference: apiv1.LocalObjectReference{
											Name: "alpha-env",
										}},
									},
								},
							},
						},
					},
				},
			}} //end deployment
		fmt.Println(deployment)

		//创建流程
		result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("created deployment %q.\n", result.GetObjectMeta().GetName())

	}
	// create deployment

	//停止流程
	//err = deploymentsClient.Delete(context.TODO(), "nginx-deployment", metav1.DeleteOptions{})
	//if err != nil {
	//	fmt.Println(err)
	//}

	//create secrets
	//var s *apiv1.Secret
	//yamlFile, err := ioutil.ReadFile("secret.yaml")
	//if err != nil {
	//	fmt.Println(err.Error())
	//} // 将读取的yaml文件解析为响应的 struct
	//
	//err = json.Unmarshal(yamlFile, &s)
	//if err != nil {
	//	fmt.Printf("unmarshall error:", err.Error())
	//	return
	//}
	var secret apiv1.Secret
	secret.Kind = "Secret"
	secret.APIVersion = "v1"
	secret.ObjectMeta = metav1.ObjectMeta{
		Name:      "test1",
		Namespace: "default",
	}
	dataMap := make(map[string][]byte)
	dataMap["API_KEY"] = []byte("API_KEY_Value")
	dataMap["API_SECRET"] = []byte("API_SECRET_Value")
	dataMap["Account_Name"] = []byte("Account_Name_Value")
	secret.Data = dataMap
	fmt.Println(secret)

	//result, err := client.CoreV1().Secrets("default").Create(context.TODO(), &secret, metav1.CreateOptions{})
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Printf("Created Secret %q.\n", secret.GetObjectMeta().GetName())
	//}
	//fmt.Printf(" * secret create: %v\n", result)

	//create secrets
	//secretClient := client.CoreV1().Secrets("maxbroker")
	//secretList, _ := secretClient.List(context.TODO(), metav1.ListOptions{})
	//for _, d := range secretList.Items {
	//	fmt.Printf(" * secret: %v %v\n", d.Name, d.Type)
	//}
	//
	//var secret apiv1.Secret
	//secret.Kind = "Secret"
	//secret.APIVersion = "v1"
	//secret.ObjectMeta = metav1.ObjectMeta{
	//	Name:      "cloud-config",
	//	Namespace: "kube-system",
	//}
	//dataMap := make(map[string][]byte)
	//dataMap["cloud.conf"] = []byte("W0dsb2JhbF0KdXNlcm5hbWUgPSBhZG1pbgpwYXNzd29yZCA9IEFkbWluX1BXRF84NjQ4NjczNTFxc2Myd2R2M2VmYjRyZ24KdGVuYW50LWlkID0gOWEyY2I0MDlmMGRhNDhlMzg1ODY4ZjI3ZmM5YzhjOWIKZG9tYWluLW5hbWUgPSBEZWZhdWx0CmF1dGgtdXJsID0gaHR0cDovL29wZW5zdGFjay1rZXlzdG9uZS12aXA6MzUzNTcvdjMKcmVnaW9uID0gcmVnaW9ub25lCg==")
	//secret.Data = dataMap
	//fmt.Println(secret)
	//result, err := secretClient.Create(context.TODO(), &secret, metav1.CreateOptions{})
	//if err != nil {
	//	fmt.Printf("err ... %v", err)
	//}
	//fmt.Printf(" * secret create: %v\n", result)

	//
	//_, err = client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	//if err != nil {
	//	fmt.Println(err)
	//}
}

//// 加载 Kubernetes 配置
//config, err := clientcmd.BuildConfigFromFlags("", filepath.Dir("config"))
//if err != nil {
//fmt.Println(err.Error())
//}
//clientset, err := kubernetes.NewForConfig(config)
//if err != nil {
//fmt.Println(err.Error())
//}
//fmt.Println(clientset)
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
//}
