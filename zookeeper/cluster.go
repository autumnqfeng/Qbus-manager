package zookeeper

import (
	"Qbus-manager/pkg/errno"
	"github.com/lexkong/log"
)

const (
	zkRoot        = "/kafka-manager"
	configsZkPath = "/configs"
)

func AddCluster(clusterName string, zkHosts string) error {
	cc := NewDefaultClusterConfig(clusterName, zkHosts)
	if err := createPersistent(zkRoot+configsZkPath+"/"+clusterName, cc); err != nil {
		log.Errorf(errno.ErrAddCluster, "ClusterName: `%v`, zkHosts: `%v`", clusterName, zkHosts)
		return errno.ErrAddCluster
	}
	return nil
}
