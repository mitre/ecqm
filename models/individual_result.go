package models

import (
	"errors"

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
	SubID                    string `bson:"measure_id,omitempty" json:"measureId"`
	EffectiveDate            int32  `bson:"effective_date,omitempty" json:"effectiveDate,omitempty" validate:"nonzero"`
}

type ResultWrapper struct {
	ID               bson.ObjectId     `bson:"_id" json:"id"`
	IndividualResult *IndividualResult `bson:"value"`
}

// FindAndPopulateIndividualResult will attempt to find an IndividualResult in
// the patient_cache based on the measure id, sub id and effective date and patient id passed in.
// If it finds the associated document in the database, it will return true
// and populate the other fields in the IndividualResult that is passed in.
// Otherwise, it will return false, and the passed in IndividualResult will remain
// unchanged.
func FindAndPopulateIndividualResult(db *mgo.Database, ir *IndividualResult) (bool, error) {
	query := bson.M{"value.measure_id": ir.MeasureID, "value.effective_date": ir.EffectiveDate, "value.patient_id": ir.PatientID}
	if ir.SubID != "" {
		query["value.sub_id"] = ir.SubID
	}
	rw := &ResultWrapper{IndividualResult: ir}
	q := db.C("patient_cache").Find(query)
	count, err := q.Count()
	if err != nil {
		return false, err
	}
	switch count {
	case 0:
		return false, nil
	case 1:
		err = q.One(rw)
		if err != nil {
			return false, err
		}
	default:
		return false, errors.New("Found more than one patient cache entry for this")
	}
	return true, nil
}
