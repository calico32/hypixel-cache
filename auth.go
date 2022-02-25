package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func ensureAuthenticated(c *gin.Context) {
	if c.GetHeader("X-Secret") != os.Getenv("API_SECRET") {
		c.AbortWithStatusJSON(401, NewFailure("Unauthorized"))
		return
	}

	c.Next()
}
