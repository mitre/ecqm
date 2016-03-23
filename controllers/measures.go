package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitre/ecqm/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func ShowMeasureHandler(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		measuresCollection := db.C("measures")
		measure := &models.Measure{}
		err := measuresCollection.Find(bson.M{"hqmf_id": id}).One(measure)
		if err == nil {
			c.JSON(http.StatusOK, measure)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}

func IndexMeasureHandler(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		measuresCollection := db.C("measures")
		var measures []models.Measure
		err := measuresCollection.Find(nil).All(&measures)
		if err == nil {
			c.JSON(http.StatusOK, measures)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}
