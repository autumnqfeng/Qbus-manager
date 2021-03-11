package zookeeper

import (
	"fmt"
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

func Broker(id int) (B, error) {
	var b B
	err := get(fmt.Sprintf("/brokers/ids/%d", id), &b)
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
	data, _, events, _ := WatchChildren("/brokers/ids")
	list := make([]HostPort, len(data))
	for i, id := range data {
		intId, err := strconv.Atoi(id)
		if err != nil {
			continue
		}

		broker, err := Broker(intId)
		if err != nil {
			continue
		}
		list[i] =HostPort{Id: intId, Host: broker.Host, Port: broker.Port}
	}

	for _, ch := range brokersListeners {
		ch <- list
	}

	_, ok := <-events
	if ok {
		watchBrokers()
	}
}
