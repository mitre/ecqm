package main

import (
	"flag"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mitre/ecqm/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
	assetPath := flag.String("assets", "", "Path to static assets to host")
	flag.Parse()

	if *assetPath == "" {
		fmt.Println("You must specify a static asset path")
		return
	}

	e := echo.New()
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	db := session.DB("fhir")

	e.Get("/QualityReport/:id", controllers.ShowQualityReportHandler(db))
	e.Post("/QualityReport", controllers.CreateQualityReportHandler(db))

	e.Get("/Measure/:id", controllers.ShowMeasureHandler(db))
	e.Get("/Measure", controllers.IndexMeasureHandler(db))
	e.Get("/UserInfo", controllers.UserInfo)

	e.Index(fmt.Sprintf("%s/index.html", *assetPath))
	e.Static("/assets", fmt.Sprintf("%s/assets", *assetPath))
	e.Use(middleware.Logger())
	e.Run(":3001")
}
