package kafka

import (
	"fmt"

	"qbus-manager/pkg/zookeeper"

	"go.uber.org/zap"
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
			zap.L().Info("broker_change", zap.Int("brokerId", hostPort.Id), zap.Any("broker", hostPort))
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
