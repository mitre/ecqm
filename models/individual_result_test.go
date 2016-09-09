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
	ir := IndividualResult{MeasureID: "abcd", EffectiveDate: 1234, PatientID: "1234", InitialPatientPopulation: 1, Last: "Jones"}
	id := bson.NewObjectId()
	rw := &ResultWrapper{ID: id, IndividualResult: ir}
	i.Database.C("patient_cache").Insert(rw)
	ir2 := IndividualResult{MeasureID: "efgh", EffectiveDate: 1234, PatientID: "1234", InitialPatientPopulation: 1, Last: "B"}
	id2 := bson.NewObjectId()
	rw2 := &ResultWrapper{ID: id2, IndividualResult: ir2}
	i.Database.C("patient_cache").Insert(rw2)
	ir3 := IndividualResult{MeasureID: "efgh", EffectiveDate: 1234, PatientID: "5678", InitialPatientPopulation: 1, Denominator: 1, Last: "A"}
	id3 := bson.NewObjectId()
	rw3 := &ResultWrapper{ID: id3, IndividualResult: ir3}
	i.Database.C("patient_cache").Insert(rw3)
}

func (i *IndividualResultSuite) TearDownTest(c *C) {
	i.Database.Session.Close()
	i.DBServer.Wipe()
}

func (i *IndividualResultSuite) TestFindAllResultsForPatient(c *C) {
	results, err := FindAllResultsForPatient(i.Database, "1234")
	util.CheckErr(err)
	c.Assert(len(results), Equals, 2)
}

func (i *IndividualResultSuite) TestFindResultsForMeasurePopulation(c *C) {
	pq := PopulationQuery{MeasureID: "efgh", EffectiveDate: 1234, Population: InitialPatientPopulation}
	pr, err := FindResultsForMeasurePopulation(i.Database, pq)
	util.CheckErr(err)
	c.Assert(pr.Total, Equals, 2)
	pt := pr.Patients[0]
	c.Assert(pt.Last, Equals, "A")
	pq.Population = Denominator
	pr, err = FindResultsForMeasurePopulation(i.Database, pq)
	util.CheckErr(err)
	c.Assert(pr.Total, Equals, 1)
	pt = pr.Patients[0]
	c.Assert(pt.Last, Equals, "A")
}
