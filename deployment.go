package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
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

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	//for {
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// Examples for error handling:
	// - Use helper functions like e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	namespace := "default"
	pod := "myapp-pod1"
	_, err = clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		//fmt.Printf("%v\n", cPod)
	}

	nList, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	fmt.Printf("Nodes %v\n", len(nList.Items))
	for i, item := range nList.Items {
		fmt.Printf("%v: %v %v\n", i, item.Name, item.Status.Capacity)
	}

	sList, err := clientset.CoreV1().Services(namespace).List(metav1.ListOptions{})
	fmt.Printf("Services %v\n", len(sList.Items))
	for i, item := range sList.Items {
		fmt.Printf("%v: %v %v\n", i, item.Name, item.Spec.ClusterIP)
	}

	dList, err := clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	fmt.Printf("Deployments %v\n", len(dList.Items))
	for i, item := range dList.Items {
		fmt.Printf("%v: %v %v %v\n", i, item.Name, item.Status.Replicas, *item.Spec.Replicas)
	}

	//	time.Sleep(10 * time.Second)
	//}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
