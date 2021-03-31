package _init

import (
	"github.com/gin-gonic/gin"
	"qbus-manager/router"
	"qbus-manager/router/middleware"
)

func ginInit() *gin.Engine {

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
