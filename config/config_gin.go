package config

import (
	"Qbus-manager/router"
	"Qbus-manager/router/middleware"
	"github.com/gin-gonic/gin"
)

func initGin() *gin.Engine {

	// Set gin mode.
	gin.SetMode(DataYaml.RunMode)

	// Create the Gin engine.
	g := gin.New()

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middleware's.
		middleware.Logging(),
	)
	return g
}