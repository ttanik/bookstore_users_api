package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ttanik/bookstore_users-api/logger"
)

var (
	// Default is an Engine that has middlewares for log and recovery
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	logger.Info("about to start the application...")
	router.Run(":8081")
}
