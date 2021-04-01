package zookeeper

import (
	"fmt"
	"strings"
	"sync"

	"qbus-manager/pkg/errno"

	"go.uber.org/zap"
)

var (
	clusterConfigMap = make(map[string]*ClusterConfig)
	rw               = sync.RWMutex{}
)

func AddCluster(clusterName string, zkHosts string) error {
	cc, err := NewDefaultClusterConfig(clusterName, zkHosts)
	if err != nil {
		zap.L().Error(fmt.Sprintf(errno.ErrAddCluster.Message), zap.Error(err))
		return err
	}

	conn, err := GetDefaultConn()
	if err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`,err_msg `%v`, clusterName: `%v`, zkHosts: `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message, clusterName, zkHosts), zap.Error(errno.ErrGetConn))
		return err
	}

	if err := createPersistent(conn, ZkRoot+ConfigsZkPath+"/"+clusterName, cc); err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`,err_msg `%v`, clusterName: `%v`, zkHosts: `%v`", errno.ErrAddCluster.Code, errno.ErrAddCluster.Message, clusterName, zkHosts), zap.Error(errno.ErrAddCluster))
		return errno.ErrAddCluster
	}
	if err = AddConn(clusterName, strings.Split(zkHosts, ",")); err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`,err_msg `%v`, clusterName: `%v`, zkHosts: `%v`", errno.ErrAddConnPoll.Code, errno.ErrAddConnPoll.Message, clusterName, zkHosts), zap.Error(errno.ErrAddConnPoll))
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
		zap.L().Error(fmt.Sprintf("Cannot disable non-existing cluster : `%v`", clusterName), zap.Error(errno.ErrClusterNotExist))
		return errno.ErrClusterNotExist
	}

	cc.Enabled = false

	conn, err := GetDefaultConn()
	if err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`,err_msg `%v`, clusterName: `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message, clusterName), zap.Error(errno.ErrGetConn))
		return err
	}

	if err := set(conn, ZkRoot+ConfigsZkPath+"/"+clusterName, cc); err != nil {
		zap.L().Error(fmt.Sprintf("ClusterName: `%v`", clusterName), zap.Error(errno.ErrDisableCluster))
		return errno.ErrDisableCluster
	}
	return nil
}

func DeleteCluster(clusterName string) error {
	rw.RLock()
	defer rw.RUnlock()
	cc, ok := clusterConfigMap[clusterName]
	if !ok {
		zap.L().Error(fmt.Sprintf("Cannot disable non-existing cluster : `%v`", clusterName), zap.Error(errno.ErrClusterNotExist))
		return errno.ErrClusterNotExist
	}

	conn, err := GetDefaultConn()
	if err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`,err_msg `%v`, clusterName: `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message, clusterName), zap.Error(errno.ErrGetConn))
		return err
	}

	if err := createPersistent(conn, ZkRoot+DeleteClustersZkPath+"/"+clusterName, cc); err != nil {
		zap.L().Error(fmt.Sprintf("ClusterName: `%v`", clusterName), zap.Error(errno.ErrDeleteCluster))
		return errno.ErrDeleteCluster
	}

	DelConn(clusterName)
	return nil
}

func ListAllCluster() ([]string, error) {
	conn, err := GetDefaultConn()
	if err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`,err_msg `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message), zap.Error(errno.ErrGetConn))
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
		zap.L().Error(fmt.Sprintf("err_code: `%v`,err_msg `%v`, clusterName: `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message, clusterName), zap.Error(errno.ErrGetConn))
		return nil, errno.ErrGetConn
	}

	path := ZkRoot + ConfigsZkPath + "/" + clusterName
	if !exists(conn, path) {
		zap.L().Error(fmt.Sprintf("`%v`, cluster name :`%v`", errno.ErrNotFoundClusterConfig.Message, clusterName), zap.Error(errno.ErrNotFoundClusterConfig))
		return nil, errno.ErrNotFoundClusterConfig
	}
	var cc ClusterConfig
	if err := get(conn, path, &cc); err != nil {
		zap.L().Error(fmt.Sprintf("cluster name :`%v`", clusterName), zap.Error(err))
		return nil, err
	}
	return &cc, nil
}

func GetAllHost(cc *ClusterConfig) ([]string, error) {
	clusterName := cc.ClusterName

	conn, err := GetConn(clusterName)
	if err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrClusterConnect.Code, errno.ErrClusterConnect.Message, clusterName), zap.Error(errno.ErrClusterConnect))
		return nil, errno.ErrClusterConnect
	}

	data, err := children(conn, "/qbus2/status")
	if err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrGetBroker.Code, errno.ErrGetBroker.Message, clusterName), zap.Error(errno.ErrGetBroker))
		return nil, errno.ErrGetBroker
	}

	return data, nil
}

func GetHostInfosByHosts(cc *ClusterConfig, hosts []string) ([]DiskInfo, error) {
	clusterName := cc.ClusterName

	conn, err := GetConn(clusterName)
	if err != nil {
		zap.L().Error(fmt.Sprintf("err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrClusterConnect.Code, errno.ErrClusterConnect.Message, clusterName), zap.Error(errno.ErrClusterConnect))
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
			zap.L().Error(fmt.Sprintf("err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrGetBroker.Code, errno.ErrGetBroker.Message, clusterName), zap.Error(errno.ErrGetBroker))
			return nil, errno.ErrGetBroker
		}
		diskInfos = append(diskInfos, *NewDiskInfo(host, disk))
	}

	return diskInfos, nil
}
