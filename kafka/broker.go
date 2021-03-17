package kafka

import (
	"Qbus-manager/pkg/errno"
	"Qbus-manager/zookeeper"
	"fmt"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	. "strconv"
)

func init() {
	go handleBrokerChanges()
}

func handleBrokerChanges() {
	brokerChanges := make(chan []zookeeper.HostPort)
	zookeeper.WatchBrokers(brokerChanges)
	defer close(brokerChanges)
	for hostPorts := range brokerChanges {
		hashMap := make(map[int]zookeeper.HostPort)
		for _, hostPort := range hostPorts {
			log.Info("broker_change", lager.Data{Itoa(hostPort.Id): hostPort})
			hashMap[hostPort.Id] = hostPort
		}
		Brokers = hashMap
	}
}

type BrokersMap map[int]zookeeper.HostPort

var Brokers = make(BrokersMap)

func (b BrokersMap) KafkaUrlList() []string {
	var res []string
	for k := range b {
		if v := b.KafkaUrl(k); v != "" {
			res = append(res, b.KafkaUrl(k))
		}
	}
	return res
}

func (b BrokersMap) KafkaUrl(k int) string {
	if b[k].Host == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", b[k].Host, b[k].Port)
}

func GetBrokerListByCluster(cc *zookeeper.ClusterConfig) ([]zookeeper.HostPort, error) {
	clusterName := cc.ClusterName

	conn, err := zookeeper.GetConn(clusterName)
	if err != nil {
		log.Errorf(errno.ErrClusterConnect, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrClusterConnect.Code, errno.ErrClusterConnect.Message, clusterName)
		return nil, errno.ErrClusterConnect
	}

	data, _, err := conn.Children("/brokers/ids")
	if err != nil {
		log.Errorf(errno.ErrGetBroker, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrGetBroker.Code, errno.ErrGetBroker.Message, clusterName)
		return nil, errno.ErrGetBroker
	}

	return zookeeper.DealBrokersChildren(conn, data), nil
}
