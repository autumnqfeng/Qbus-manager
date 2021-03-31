package zookeeper

import (
	"encoding/json"
	"errors"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	ZkRoot               = "/kafka-manager"
	ConfigsZkPath        = "/configs"
	BaseClusterZkPath    = "/clusters"
	DeleteClustersZkPath = "/deleteClusters"
)

var pathDoesNotExistsErr = errors.New("zookeeper: Path does not exists")

func exists(conn *zk.Conn, path string) bool {
	exists, _, _ := conn.Exists(path)
	return exists
}

func watchChildren(conn *zk.Conn, path string) ([]string, *zk.Stat, <-chan zk.Event, error) {
	return conn.ChildrenW(path)
}

func children(conn *zk.Conn, path string) ([]string, error) {
	data, _, err := conn.Children(path)
	return data, err
}

func get(conn *zk.Conn, path string, v interface{}) error {
	if exists := exists(conn, path); !exists {
		return pathDoesNotExistsErr
	}
	data, _, err := conn.Get(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func all(conn *zk.Conn, path string, fn PermissionFunc) ([]string, error) {
	rows := make([]string, 0)
	if exists, _, _ := conn.Exists(path); !exists {
		return rows, pathDoesNotExistsErr
	}
	children, _, err := conn.Children(path)
	if err != nil {
		return rows, err
	}
	for _, c := range children {
		if fn(c) {
			rows = append(rows, c)
		}
	}
	return rows, nil
}

func set(conn *zk.Conn, path string, data interface{}) error {
	_, stat, err := conn.Exists(path)
	if err != nil {
		return err
	}
	enc, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Set(path, enc, stat.Version)
	return err
}

func createPersistent(conn *zk.Conn, path string, data interface{}) error {
	return create(conn, path, data, 0)
}

func create(conn *zk.Conn, path string, data interface{}, flag int) error {
	var (
		bytes []byte
		err   error
	)
	if str, ok := data.(string); ok {
		bytes = []byte(str)
	} else {
		bytes, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}
	_, err = conn.Create(path, bytes, int32(flag), zk.WorldACL(zk.PermAll))
	return err
}
