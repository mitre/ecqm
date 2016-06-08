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

func (i *IndividualResultSuite) TestFindAllResultsForPatient(c *C) {
	ir := IndividualResult{MeasureID: "abcd", EffectiveDate: 1234, PatientID: "1234", InitialPatientPopulation: 1}
	id := bson.NewObjectId()
	rw := &ResultWrapper{ID: id, IndividualResult: ir}
	i.Database.C("patient_cache").Insert(rw)
	ir2 := IndividualResult{MeasureID: "efgh", EffectiveDate: 1234, PatientID: "1234", InitialPatientPopulation: 0}
	id2 := bson.NewObjectId()
	rw2 := &ResultWrapper{ID: id2, IndividualResult: ir2}
	i.Database.C("patient_cache").Insert(rw2)
	results, err := FindAllResultsForPatient(i.Database, "1234")
	util.CheckErr(err)
	c.Assert(len(results), Equals, 2)
}
