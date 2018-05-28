package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var cs *kubernetes.Clientset

func newKubernetesClient() (*kubernetes.Clientset, error) {
	if cs != nil {
		return cs, nil
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		cs = client
	}

	return client, err
}
