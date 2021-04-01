package _init

import (
	"github.com/gin-gonic/gin"
	"qbus-manager/configs"
	"qbus-manager/pkg/logger"
	"qbus-manager/router"
)

func ginInit() *gin.Engine {

	// Set gin mode.
	gin.SetMode(configs.Conf.Mode)

	// Create the Gin engine.
	g := gin.New()

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middleware's.
		logger.GinLogger(),
		logger.GinRecovery(true),
	)
	return g
}
