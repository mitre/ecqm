package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/intervention-engine/fhir/auth"
	"github.com/intervention-engine/fhir/server"
	"github.com/mitre/ecqm/controllers"
	"github.com/mitre/ptmatch/middleware"
	ptmatch "github.com/mitre/ptmatch/server"
	"gopkg.in/mgo.v2"
)

func main() {
	assetPath := flag.String("assets", "", "Path to static assets to host")
	dbHostArg := flag.String("mongoHost", "localhost", "MongoDB server address, defaults to localhost")
	jwkPath := flag.String("heartJWK", "", "Path the JWK for the HEART client")
	clientID := flag.String("heartClientID", "", "Client ID registered with the OP")
	opURL := flag.String("heartOP", "", "URL for the OpenID Provider")
	sessionSecret := flag.String("secret", "", "Secret for the cookie session")
	flag.Parse()

	s := server.NewServer(*dbHostArg)
	session, err := mgo.Dial(*dbHostArg)
	if err != nil {
		panic(err)
	}
	db := session.DB("fhir")

	var authConfig auth.Config

	if *jwkPath != "" {
		if *clientID == "" || *opURL == "" {
			fmt.Println("You must provide both a client ID and OP URL for HEART mode")
			return
		}
		secret := *sessionSecret
		if secret == "" {
			secret = "reallySekret"
		}
		authConfig = auth.HEART(*clientID, *jwkPath, *opURL, secret)
	} else {
		authConfig = auth.None()
	}

	ar := func(e *gin.Engine) {
		e.GET("/QualityReport/:id", controllers.ShowQualityReportHandler(db))
		e.POST("/QualityReport", controllers.CreateQualityReportHandler(db))
		e.GET("/PatientReport/:id", controllers.ShowIndividualResultsForPatientHandler(db))
		e.GET("/QualityReport/:id/PatientResults", controllers.ShowQualityReportPatientsHandler(db))

		s.Engine.GET("/CQMMeasure/:id", controllers.ShowMeasureHandler(db))
		e.GET("/CQMMeasure", controllers.IndexMeasureHandler(db))
		e.GET("/UserInfo", controllers.UserInfo)

		ptmatch.Setup(e)

		if *assetPath != "" {
			e.StaticFile("/", fmt.Sprintf("%s/index.html", *assetPath))
			e.Static("/assets", fmt.Sprintf("%s/assets", *assetPath))
		}
	}
	recMatchWatch := middleware.PostProcessRecordMatchResponse()
	s.AddMiddleware("Bundle", recMatchWatch)

	s.AfterRoutes = append(s.AfterRoutes, ar)

	s.Run(server.Config{Auth: authConfig, ServerURL: "http://localhost:3001", DatabaseName: "fhir"})
}
