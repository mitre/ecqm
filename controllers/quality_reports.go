package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mitre/ecqm/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

func ShowQualityReportHandler(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id bson.ObjectId

		idString := c.Param("id")
		if bson.IsObjectIdHex(idString) {
			id = bson.ObjectIdHex(idString)
		} else {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		queryCache := db.C("query_cache")
		qualityReport := &models.QualityReport{}
		err := queryCache.FindId(id).One(qualityReport)
		if err != nil {
			if err == mgo.ErrNotFound {
				c.String(http.StatusNotFound, "Not found")
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, qualityReport)
	}
}

func CreateQualityReportHandler(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		qualityReport := &models.QualityReport{}
		qualityReport.MeasureID = c.PostForm("measureId")
		qualityReport.SubID = c.PostForm("subId")
		ed := c.PostForm("effectiveDate")
		edInt, err := strconv.ParseInt(ed, 10, 32)
		if err != nil {
			c.String(http.StatusBadRequest, "Could not convert the effective date into an int32")
			return
		}
		qualityReport.EffectiveDate = int32(edInt)
		err = validator.Validate(qualityReport)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		err = models.FindOrCreateQualityReport(db, qualityReport)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, qualityReport)
	}
}
