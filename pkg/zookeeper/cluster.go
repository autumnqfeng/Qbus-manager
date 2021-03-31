package zookeeper

import (
	"fmt"
	"strings"
	"sync"

	"qbus-manager/pkg/errno"

	"github.com/lexkong/log"
)

var (
	clusterConfigMap = make(map[string]*ClusterConfig)
	rw               = sync.RWMutex{}
)

func AddCluster(clusterName string, zkHosts string) error {
	cc, err := NewDefaultClusterConfig(clusterName, zkHosts)
	if err != nil {
		log.Errorf(err, errno.ErrAddCluster.Message)
		return err
	}

	conn, err := GetDefaultConn()
	if err != nil {
		log.Errorf(errno.ErrGetConn, "err_code: `%v`,err_msg `%v`, clusterName: `%v`, zkHosts: `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message, clusterName, zkHosts)
		return err
	}

	if err := createPersistent(conn, ZkRoot+ConfigsZkPath+"/"+clusterName, cc); err != nil {
		log.Errorf(errno.ErrAddCluster, "err_code: `%v`,err_msg `%v`, clusterName: `%v`, zkHosts: `%v`", errno.ErrAddCluster.Code, errno.ErrAddCluster.Message, clusterName, zkHosts)
		return errno.ErrAddCluster
	}
	if err = AddConn(clusterName, strings.Split(zkHosts, ",")); err != nil {
		log.Errorf(errno.ErrAddConnPoll, "err_code: `%v`,err_msg `%v`, clusterName: `%v`, zkHosts: `%v`", errno.ErrAddConnPoll.Code, errno.ErrAddConnPoll.Message, clusterName, zkHosts)
		return errno.ErrAddConnPoll
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

	conn, err := GetDefaultConn()
	if err != nil {
		log.Errorf(errno.ErrGetConn, "err_code: `%v`,err_msg `%v`, clusterName: `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message, clusterName)
		return err
	}

	if err := set(conn, ZkRoot+ConfigsZkPath+"/"+clusterName, cc); err != nil {
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

	conn, err := GetDefaultConn()
	if err != nil {
		log.Errorf(errno.ErrGetConn, "err_code: `%v`,err_msg `%v`, clusterName: `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message, clusterName)
		return err
	}

	if err := createPersistent(conn, ZkRoot+DeleteClustersZkPath+"/"+clusterName, cc); err != nil {
		log.Errorf(errno.ErrDeleteCluster, "ClusterName: `%v`", clusterName)
		return errno.ErrDeleteCluster
	}

	DelConn(clusterName)
	return nil
}

func ListAllCluster() ([]string, error) {
	conn, err := GetDefaultConn()
	if err != nil {
		log.Errorf(errno.ErrGetConn, "err_code: `%v`,err_msg `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message)
		return nil, errno.ErrGetConn
	}

	clusters, err := all(conn, ZkRoot+BaseClusterZkPath, func(string) bool { return true })
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func GetClusterConfig(clusterName string) (*ClusterConfig, error) {
	conn, err := GetDefaultConn()
	if err != nil {
		log.Errorf(errno.ErrGetConn, "err_code: `%v`,err_msg `%v`, clusterName: `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message, clusterName)
		return nil, errno.ErrGetConn
	}

	path := ZkRoot + ConfigsZkPath + "/" + clusterName
	if !exists(conn, path) {
		log.Errorf(errno.ErrNotFoundClusterConfig, "`%v`, cluster name :`%v`", errno.ErrNotFoundClusterConfig.Message, clusterName)
		return nil, errno.ErrNotFoundClusterConfig
	}
	var cc ClusterConfig
	if err := get(conn, path, &cc); err != nil {
		log.Errorf(err, ", cluster name :`%v`", clusterName)
		return nil, err
	}
	return &cc, nil
}

func GetAllHost(cc *ClusterConfig) ([]string, error) {
	clusterName := cc.ClusterName

	conn, err := GetConn(clusterName)
	if err != nil {
		log.Errorf(errno.ErrClusterConnect, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrClusterConnect.Code, errno.ErrClusterConnect.Message, clusterName)
		return nil, errno.ErrClusterConnect
	}

	data, err := children(conn, "/qbus2/status")
	if err != nil {
		log.Errorf(errno.ErrGetBroker, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrGetBroker.Code, errno.ErrGetBroker.Message, clusterName)
		return nil, errno.ErrGetBroker
	}

	return data, nil
}

func GetHostInfosByHosts(cc *ClusterConfig, hosts []string) ([]DiskInfo, error) {
	clusterName := cc.ClusterName

	conn, err := GetConn(clusterName)
	if err != nil {
		log.Errorf(errno.ErrClusterConnect, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrClusterConnect.Code, errno.ErrClusterConnect.Message, clusterName)
		return nil, errno.ErrClusterConnect
	}

	diskInfos := make([]DiskInfo, 0)

	for _, host := range hosts {
		path := fmt.Sprintf("/qbus2/status/%v/disk", host)
		if !exists(conn, path) {
			continue
		}

		var disk int
		if err := get(conn, path, &disk); err != nil {
			log.Errorf(errno.ErrGetBroker, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrGetBroker.Code, errno.ErrGetBroker.Message, clusterName)
			return nil, errno.ErrGetBroker
		}
		diskInfos = append(diskInfos, *NewDiskInfo(host, disk))
	}

	return diskInfos, nil
}
