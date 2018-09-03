package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/meyskens/helmet/template"
)

type putData struct {
	Values map[interface{}]interface{} `json:"values"`
}

func putDeployment(c echo.Context) error {
	data := putData{}
	name := c.Param("name")
	c.Bind(&data)

	newChart, err := chart.Clone()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	newChart.MergeValues(data.Values)
	files, _, err := newChart.CreateManifests(template.NewRelease(name, namespace))

	manifests := [][]byte{}
	for _, file := range files {
		for _, yamlFile := range strings.Split(string(file), "---") {
			if strings.TrimSpace(yamlFile) != "" {
				manifests = append(manifests, []byte(yamlFile))
			}
		}
	}

	for _, manifest := range manifests {
		err = applyManifest(namespace, manifest)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"result": "success"})
}

func deleteDeployment(c echo.Context) error {
	data := putData{}
	name := c.Param("name")
	c.Bind(&data)

	newChart, err := chart.Clone()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	newChart.MergeValues(data.Values)
	manifests, _, err := newChart.CreateManifests(template.NewRelease(name, namespace))
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
