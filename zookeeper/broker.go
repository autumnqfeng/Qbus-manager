package zookeeper

import (
	"Qbus-manager/pkg/errno"
	"fmt"
	"github.com/lexkong/log"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
)

type B struct {
	Version   int      `json:"-"`
	JmxPort   int      `json:"jmx_port"`
	Timestamp string   `json:"timestamp"`
	Endpoints []string `json:"endpoints"`
	Host      string   `json:"host"`
	Port      int      `json:"port"`
	Id        int      `json:"id"`
}

func Broker(conn *zk.Conn, id int) (B, error) {
	var b B
	err := get(conn, fmt.Sprintf("/brokers/ids/%d", id), &b)
	b.Id = id
	return b, err
}

type HostPort struct {
	Id   int
	Host string
	Port int
}

var brokersListeners = make([]chan []HostPort, 0, 10)

func WatchBrokers(ch chan []HostPort) {
	brokersListeners = append(brokersListeners, ch)
}

func watchBrokers() {
	conn, err := GetDefaultConn()
	if err != nil {
		log.Errorf(errno.ErrGetConn, "err_code: `%v`, watch brokers `%v`", errno.ErrGetConn.Code, errno.ErrGetConn.Message)
		return
	}
	data, _, events, _ := watchChildren(conn, "/brokers/ids")

	for _, ch := range brokersListeners {
		ch <- DealBrokersChildren(conn, data)
	}

	_, ok := <-events
	if ok {
		watchBrokers()
	}
}

func DealBrokersChildren(conn *zk.Conn, data []string) []HostPort {
	hostPorts := make([]HostPort, len(data))
	for i, id := range data {
		intId, err := strconv.Atoi(id)
		if err != nil {
			continue
		}

		broker, err := Broker(conn, intId)
		if err != nil {
			continue
		}
		hostPorts[i] = HostPort{Id: intId, Host: broker.Host, Port: broker.Port}
	}
	return hostPorts
}

func GetBrokerListByCluster(cc *ClusterConfig) ([]HostPort, error) {
	clusterName := cc.ClusterName

	conn, err := GetConn(clusterName)
	if err != nil {
		log.Errorf(errno.ErrClusterConnect, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrClusterConnect.Code, errno.ErrClusterConnect.Message, clusterName)
		return nil, errno.ErrClusterConnect
	}

	data, err := children(conn, "/brokers/ids")
	if err != nil {
		log.Errorf(errno.ErrGetBroker, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrGetBroker.Code, errno.ErrGetBroker.Message, clusterName)
		return nil, errno.ErrGetBroker
	}

	return DealBrokersChildren(conn, data), nil
}
