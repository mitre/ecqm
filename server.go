package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
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

	e := gin.Default()
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	db := session.DB("fhir")

	e.GET("/QualityReport/:id", controllers.ShowQualityReportHandler(db))
	e.POST("/QualityReport", controllers.CreateQualityReportHandler(db))

	e.GET("/Measure/:id", controllers.ShowMeasureHandler(db))
	e.GET("/Measure", controllers.IndexMeasureHandler(db))
	e.GET("/UserInfo", controllers.UserInfo)

	e.StaticFile("/", fmt.Sprintf("%s/index.html", *assetPath))
	e.Static("/assets", fmt.Sprintf("%s/assets", *assetPath))
	e.Run(":3001")
}
