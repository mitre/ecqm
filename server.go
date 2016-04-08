package main

import (
	"flag"
	"fmt"

	"github.com/intervention-engine/fhir/server"
	"github.com/mitre/ecqm/controllers"
	ptmatch "github.com/mitre/ptmatch/server"
	"gopkg.in/mgo.v2"
)

func main() {

	s := server.NewServer("localhost")
	assetPath := flag.String("assets", "", "Path to static assets to host")
	flag.Parse()

	if *assetPath == "" {
		fmt.Println("You must specify a static asset path")
		return
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	db := session.DB("fhir")

	s.Engine.GET("/QualityReport/:id", controllers.ShowQualityReportHandler(db))
	s.Engine.POST("/QualityReport", controllers.CreateQualityReportHandler(db))

	s.Engine.GET("/Measure/:id", controllers.ShowMeasureHandler(db))
	s.Engine.GET("/Measure", controllers.IndexMeasureHandler(db))
	s.Engine.GET("/UserInfo", controllers.UserInfo)

	s.Engine.StaticFile("/", fmt.Sprintf("%s/index.html", *assetPath))
	s.Engine.Static("/assets", fmt.Sprintf("%s/assets", *assetPath))
	ptmatch.Setup(s)

	s.Run(server.Config{})
}
