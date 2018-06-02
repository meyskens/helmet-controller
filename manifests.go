package main

import (
	"errors"
	"fmt"
	"log"

	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	ext "k8s.io/api/extensions/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

func applyManifest(namespace string, in []byte) error {
	client, err := newKubernetesClient()
	if err != nil {
		return err
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(in, nil, nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
		return errors.New("Error parsing YAML")
	}

	// now use switch over the type of the object
	// and match each type-case
	// more to be added soon
	switch o := obj.(type) {
	case *apps.Deployment:
		_, err := client.AppsV1().Deployments(namespace).Get(o.Name, meta.GetOptions{})
		if err != nil {
			// need to create
			_, err = client.AppsV1().Deployments(namespace).Create(o)
		} else {
			_, err = client.AppsV1().Deployments(namespace).Update(o)
		}
		log.Printf("Error %s in:\n%s\n", err, string(in))
		return err
	case *ext.Ingress:
		_, err := client.ExtensionsV1beta1().Ingresses(namespace).Get(o.Name, meta.GetOptions{})
		if err != nil {
			// need to create
			_, err = client.ExtensionsV1beta1().Ingresses(namespace).Create(o)
		} else {
			_, err = client.ExtensionsV1beta1().Ingresses(namespace).Update(o)
		}
		log.Printf("Error %s in:\n%s\n", err, string(in))
		return err
	case *v1.Service:
		_, err := client.CoreV1().Services(namespace).Get(o.Name, meta.GetOptions{})
		if err == nil {
			// there is a bug here that doesn't allow updates
			err = client.CoreV1().Services(namespace).Delete(o.Name, &meta.DeleteOptions{})
		}
		_, err = client.CoreV1().Services(namespace).Create(o)
		log.Printf("Error %s in:\n%s\n", err, string(in))
		return err
	case *batch.Job:
		_, err := client.BatchV1().Jobs(namespace).Get(o.Name, meta.GetOptions{})
		if err != nil {
			// need to create
			_, err = client.BatchV1().Jobs(namespace).Create(o)
		} else {
			_, err = client.BatchV1().Jobs(namespace).Update(o)
		}
		log.Printf("Error %s in:\n%s\n", err, string(in))
		return err
	default:
		//o is unknown for us
		log.Printf("Unknown type for:\n%s\n", string(in))
		return errors.New("Unknown type")
	}
}

func deleteManifest(namespace string, in []byte) error {
	client, err := newKubernetesClient()
	if err != nil {
		return err
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(in, nil, nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
	}

	// now use switch over the type of the object
	// and match each type-case
	// more to be added soon
	switch o := obj.(type) {
	case *apps.Deployment:
		err = client.AppsV1().Deployments(namespace).Delete(o.Name, &meta.DeleteOptions{})
		break
	case *ext.Ingress:
		err = client.ExtensionsV1beta1().Ingresses(namespace).Delete(o.Name, &meta.DeleteOptions{})
		break
	case *v1.Service:
		err = client.CoreV1().Services(namespace).Delete(o.Name, &meta.DeleteOptions{})
		break
	default:
		//o is unknown for us
		return errors.New("Unknown type")
	}
	return err
}
