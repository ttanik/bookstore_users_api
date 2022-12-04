package app

import (
	"github.com/gin-gonic/gin"
)

var (
	// Default is an Engine that has middlewares for log and recovery
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	router.Run(":8080")
}
