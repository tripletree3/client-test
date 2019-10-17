package main

import (
	"flag"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"

	//"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	//v1lister "k8s.io/client-go/listers/core/v1"
	//"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"time"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	sharedInfomerFactory := informers.NewSharedInformerFactory(clientset, 0)

	podLister := sharedInfomerFactory.Core().V1().Pods().Lister()
	//
	stopCh := make(chan struct{})

	defer close(stopCh)

	go sharedInfomerFactory.Start(stopCh)

	podInformer := sharedInfomerFactory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    handlePodAdd,
		UpdateFunc: handlePodUpdate,
		DeleteFunc: handlePodDelete,
	})

	//time.Sleep(time.Second * 11)
	for {
		podList, err := podLister.List(labels.Everything())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(len(podList))
		}
		fmt.Println("sleep 10s")
		time.Sleep(time.Second * 10)
	}
	//pList, err := podLister.Pods("kube-system").List(labels.Nothing())
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(pList)
	//}
	//selector := fields.ParseSelectorOrDie("spec.nodeName!=" + "" + ",status.phase!=" + string(apiv1.PodSucceeded) +
	//	",status.phase!=" + string(apiv1.PodFailed))
	//podListWatch := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", apiv1.NamespaceAll, selector)
	//
	//store, reflector := cache.NewNamespaceKeyedIndexerAndReflector(podListWatch, &apiv1.Pod{}, time.Hour)
	//podLister := v1lister.NewPodLister(store)
	//go reflector.Run(stopCh)
	//
	////time.Sleep(30)
	//
	//for {
	//	pods, _ := podLister.List(labels.Everything())
	//	fmt.Println(len(pods))
	//	fmt.Println("sleep 10s")
	//	time.Sleep(time.Second * 10)
	//}

	//podsList, err := clientset.CoreV1().Pods("kube-system").List(metav1.ListOptions{})
	//fmt.Println(podsList)

}

func handlePodAdd(obj interface{}) {
	newPod := obj.(*apiv1.Pod)
	fmt.Println("add pod: " + newPod.Namespace + ": " + newPod.Name)
}

func handlePodUpdate(oldObj, newObj interface{}) {
	oldPod := oldObj.(*apiv1.Pod)
	newPod := newObj.(*apiv1.Pod)
	fmt.Println("update old pod: " + oldPod.Namespace + ": " + oldPod.Name + " " + string(oldPod.Status.Phase))
	fmt.Println("update to new pod: " + newPod.Namespace + ": " + newPod.Name + " " + string(newPod.Status.Phase))
}

func handlePodDelete(obj interface{}) {
	pod := obj.(*apiv1.Pod)
	fmt.Println("delete pod: " + pod.Namespace + ": " + pod.Name + " " + string(pod.Status.Phase))
}
