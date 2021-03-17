package cluster

import (
	"Qbus-manager/handler"
	"Qbus-manager/pkg/errno"
	"Qbus-manager/zookeeper"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func AddCluster(c *gin.Context) {
	log.Info("Add cluster function called.")
	var acv AddClusterVo
	if err := c.Bind(&acv); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	if err := zookeeper.AddCluster(acv.Name, acv.ZookeeperList); err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	handler.SendResponse(c, errno.OK, nil)
}

func DeleteCluster(c *gin.Context) {
	clusterName := c.Query("clustername")
	if clusterName == "" || len(clusterName) <= 0 {
		handler.SendResponse(c, errors.New("param is empty"), nil)
		return
	}
	if err := zookeeper.DisableCluster(clusterName); err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	if err := zookeeper.DeleteCluster(clusterName); err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	handler.SendResponse(c, errno.OK, nil)
}

func ListAllCluster(c *gin.Context) {
	clusters, err := zookeeper.ListAllCluster()
	if err != nil {
		handler.SendResponse(c, err, nil)
		return
	}
	handler.SendResponse(c, errno.OK, clusters)
}
