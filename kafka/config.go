package kafka

import (
	"Qbus-manager/zookeeper"

	"fmt"
	. "strconv"

	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
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
		BrokerUrls = hashMap
	}
}

type BrokerURLs map[int]zookeeper.HostPort

var BrokerUrls = make(BrokerURLs)

func (b BrokerURLs) List() []string {
	var res []string
	for k := range b {
		if v := b.KafkaUrl(k); v != "" {
			res = append(res, b.KafkaUrl(k))
		}
	}
	return res
}

func (b BrokerURLs) KafkaUrl(k int) string {
	if b[k].Host == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", b[k].Host, b[k].Port)
}
