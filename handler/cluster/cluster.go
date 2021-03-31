package cluster

import (
	"qbus-manager/handler"
	"qbus-manager/pkg/errno"
	"qbus-manager/pkg/zookeeper"

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
		handler.SendResponse(c, errno.ErrValidation, nil)
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
	result := make(map[string][]string)
	result["list"] = clusters
	handler.SendResponse(c, errno.OK, result)
}

func GetClusterDetail(c *gin.Context) {
	clusterName := c.Query("clustername")
	if clusterName == "" || len(clusterName) <= 0 {
		handler.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	cc, err := zookeeper.GetClusterConfig(clusterName)
	if err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	brokers, err := zookeeper.GetBrokerListByCluster(cc)
	if err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	topics, err := zookeeper.GetTopicByCluster(cc)
	if err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	result := Detail{
		ClusterName:   cc.ClusterName,
		ZookeeperList: cc.CuratorConfig.ZkConnect,
		BrokerList:    brokers,
		TopicList:     topics,
	}
	handler.SendResponse(c, errno.OK, result)
}

func GetClusterDiskInfo(c *gin.Context) {
	clusterName := c.Query("clustername")
	if clusterName == "" || len(clusterName) <= 0 {
		handler.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	cc, err := zookeeper.GetClusterConfig(clusterName)
	if err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	hosts, err := zookeeper.GetAllHost(cc)
	if err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	result, err := zookeeper.GetHostInfosByHosts(cc, hosts)
	if err != nil {
		handler.SendResponse(c, err, nil)
		return
	}

	handler.SendResponse(c, errno.OK, result)
}
