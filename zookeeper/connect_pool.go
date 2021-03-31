package zookeeper

import (
	"errors"
	"github.com/samuel/go-zookeeper/zk"
	"qbus-manager/pkg/errno"
	"strings"
	"time"
)

// zookeeper connect poll
var (
	zkConnects = make(map[string]*zk.Conn)
	timeOut    = 30 * time.Second
)

func Init(urls []string) error {
	err := AddConn("default", urls)
	go watchBrokers()
	return err
}

func GetDefaultConn() (*zk.Conn, error) {
	return GetConn("default")
}

func GetConn(clusterName string) (*zk.Conn, error) {
	if zkConnects[clusterName] == nil {
		if err := getConn(clusterName); err != nil {
			return nil, err
		}
	}
	if zkConnects[clusterName] == nil {
		return nil, errors.New(clusterName + ": get connection failed")
	}
	return zkConnects[clusterName], nil
}

func getConn(clusterName string) error {
	cc, err := GetClusterConfig(clusterName)
	if err != nil {
		return errno.ErrClusterMsgNotInKM
	}
	return AddConn(cc.ClusterName, strings.Split(cc.CuratorConfig.ZkConnect, ","))
}

func AddConn(clusterName string, urls []string) error {
	if zkConnects[clusterName] != nil {
		DelConn(clusterName)
	}

	conn, _, err := zk.Connect(urls, timeOut)
	if err != nil {
		return err
	}

	zkConnects[clusterName] = conn
	return nil
}

func DelConn(clusterName string) {
	conn := zkConnects[clusterName]
	if conn != nil {
		conn.Close()
		delete(zkConnects, clusterName)
	}
}

func StopAll() {
	for clusterName, conn := range zkConnects {
		conn.Close()
		delete(zkConnects, clusterName)
	}
}
