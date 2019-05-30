package main

import (
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "default"
	listDeployments(namespace, clientset)
	createDeployment(namespace, clientset)
}

func listDeployments(ns string, cs *kubernetes.Clientset) {
	dList, err := cs.AppsV1().Deployments(ns).List(metav1.ListOptions{})
	if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error get deployment list in namespace %s: %v", ns, statusError)
		return
	}
	fmt.Printf("Deployments %v\n", len(dList.Items))
	for i, item := range dList.Items {
		fmt.Printf("%v: %v %v %v\n", i, item.Name, item.Status.Replicas, item.APIVersion)
	}
}

func createDeployment(ns string, cs *kubernetes.Clientset) {
	var replicas int32 = 2
	var deployment = appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"service": "test"},
			},
		},
	}
	fmt.Printf("%v", deployment)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
