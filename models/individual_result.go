package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type IndividualResult struct {
	PatientID                string `bson:"patient_id,omitempty" json:"patientId,omitempty"`
	MedicalRecordID          string `bson:"medical_record_id,omitempty" json:"medicalRecordId,omitempty"`
	First                    string `bson:"first,omitempty" json:"first,omitempty"`
	Last                     string `bson:"last,omitempty" json:"last,omitempty"`
	Gender                   string `bson:"gender,omitempty" json:"gender,omitempty"`
	InitialPatientPopulation int32  `bson:"IPP" json:"initialPatientPopulation"`
	Denominator              int32  `bson:"DENOM,omitempty" json:"denominator"`
	Exception                int32  `bson:"DENEXCP,omitempty" json:"exception"`
	Exclusion                int32  `bson:"DENEX,omitempty" json:"exclusion"`
	Numerator                int32  `bson:"NUMER,omitempty" json:"numerator"`
	MeasureID                string `bson:"measure_id,omitempty" json:"measureId"`
	SubID                    string `bson:"sub_id,omitempty" json:"subId"`
	EffectiveDate            int32  `bson:"effective_date,omitempty" json:"effectiveDate,omitempty" validate:"nonzero"`
}

type ResultWrapper struct {
	ID               bson.ObjectId    `bson:"_id" json:"id"`
	IndividualResult IndividualResult `bson:"value"`
}

func FindAllResultsForPatient(db *mgo.Database, patientID string) ([]IndividualResult, error) {
	q := db.C("patient_cache").Find(bson.M{"value.patient_id": patientID})
	var wrappedResults []ResultWrapper
	err := q.All(&wrappedResults)
	if err != nil {
		return nil, err
	}
	var result []IndividualResult
	for _, wr := range wrappedResults {
		result = append(result, wr.IndividualResult)
	}

	return result, nil
}
