package zookeeper

type clusterConfig struct {
	ClusterName              string      `json:"ClusterName"`
	KafkaVersion             string      `json:"kafkaVersion"`
	ZkHosts                  string      `json:"zkHosts"`
	JmxEnabled               bool        `json:"jmxEnabled"`
	JmxUser                  interface{} `json:"jmxUser"`
	JmxPass                  interface{} `json:"jmxPass"`
	JmxSsl                   bool        `json:"jmxSsl"`
	PollConsumers            bool        `json:"pollConsumers"`
	FilterConsumers          bool        `json:"filterConsumers"`
	Tuning                   interface{} `json:"tuning"`
	LogKafkaEnabled          bool        `json:"logkafkaEnabled"`
	ActiveOffsetCacheEnabled bool        `json:"activeOffsetCacheEnabled"`
	DisplaySizeEnabled       bool        `json:"displaySizeEnabled"`
}

func NewDefaultClusterConfig(clusterName string, zkHosts string) *clusterConfig {
	return &clusterConfig{
		ClusterName:              clusterName,
		KafkaVersion:             "0.9.0.1",
		ZkHosts:                  zkHosts,
		JmxEnabled:               false,
		JmxUser:                  nil,
		JmxPass:                  nil,
		JmxSsl:                   false,
		PollConsumers:            true,
		FilterConsumers:          true,
		Tuning:                   nil,
		LogKafkaEnabled:          false,
		ActiveOffsetCacheEnabled: true,
		DisplaySizeEnabled:       false,
	}
}
