package cluster

import "qbus-manager/pkg/zookeeper"

type AddClusterVo struct {
	Name          string `json:"name"`
	ZookeeperList string `json:"zookeeper_list"`
}

type Detail struct {
	ClusterName   string               `json:"name"`
	ZookeeperList string               `json:"zookeeper_list"`
	BrokerList    []zookeeper.HostPort `json:"broker_list"`
	TopicList     []string             `json:"topic_list"`
}
