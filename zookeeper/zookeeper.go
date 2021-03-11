package zookeeper

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	conn  *zk.Conn
	PathDoesNotExistsErr = errors.New("zookeeper: Path does not exists")
)

func Stop() {
	if conn != nil {
		conn.Close()
	}
}

func Connect(urls []string) error {
	var err error
	Stop()
	conn, _, err = zk.Connect(urls, 30*time.Second)
	if err != nil {
		return err
	}
	go watchBrokers()
	return nil
}

func WatchChildren(path string) ([]string, *zk.Stat, <-chan zk.Event, error) {
	return conn.ChildrenW(path)
}

func get(path string, v interface{}) error {
	if exists, _, _ := conn.Exists(path); !exists {
		return PathDoesNotExistsErr
	}
	data, _, err := conn.Get(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}