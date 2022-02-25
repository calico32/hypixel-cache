package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func ensureAuthenticated(c *gin.Context) {
	if c.GetHeader("X-Secret") != os.Getenv("API_SECRET") {
		finish(c, 401, NewFailure("unauthorized"))
		return
	}

	c.Next()
}
