package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type successPlayerFound struct {
	Success   bool   `json:"success"`
	FetchedAt string `json:"fetchedAt"`
	Cached    bool   `json:"cached"`
	Username  string `json:"username"`
	Uuid      string `json:"uuid"`
	Player    Json   `json:"player"`
}

func NewSuccessPlayerFound(data CachedPlayer, cached bool) successPlayerFound {
	return successPlayerFound{
		Success:   true,
		FetchedAt: data.FetchedAt.Format(time.RFC3339),
		Cached:    cached,
		Username:  data.Player["displayname"].(string),
		Uuid:      data.Player["uuid"].(string),
		Player:    data.Player,
	}
}

type successNotFound struct {
	Success   bool   `json:"success"`
	FetchedAt string `json:"fetchedAt"`
	Cached    bool   `json:"cached"`
}

func NewSuccessNotFound(fetchedAt time.Time, cached bool) successNotFound {
	return successNotFound{
		Success:   true,
		FetchedAt: fetchedAt.Format(time.RFC3339),
		Cached:    cached,
	}
}

type failure struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func NewFailure(reason string) failure {
	return failure{
		Success: false,
		Error:   reason,
	}
}

func cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET")
	c.Header("Access-Control-Expose-Headers", "X-Response-Time")
}

func responseTimeStart(c *gin.Context) {
	c.Header("X-Response-Start", strconv.FormatInt(time.Now().UnixMicro(), 10))
}

func responseTimeEnd(c *gin.Context) {
	now := time.Now().UnixMicro()
	start, err := strconv.ParseInt(c.Writer.Header().Get("X-Response-Start"), 10, 64)
	if err != nil {
		panic("error parsing X-Response-Start:" + err.Error())
	}

	duration := now - start

	c.Header("X-Response-Start", "")
	c.Header("X-Response-Time", strconv.FormatInt(duration/1000, 10)+"ms")
}

func finish(c *gin.Context, code int, json interface{}) {
	responseTimeEnd(c)
	c.AbortWithStatusJSON(code, json)
}
