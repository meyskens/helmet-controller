package main

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var cs *kubernetes.Clientset

func newKubernetesClient() (*kubernetes.Clientset, error) {
	if cs != nil {
		return cs, nil
	}
	var config *rest.Config
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("USE_KUBE_CONFIG") == "" {
		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		log.Info("Not in cluster, using .kube/config")
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		kubeconfig := filepath.Join(home, ".kube", "config")
		// use the current context in kubeconfig
		var err error
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		cs = client
	}

	return client, err
}
