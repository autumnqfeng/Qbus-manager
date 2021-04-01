package zookeeper

import (
	"fmt"

	"qbus-manager/pkg/errno"

	"go.uber.org/zap"
)

func GetTopicByCluster(cc *ClusterConfig) ([]string, error) {
	clusterName := cc.ClusterName

	conn, err := GetConn(clusterName)
	if err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrClusterConnect.Code, errno.ErrClusterConnect.Message, clusterName),
			zap.Error(errno.ErrClusterConnect))
		return nil, errno.ErrClusterConnect
	}

	data, err := children(conn, "/brokers/topics")
	if err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrGetBroker.Code, errno.ErrGetBroker.Message, clusterName),
			zap.Error(errno.ErrGetBroker))
		return nil, errno.ErrGetBroker
	}
	return data, nil
}
