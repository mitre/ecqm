package models

import (
	"github.com/pebbe/util"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/dbtest"
)

type IndividualResultSuite struct {
	DBServer *dbtest.DBServer
	Database *mgo.Database
}

var _ = Suite(&IndividualResultSuite{})

func (i *IndividualResultSuite) SetUpSuite(c *C) {
	//set up dbtest server
	i.DBServer = &dbtest.DBServer{}
	i.DBServer.SetPath(c.MkDir())
}

func (i *IndividualResultSuite) SetUpTest(c *C) {
	session := i.DBServer.Session()
	i.Database = session.DB("qme-test")
}

func (i *IndividualResultSuite) TearDownTest(c *C) {
	i.Database.Session.Close()
	i.DBServer.Wipe()
}

func (i *IndividualResultSuite) TestFindAndPopulateIndividualResult(c *C) {
	ir := &IndividualResult{MeasureID: "abcd", EffectiveDate: 1234, PatientID: "1234", InitialPatientPopulation: 5}
	id := bson.NewObjectId()
	rw := &ResultWrapper{ID: id, IndividualResult: ir}
	i.Database.C("patient_cache").Insert(rw)
	irToFind := &IndividualResult{MeasureID: "abcd", EffectiveDate: 1234, PatientID: "1234"}
	exists, err := FindAndPopulateIndividualResult(i.Database, irToFind)
	util.CheckErr(err)
	c.Assert(exists, Equals, true)
	c.Assert(irToFind.InitialPatientPopulation, Equals, 5)
	irDoesntExist := &IndividualResult{MeasureID: "foobar", EffectiveDate: 1234, PatientID: "1234"}
	exists, err = FindAndPopulateIndividualResult(i.Database, irDoesntExist)
	util.CheckErr(err)
	c.Assert(exists, Equals, false)
}
