package zookeeper

import (
	"qbus-manager/pkg/errno"

	"github.com/lexkong/log"
)

func GetTopicByCluster(cc *ClusterConfig) ([]string, error) {
	clusterName := cc.ClusterName

	conn, err := GetConn(clusterName)
	if err != nil {
		log.Errorf(errno.ErrClusterConnect, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrClusterConnect.Code, errno.ErrClusterConnect.Message, clusterName)
		return nil, errno.ErrClusterConnect
	}

	data, err := children(conn, "/brokers/topics")
	if err != nil {
		log.Errorf(errno.ErrGetBroker, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrGetBroker.Code, errno.ErrGetBroker.Message, clusterName)
		return nil, errno.ErrGetBroker
	}
	return data, nil
}
