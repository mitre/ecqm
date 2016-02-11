package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/mitre/ecqm/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

func ShowQualityReportHandler(db *mgo.Database) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var id bson.ObjectId

		idString := c.Param("id")
		if bson.IsObjectIdHex(idString) {
			id = bson.ObjectIdHex(idString)
		} else {
			return errors.New("Invalid id")
		}

		queryCache := db.C("query_cache")
		qualityReport := &models.QualityReport{}
		err := queryCache.FindId(id).One(qualityReport)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, qualityReport)
	}
}

func CreateQualityReportHandler(db *mgo.Database) echo.HandlerFunc {
	return func(c *echo.Context) error {
		qualityReport := &models.QualityReport{}
		qualityReport.MeasureID = c.Form("measureId")
		qualityReport.SubID = c.Form("subId")
		ed := c.Form("effectiveDate")
		edInt, err := strconv.ParseInt(ed, 10, 32)
		if err != nil {
			return errors.New("Could not convert the effective date into an int32")
		}
		qualityReport.EffectiveDate = int32(edInt)
		err = validator.Validate(qualityReport)
		if err != nil {
			return err
		}
		qualityReport.ID = bson.NewObjectId()
		err = db.C("query_cache").Insert(qualityReport)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, qualityReport)
	}
}
