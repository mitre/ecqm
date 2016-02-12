package main

import (
	"github.com/labstack/echo"
	"github.com/mitre/ecqm/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
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
	e.Run(":3001")
}
