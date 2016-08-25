package main

import (
	"flag"
	"fmt"

	"github.com/intervention-engine/fhir/server"
	"github.com/mitre/ecqm/controllers"
	"github.com/mitre/heart"
	ptmatch "github.com/mitre/ptmatch/server"
	"gopkg.in/mgo.v2"
)

func main() {

	s := server.NewServer("localhost")
	assetPath := flag.String("assets", "", "Path to static assets to host")
	jwkPath := flag.String("heartJWK", "", "Path the JWK for the HEART client")
	clientID := flag.String("heartClientID", "", "Client ID registered with the OP")
	opURL := flag.String("heartOP", "", "URL for the OpenID Provider")
	sessionSecret := flag.String("secret", "", "Secret for the cookie session")
	flag.Parse()

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	db := session.DB("fhir")

	if *jwkPath != "" {
		if *clientID == "" || *opURL == "" {
			fmt.Println("You must provide both a client ID and OP URL for HEART mode")
			return
		}
		secret := *sessionSecret
		if secret == "" {
			secret = "reallySekret"
		}
		heart.SetUpRoutes(*jwkPath, *clientID, *opURL,
			"http://localhost:3001", secret, s.Engine)
	}

	s.Engine.GET("/QualityReport/:id", controllers.ShowQualityReportHandler(db))
	s.Engine.POST("/QualityReport", controllers.CreateQualityReportHandler(db))
	s.Engine.GET("/PatientReport/:id", controllers.ShowIndividualResultsForPatientHandler(db))

	s.Engine.GET("/Measure/:id", controllers.ShowMeasureHandler(db))
	s.Engine.GET("/Measure", controllers.IndexMeasureHandler(db))
	s.Engine.GET("/UserInfo", controllers.UserInfo)

	if *assetPath != "" {
		s.Engine.StaticFile("/", fmt.Sprintf("%s/index.html", *assetPath))
		s.Engine.Static("/assets", fmt.Sprintf("%s/assets", *assetPath))
	}

	ptmatch.Setup(s)

	s.Run(server.Config{})
}
