package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thomaspepio/hn-queries/index"
)

func main() {
	index := ingestHnLogs()
	startEndpoints(index)
}

func ingestHnLogs() *index.Index {
	return nil
}

func startEndpoints(index *index.Index) {
	router := gin.Default()

	router.GET("/1/queries/count/:datePrefix", func(context *gin.Context) {
		name := context.Param("datePrefix")
		context.String(http.StatusOK, "Hello %s", name)
	})

	router.GET("/1/queries/popular/:datePrefix", func(context *gin.Context) {
		name := context.Param("datePrefix")
		size := context.Query("size")
		context.String(http.StatusOK, "Hello %s %s", name, size)
	})

	router.Run(":8080")
}
