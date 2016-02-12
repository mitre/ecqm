package models

import "gopkg.in/mgo.v2/bson"

type Measure struct {
	ID                 bson.ObjectId `bson:"_id" json:"id"`
	HQMFID             string        `bson:"hqmf_id,omitempty" json:"hqmfId,omitempty"`
	HQMFSetID          string        `bson:"hqmf_set_id,omitempty" json:"hqmfSetId,omitempty"`
	HQMFVersionNumber  int           `bson:"hqmf_version_number,omitempty" json:"hqmfVersionNumber,omitempty"`
	CMSID              string        `bson:"cms_id,omitempty" json:"cmsId,omitempty"`
	Name               string        `bson:"name,omitempty" json:"name,omitempty"`
	Description        string        `bson:"description,omitempty" json:"description,omitempty"`
	Type               string        `bson:"type,omitempty" json:"type,omitempty"`
	Category           string        `bson:"category,omitempty" json:"category,omitempty"`
	ContinuousVariable bool          `bson:"continuous_variable" json:"continuousVariable"`
	EpisodeOfCare      bool          `bson:"episode_of_care" json:"episodeOfCare"`
	PopulationIDs      PopulationIDs `bson:"population_ids,omitempty" json:"populationIds,omitempty"`
}
