package main

import (
	"errors"
	"fmt"
	"log"

	apps "k8s.io/api/apps/v1"
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
	}

	// now use switch over the type of the object
	// and match each type-case
	// more to be added soon
	switch o := obj.(type) {
	case *apps.Deployment:
		_, err := client.Apps().Deployments(namespace).Get(o.Name, meta.GetOptions{})
		if err != nil {
			// need to create
			_, err = client.Apps().Deployments(namespace).Create(o)
		} else {
			_, err = client.Apps().Deployments(namespace).Update(o)
		}
		return err
		break
	case *ext.Ingress:
		_, err := client.ExtensionsV1beta1().Ingresses(namespace).Get(o.Name, meta.GetOptions{})
		if err != nil {
			// need to create
			_, err = client.ExtensionsV1beta1().Ingresses(namespace).Create(o)
		} else {
			_, err = client.ExtensionsV1beta1().Ingresses(namespace).Update(o)
		}
		return err
		break
	case *v1.Service:
		_, err := client.CoreV1().Services(namespace).Get(o.Name, meta.GetOptions{})
		if err != nil {
			// need to create
			_, err = client.CoreV1().Services(namespace).Create(o)
		} else {
			_, err = client.CoreV1().Services(namespace).Update(o)
		}
		return err
		break
	default:
		//o is unknown for us
		return errors.New("Unknown type")
	}
	return nil
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
		err = client.Apps().Deployments(namespace).Delete(o.Name, &meta.DeleteOptions{})
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
