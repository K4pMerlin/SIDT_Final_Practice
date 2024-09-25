package api

import (
	"CengkeHelper/Backend/process"
	"net/http"

	"github.com/gin-gonic/gin"
)

var HelperApp *gin.Engine

func init() {
	HelperApp.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(400, gin.H{"error": "Query parameter 'q' is required"})
			return
		}

		results := process.SearchCourses(query)
		c.JSON(200, results)
	})

}

func Listen(w http.ResponseWriter, r *http.Request) {
	HelperApp.ServeHTTP(w, r)
}
