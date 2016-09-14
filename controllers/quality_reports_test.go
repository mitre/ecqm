package controllers

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/matryer/silk/runner"
	"github.com/mitre/ecqm/models"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/dbtest"
)

type QualityReportSuite struct {
	DBServer *dbtest.DBServer
	Engine   *gin.Engine
	Database *mgo.Database
}

var _ = Suite(&QualityReportSuite{})

func TestAPI(t *testing.T) { TestingT(t) }

func (q *QualityReportSuite) SetUpSuite(c *C) {
	//set up dbtest server
	q.DBServer = &dbtest.DBServer{}
	q.DBServer.SetPath(c.MkDir())
	session := q.DBServer.Session()
	q.Database = session.DB("qme-test")
	qr := &models.QualityReport{MeasureID: "efg", EffectiveDate: 5678}
	qrID := bson.ObjectIdHex("56bd06841cd462774f2af485")
	qr.ID = qrID
	q.Database.C("query_cache").Insert(qr)
	ir := models.IndividualResult{MeasureID: "abcd", EffectiveDate: 1234, PatientID: "1234", InitialPatientPopulation: 1, Last: "Jones"}
	id := bson.NewObjectId()
	rw := &models.ResultWrapper{ID: id, IndividualResult: ir}
	q.Database.C("patient_cache").Insert(rw)
	ir2 := models.IndividualResult{MeasureID: "efg", EffectiveDate: 5678, PatientID: "1234", InitialPatientPopulation: 1, Last: "B"}
	id2 := bson.NewObjectId()
	rw2 := &models.ResultWrapper{ID: id2, IndividualResult: ir2}
	q.Database.C("patient_cache").Insert(rw2)
	ir3 := models.IndividualResult{MeasureID: "efg", EffectiveDate: 5678, PatientID: "5678", InitialPatientPopulation: 1, Denominator: 1, Last: "A"}
	id3 := bson.NewObjectId()
	rw3 := &models.ResultWrapper{ID: id3, IndividualResult: ir3}
	q.Database.C("patient_cache").Insert(rw3)
	e := gin.New()
	e.GET("/QualityReport/:id", ShowQualityReportHandler(q.Database))
	e.POST("/QualityReport", CreateQualityReportHandler(q.Database))
	e.GET("/QualityReport/:id/PatientResults", ShowQualityReportPatientsHandler(q.Database))
	q.Engine = e
}

func (q *QualityReportSuite) TestAPI(c *C) {
	s := httptest.NewServer(q.Engine)
	defer s.Close()
	runner.New(c, s.URL).RunGlob(filepath.Glob("../api_doc/*.silk.md"))
}

func (q *QualityReportSuite) TearDownSuite(c *C) {
	q.Database.Session.Close()
	q.DBServer.Wipe()
}
