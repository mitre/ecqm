package models

import (
	"errors"
	"fmt"

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
	Exception                int32  `bson:"DENEXCEP,omitempty" json:"exception"`
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

type Population string

const (
	InitialPatientPopulation Population = "IPP"
	Denominator                         = "DENOM"
	Numerator                           = "NUMER"
	Exception                           = "DENEXCEP"
	Exclusion                           = "DENEX"
	Outlier                             = "antinumerator"
)

type PopulationQuery struct {
	MeasureID     string     `json:"measureId"`
	SubID         string     `json:"subId"`
	EffectiveDate int32      `json:"effectiveDate"`
	Limit         int        `json:"limit"`
	Offset        int        `json:"offset"`
	Population    Population `json:"population"`
}

func (p PopulationQuery) ToQuery() bson.M {
	query := bson.M{"value.measure_id": p.MeasureID, "value.effective_date": p.EffectiveDate}
	if p.SubID != "" {
		query["value.sub_id"] = p.SubID
	}
	populationKey := fmt.Sprintf("value.%s", p.Population)
	query[populationKey] = bson.M{"$gte": 1}
	return query
}

func ParsePopulationQuery(queryParam string) (Population, error) {
	switch queryParam {
	case "initialPatientPopulation":
		return InitialPatientPopulation, nil
	case "denominator":
		return Denominator, nil
	case "numerator":
		return Numerator, nil
	case "exception":
		return Exception, nil
	case "exclusion":
		return Exclusion, nil
	case "outlier":
		return Outlier, nil
	}
	return "", errors.New("Invalid popualtion type")
}

type PopulationResult struct {
	Total           int                `json:"total"`
	PopulationQuery PopulationQuery    `json:"populationQuery"`
	Patients        []IndividualResult `json:"patients"`
}

func FindResultsForMeasurePopulation(db *mgo.Database, query PopulationQuery) (*PopulationResult, error) {
	pr := &PopulationResult{PopulationQuery: query}
	q := db.C("patient_cache").Find(query.ToQuery())
	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	pr.Total = count
	q.Sort("value.last") // Sort by last name to make sure that offset is consistent between queries
	q.Skip(query.Offset)
	if query.Limit != 0 {
		q.Limit(query.Limit)
	} else {
		q.Limit(50)
	}
	var wrappedResults []ResultWrapper
	err = q.All(&wrappedResults)
	if err != nil {
		return nil, err
	}
	var result []IndividualResult
	for _, wr := range wrappedResults {
		result = append(result, wr.IndividualResult)
	}
	pr.Patients = result

	return pr, nil
}
