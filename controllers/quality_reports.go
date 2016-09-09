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

func findQualityReport(db *mgo.Database, c *gin.Context) *models.QualityReport {
	var id bson.ObjectId

	idString := c.Param("id")
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		c.String(http.StatusBadRequest, "Invalid ID")
		return nil
	}

	queryCache := db.C("query_cache")
	qualityReport := &models.QualityReport{}
	err := queryCache.FindId(id).One(qualityReport)
	if err != nil {
		if err == mgo.ErrNotFound {
			c.String(http.StatusNotFound, "Not found")
			return nil
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}
	return qualityReport
}

func ShowQualityReportHandler(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		qualityReport := findQualityReport(db, c)
		if qualityReport != nil {
			c.JSON(http.StatusOK, qualityReport)
		}
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

func ShowQualityReportPatientsHandler(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		populationParam, popPresent := c.GetQuery("population")
		if !popPresent {
			c.String(http.StatusBadRequest, "A population must be specified as a query parameter")
			return
		}
		population, err := models.ParsePopulationQuery(populationParam)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid population specified")
			return
		}
		qualityReport := findQualityReport(db, c)
		if qualityReport != nil {
			pq := models.PopulationQuery{MeasureID: qualityReport.MeasureID, SubID: qualityReport.SubID,
				EffectiveDate: qualityReport.EffectiveDate, Population: population}
			limitParam, limitPresent := c.GetQuery("limit")
			if limitPresent {
				limit, err := strconv.ParseInt(limitParam, 10, 32)
				if err != nil {
					c.String(http.StatusBadRequest, "Invalid limit specified")
					return
				}
				pq.Limit = int(limit)
			}
			offsetParam, offsetPresent := c.GetQuery("offset")
			if offsetPresent {
				offset, err := strconv.ParseInt(offsetParam, 10, 32)
				if err != nil {
					c.String(http.StatusBadRequest, "Invalid offset specified")
					return
				}
				pq.Offset = int(offset)
			}
			pr, err := models.FindResultsForMeasurePopulation(db, pq)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.JSON(http.StatusOK, pr)
		}
	}
}
