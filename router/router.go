package router

import (
	"Qbus-manager/handler/check"
	"Qbus-manager/handler/cluster"
	"Qbus-manager/handler/topic"
	"net/http"

	"Qbus-manager/router/middleware"

	// "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// Load loads the middleware's, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middleware's
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// pprof router
	//pprof.Register(g)

	clusterGroup := g.Group("/qbus/clusters")
	{
		clusterGroup.POST("/addcluster", cluster.AddCluster)
		clusterGroup.GET("/deletecluster", cluster.DeleteCluster)
		clusterGroup.GET("/listall", cluster.ListAllCluster)
		clusterGroup.GET("/getclusterdetail", cluster.GetClusterDetail)
	}
	//clusterGroupV1 := g.Group("/v1/qbus/clusters")
	//{
	//	clusterGroupV1.GET("")
	//}

	topicGroup := g.Group("/qbus/topics")
	{
		topicGroup.POST("/addtopic", topic.Create)
		topicGroup.GET("/deletetopic", topic.Delete)
	}

	// The health check handlers
	systemCheck := g.Group("/check")
	{
		systemCheck.GET("/health", check.HealthCheck)
		systemCheck.GET("/disk", check.DiskCheck)
		systemCheck.GET("/cpu", check.CPUCheck)
		systemCheck.GET("/ram", check.RAMCheck)
	}

	return g
}
