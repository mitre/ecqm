package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mitre/ecqm/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func ShowMeasureHandler(db *mgo.Database) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")
		measuresCollection := db.C("measures")
		measure := &models.Measure{}
		err := measuresCollection.Find(bson.M{"hqmf_id": id}).One(measure)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, measure)
	}
}

func IndexMeasureHandler(db *mgo.Database) echo.HandlerFunc {
	return func(c *echo.Context) error {
		measuresCollection := db.C("measures")
		var measures []models.Measure
		err := measuresCollection.Find(nil).All(&measures)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, measures)
	}
}
