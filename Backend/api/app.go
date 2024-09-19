package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


var HelperApp *gin.Engine


func Listen(w http.ResponseWriter, r *http.Request) {
	HelperApp.ServeHTTP(w, r)
}
