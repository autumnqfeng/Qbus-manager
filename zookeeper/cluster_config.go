package zookeeper

import (
	"errors"
	"regexp"
	"strings"
)

type curatorConfig struct {
	ZkConnect       string `json:"zkConnect"`
	ZkMaxRetry      int    `json:"zkMaxRetry"`
	BaseSleepTimeMs int    `json:"baseSleepTimeMs"`
	MaxSleepTimeMs  int    `json:"maxSleepTimeMs"`
}

type clusterConfig struct {
	ClusterName              string        `json:"name"`
	CuratorConfig            curatorConfig `json:"curatorConfig"`
	Enabled                  bool          `json:"enabled"`
	Version                  string        `json:"kafkaVersion"`
	JmxEnabled               bool          `json:"jmxEnabled"`
	JmxUser                  interface{}   `json:"jmxUser"`
	JmxPass                  interface{}   `json:"jmxPass"`
	JmxSsl                   bool          `json:"jmxSsl"`
	PollConsumers            bool          `json:"pollConsumers"`
	FilterConsumers          bool          `json:"filterConsumers"`
	LogKafkaEnabled          bool          `json:"logkafkaEnabled"`
	ActiveOffsetCacheEnabled bool          `json:"activeOffsetCacheEnabled"`
	DisplaySizeEnabled       bool          `json:"displaySizeEnabled"`
	Tuning                   interface{}   `json:"tuning"`
}

func NewDefaultClusterConfig(clusterName string, zkHosts string) (*clusterConfig, error) {
	if err := validateClusterName(clusterName); err != nil {
		return nil, err
	}
	if err := validateZkHosts(zkHosts); err != nil {
		return nil, err
	}

	return &clusterConfig{
		ClusterName: clusterName,
		CuratorConfig: curatorConfig{
			ZkConnect:       zkHosts,
			ZkMaxRetry:      100,
			BaseSleepTimeMs: 100,
			MaxSleepTimeMs:  1000,
		},
		Enabled:                  true,
		Version:                  "0.9.0.1",
		JmxEnabled:               false,
		JmxUser:                  nil,
		JmxPass:                  nil,
		JmxSsl:                   false,
		PollConsumers:            true,
		FilterConsumers:          true,
		LogKafkaEnabled:          false,
		ActiveOffsetCacheEnabled: true,
		DisplaySizeEnabled:       false,
		Tuning:                   nil,
	}, nil
}

func validateClusterName(clusterName string) error {
	if len(clusterName) <= 0 {
		return errors.New("cluster name is illegal, can't be empty")
	}
	if strings.EqualFold(clusterName, ".") || strings.EqualFold(clusterName, "..") {
		return errors.New("cluster name cannot be \".\" or \"..\"")
	}
	if len(clusterName) > 255 {
		return errors.New("cluster name is illegal, can't be longer than 255 characters")
	}

	if matched, err := regexp.MatchString("[a-zA-Z0-9._\\-]", clusterName); err != nil && !matched {
		return errors.New("cluster name " + clusterName + " is illegal, contains a character other than ASCII alphanumerics, '.', '_' and '-'")
	}
	return nil
}

func validateZkHosts(zkHosts string) error {
	if len(zkHosts) <= 0 {
		return errors.New("cluster zk hosts is illegal, can't be empty")
	}
	return nil
}
