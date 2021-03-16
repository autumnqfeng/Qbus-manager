package zookeeper

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	conn                 *zk.Conn
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

func Exists(path string) bool {
	exists, _, _ := conn.Exists(path)
	return exists
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

func set(path string, data interface{}) error {
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

func createPersistent(path string, data interface{}) error {
	return create(path, data, 0)
}

func create(path string, data interface{}, flag int) error {
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
