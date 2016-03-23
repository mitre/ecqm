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
	id := bson.ObjectIdHex("56bd06841cd462774f2af485")
	qr.ID = id
	q.Database.C("query_cache").Insert(qr)
	e := gin.New()
	e.GET("/QualityReport/:id", ShowQualityReportHandler(q.Database))
	e.POST("/QualityReport", CreateQualityReportHandler(q.Database))
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
