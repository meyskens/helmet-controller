package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/labstack/echo"
	"github.com/meyskens/helmet/template"
)

type putData struct {
	Values map[string]interface{} `json:"values"`
}

func mapStringInterfaceToMapInterfaceInterace(in map[string]interface{}) map[interface{}]interface{} {
	m := map[interface{}]interface{}{}
	for k, v := range in {
		m[k] = v
		switch v.(type) {
		case map[string]interface{}:
			m[k] = mapStringInterfaceToMapInterfaceInterace(v.(map[string]interface{}))
			break
		default:
			m[k] = v
		}
	}

	return m
}

func putDeployment(c echo.Context) error {
	data := putData{}
	name := c.Param("name")
	err := c.Bind(&data)
	if err != nil {
		log.Println(err)
	}

	manifests, err := buildManifests(c, data, name)
	if err != nil {
		return err
	}

	for _, manifest := range manifests {
		err = applyManifest(namespace, manifest)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"result": "success"})
}

func buildManifests(c echo.Context, data putData, name string) ([][]byte, error) {
	newChart, err := chart.Clone()
	if err != nil {
		return nil, c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	spew.Dump(mapStringInterfaceToMapInterfaceInterace(data.Values))
	newChart.MergeValues(mapStringInterfaceToMapInterfaceInterace(data.Values))
	files, _, err := newChart.CreateManifests(template.NewRelease(name, namespace))
	if err != nil {
		log.Println(err)
		return nil, c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	manifests := [][]byte{}
	for _, file := range files {
		for _, yamlFile := range strings.Split(string(file), "---") {
			if strings.TrimSpace(yamlFile) != "" {
				manifests = append(manifests, []byte(yamlFile))
				log.Println(yamlFile)
			}
		}
	}
	return manifests, nil
}

func deleteDeployment(c echo.Context) error {
	data := putData{}
	name := c.Param("name")
	c.Bind(&data)

	newChart, err := chart.Clone()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	//newChart.MergeValues(data.Values)
	manifests, _, err := newChart.CreateManifests(template.NewRelease(name, namespace))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	for _, manifest := range manifests {
		err = deleteManifest(namespace, manifest)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"result": "success"})
}
