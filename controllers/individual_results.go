package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitre/ecqm/models"
	"gopkg.in/mgo.v2"
)

func ShowIndividualResultsForPatientHandler(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientIdString := c.Param("id")
		results, err := models.FindAllResultsForPatient(db, patientIdString)

		if err != nil {
			if err == mgo.ErrNotFound {
				c.String(http.StatusNotFound, "Not found")
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, results)
	}
}
