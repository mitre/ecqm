package models

import (
	"testing"

	"github.com/pebbe/util"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/dbtest"
)

type QualityReportSuite struct {
	DBServer *dbtest.DBServer
	Database *mgo.Database
}

var _ = Suite(&QualityReportSuite{})

func Test(t *testing.T) { TestingT(t) }

func (q *QualityReportSuite) SetUpSuite(c *C) {
	//set up dbtest server
	q.DBServer = &dbtest.DBServer{}
	q.DBServer.SetPath(c.MkDir())
}

func (q *QualityReportSuite) SetUpTest(c *C) {
	session := q.DBServer.Session()
	q.Database = session.DB("qme-test")
}

func (q *QualityReportSuite) TearDownTest(c *C) {
	q.Database.Session.Close()
	q.DBServer.Wipe()
}

func (q *QualityReportSuite) TestFindQualityAndPopulateQualityReport(c *C) {
	qr := &QualityReport{MeasureID: "abcd", EffectiveDate: 1234, NPI: "efg"}
	id := bson.NewObjectId()
	qr.ID = id
	q.Database.C("query_cache").Insert(qr)
	qrToFind := &QualityReport{MeasureID: "abcd", EffectiveDate: 1234}
	exists, err := FindQualityAndPopulateQualityReport(q.Database, qrToFind)
	util.CheckErr(err)
	c.Assert(exists, Equals, true)
	c.Assert(qrToFind.NPI, Equals, "efg")
	qrDoesntExist := &QualityReport{MeasureID: "foobar", EffectiveDate: 1234}
	exists, err = FindQualityAndPopulateQualityReport(q.Database, qrDoesntExist)
	util.CheckErr(err)
	c.Assert(exists, Equals, false)
}

func (q *QualityReportSuite) TestFindOrCreateQualityReport(c *C) {
	qr := &QualityReport{MeasureID: "abcd", EffectiveDate: 1234, NPI: "efg"}
	id := bson.NewObjectId()
	qr.ID = id
	q.Database.C("query_cache").Insert(qr)
	qrToFind := &QualityReport{MeasureID: "abcd", EffectiveDate: 1234}
	err := FindOrCreateQualityReport(q.Database, qrToFind)
	util.CheckErr(err)
	c.Assert(qrToFind.NPI, Equals, "efg")
	qrDoesntExist := &QualityReport{MeasureID: "foobar", EffectiveDate: 1234}
	err = FindOrCreateQualityReport(q.Database, qrDoesntExist)
	util.CheckErr(err)
	count, _ := q.Database.C("query_cache").Count()
	c.Assert(count, Equals, 2)
}
