package cluster

import (
	"Qbus-manager/handler"
	"Qbus-manager/pkg/errno"
	"Qbus-manager/zookeeper"
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
