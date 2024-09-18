package api

import (
	"CengkeHelper/logger"
	"CengkeHelper/process"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)


var HelperApp *gin.Engine

func Listen(w http.ResponseWriter, r *http.Request) {
	HelperApp.ServeHTTP(w, r)
}
