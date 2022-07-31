package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiMux() http.Handler {
	r := gin.Default()
	rg := r.Group("v1")
	rg.GET("/user", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello %s", "Joshua")
	})
	return r
}
