package zookeeper

import (
	"Qbus-manager/pkg/errno"
	"github.com/lexkong/log"
	"sync"
)

const (
	zkRoot               = "/kafka-manager"
	configsZkPath        = "/configs"
	deleteClustersZkPath = "/deleteClusters"
)

var (
	clusterConfigMap = make(map[string]*clusterConfig)
	rw               = sync.RWMutex{}
)

func AddCluster(clusterName string, zkHosts string) error {
	cc, err := NewDefaultClusterConfig(clusterName, zkHosts)
	if err != nil {
		log.Errorf(err, errno.ErrAddCluster.Message)
		return err
	}

	if err := createPersistent(zkRoot+configsZkPath+"/"+clusterName, cc); err != nil {
		log.Errorf(errno.ErrAddCluster, "ClusterName: `%v`, zkHosts: `%v`", clusterName, zkHosts)
		return errno.ErrAddCluster
	}
	rw.Lock()
	defer rw.Unlock()
	clusterConfigMap[clusterName] = cc
	return nil
}

func DisableCluster(clusterName string) error {
	rw.RLock()
	defer rw.RUnlock()
	cc, ok := clusterConfigMap[clusterName]
	if !ok {
		log.Errorf(errno.ErrClusterNotExist, "Cannot disable non-existing cluster : `%v`", clusterName)
		return errno.ErrClusterNotExist
	}

	cc.Enabled = false

	if err := set(zkRoot+configsZkPath+"/"+clusterName, cc); err != nil {
		log.Errorf(errno.ErrDisableCluster, "ClusterName: `%v`", clusterName)
		return errno.ErrDisableCluster
	}
	return nil
}

func DeleteCluster(clusterName string) error {
	rw.RLock()
	defer rw.RUnlock()
	cc, ok := clusterConfigMap[clusterName]
	if !ok {
		log.Errorf(errno.ErrClusterNotExist, "Cannot disable non-existing cluster : `%v`", clusterName)
		return errno.ErrClusterNotExist
	}

	if err := createPersistent(zkRoot+deleteClustersZkPath+"/"+clusterName, cc); err != nil {
		log.Errorf(errno.ErrDeleteCluster, "ClusterName: `%v`", clusterName)
		return errno.ErrDeleteCluster
	}
	return nil
}
