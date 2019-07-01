package main

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
)

func getNoOfPods(w http.ResponseWriter, r *http.Request) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	w.Write([]byte(fmt.Sprintf("There are %d pods in the cluster\n", len(pods.Items))))
}

func startAppHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Could not parse form.", http.StatusBadRequest)
			return
		}
		imagePath := r.PostForm.Get("imagePath")
		username := r.PostForm.Get("username")
		appName := r.PostForm.Get("appName")
		startApp(imagePath, username, appName)
	}
}

func startApp(imagePath string, username string, appName string) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName + "-" + username,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": appName + "-" + username,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": appName + "-" + username,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  appName + "-" + username,
							Image: imagePath,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 8000,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "COLUMBUS_USERNAME",
									Value: username,
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	servicesClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName + "-" + username,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"name": appName + "-" + username,
			},
			Ports: []apiv1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(8000),
				},
			},
			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}
	serviceResult, err := servicesClient.Create(service)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service %q.\n", serviceResult.GetObjectMeta().GetName())

	ingressClient := clientset.ExtensionsV1beta1().Ingresses(apiv1.NamespaceDefault)
	ingress := &extensions.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName + "-" + username,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                "nginx",
				"nginx.ingress.kubernetes.io/rewrite-target": "/$2",
			},
		},
		Spec: extensions.IngressSpec{
			Rules: []extensions.IngressRule{
				{
					Host: "applications.columbusecosystem.com",
					IngressRuleValue: extensions.IngressRuleValue{
						HTTP: &extensions.HTTPIngressRuleValue{
							Paths: []extensions.HTTPIngressPath{
								{
									Path: "/" + username + "/" + appName + "(/|$)(.*)",
									Backend: extensions.IngressBackend{
										ServiceName: appName + "-" + username,
										ServicePort: intstr.FromInt(8000),
									},
								},
							},
						},
					},
				},
			},
			TLS: []extensions.IngressTLS{
				{
					Hosts: []string{
						"applications.columbusecosystem.com",
					},
					SecretName: "tls-staging-cert",
				},
			},
		},
	}
	ingressResult, err := ingressClient.Create(ingress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service %q.\n", ingressResult.GetObjectMeta().GetName())
}

func int32Ptr(i int32) *int32 { return &i }
