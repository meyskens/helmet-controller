package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/meyskens/helmet/template"
)

type putData struct {
	Values map[string]interface{} `json:"values"`
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
	manifests, _, err := newChart.CreateManifests(template.NewRelease(name, namespace))
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
