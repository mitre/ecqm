package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// QualityReport is a representation of a calculation of an electronic
// clinical quality measure
type QualityReport struct {
	ID              bson.ObjectId `bson:"_id" json:"id"`
	NPI             string        `bson:"npi,omitempty" json:"npi,omitempty"`
	CalculationTime time.Time     `bson:"calculation_time,omitempty" json:"calculationTime,omitempty"`
	Status          Status        `bson:"status,omitempty" json:"status,omitempty"`
	MeasureID       string        `bson:"measure_id,omitempty" json:"measureId,omitempty" validate:"nonzero"`
	SubID           string        `bson:"sub_id,omitempty" json:"subId,omitempty"`
	EffectiveDate   int32         `bson:"effective_date,omitempty" json:"effectiveDate,omitempty" validate:"nonzero"`
}

type Status struct {
	State string   `bson:"state,omitempty" json:"state,omitempty"`
	Log   []string `bson:"log,omitempty" json:"log,omitempty"`
}

type QualityReportResult struct {
	PopulationIDs            PopulationIDs `bson:"population_ids,omitempty" json:"populationIds,omitempty"`
	InitialPatientPopulation int32         `bson:"IPP" json:"initialPatientPopulation"`
	Denominator              int32         `bson:"DENOM,omitempty" json:"denominator,omitempty"`
	Exception                int32         `bson:"DENEXCP,omitempty" json:"exception,omitempty"`
	Exclusion                int32         `bson:"DENEX,omitempty" json:"exclusion,omitempty"`
	Numerator                int32         `bson:"NUMER,omitempty" json:"numerator,omitempty"`
	AntiNumerator            int32         `bson:"antinumerator,omitempty" json:"antinumerator,omitempty"`
	MeasurePopulation        int32         `bson:"MSRPOPL,omitempty" json:"measurePopulation,omitempty"`
	Observation              float32       `bson:"OBSERV,omitempty" json:"Observation,omitempty"`
}

type PopulationIDs struct {
	InitialPatientPopulation string `bson:"IPP,omitempty" json:"initialPatientPopulation,omitempty"`
	Denominator              string `bson:"DENOM,omitempty" json:"denominator,omitempty"`
	Exception                string `bson:"DENEXCP,omitempty" json:"exception,omitempty"`
	Exclusion                string `bson:"DENEX,omitempty" json:"exclusion,omitempty"`
	Numerator                string `bson:"NUMER,omitempty" json:"numerator,omitempty"`
	MeasurePopulation        string `bson:"MSRPOPL,omitempty" json:"measurePopulation,omitempty"`
	Observation              string `bson:"OBSERV,omitempty" json:"observation,omitempty"`
}
