package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/meyskens/helmet/template"
)

var chart *template.Chart
var namespace string

func main() {
	log.Println("Helmet Controller")
	err := loadFiles()
	if err != nil {
		log.Fatal(err)
	}

	namespace = os.Getenv("NAMESPACE")
	if namespace == "" {
		log.Fatal("NAMESPACE is not defined")
	}

	e := echo.New()
	e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		correctKey := os.Getenv("AUTH_TOKEN")
		if correctKey == "" {
			return false, errors.New("Token not set")
		}
		return key == "valid-key", nil
	}))

	e.GET("/", getRoot)
	e.PUT("/deployment/:name", putDeployment)
	e.DELETE("/deployment/:name", deleteDeployment)

	e.Logger.Fatal(e.Start(":8080"))
}

func getRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Helmet Controller")
}

func loadFiles() error {
	chartPath := os.Getenv("CHART_PATH")
	if chartPath == "" {
		chartPath = "./chart"
	}
	// read in values
	chartFile, err := ioutil.ReadFile(path.Join(chartPath, "Chart.yaml"))
	if err != nil {
		return err
	}
	valuesFile, err := ioutil.ReadFile(path.Join(chartPath, "values.yaml"))
	if err != nil {
		return err
	}

	// parse values
	chartInfo, err := template.ParseChartFile(chartFile)
	if err != nil {
		return err
	}

	values, err := template.ParseValuesFile(valuesFile)
	if err != nil {
		return err
	}

	// read templates
	files, err := ioutil.ReadDir(path.Join(chartPath, "templates"))
	if err != nil {
		return err
	}

	filesMap := map[string][]byte{}
	for _, f := range files {
		content, err := ioutil.ReadFile(path.Join(chartPath, "templates", f.Name()))
		if err != nil {
			return err
		}
		filesMap[f.Name()] = content
	}

	chart = template.New(chartInfo, values, filesMap)

	return nil
}
