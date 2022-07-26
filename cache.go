package main

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var playerCache *cache.Cache = cache.New(15*time.Minute, 30*time.Minute)
var uuidCache *cache.Cache = cache.New(60*time.Minute, 120*time.Minute)

var uuidRegexp *regexp.Regexp = regexp.MustCompile("^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$")
var usernameRegexp *regexp.Regexp = regexp.MustCompile("^[a-zA-Z0-9_]{3,16}$")

type CachedPlayer struct {
	Player    map[string]interface{}
	FetchedAt time.Time
}

func findCachedPlayer(c *gin.Context) {
	if c.Writer.Written() {
		c.Next()
		return
	}

	uuid, err := getUuid(c.Param("type"), c.Param("identifier"))
	if err != nil {
		finish(c, 500, NewFailure(err.Error()))
	}

	if cached, ok := playerCache.Get(uuid); ok {
		cachedPlayer := cached.(CachedPlayer)
		if cachedPlayer.Player != nil {
			finish(c, 200, NewSuccessPlayerFound(cachedPlayer, true))
		} else {
			finish(c, 200, NewSuccessNotFound(cachedPlayer.FetchedAt, false))
		}
	} else {
		apiResolvedUuids[c.Request] = uuid
		c.Next()
	}
}
