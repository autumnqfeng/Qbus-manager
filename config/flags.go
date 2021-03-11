package config

import (
	"github.com/spf13/pflag"
	"strings"
)

var (
	cfg     = pflag.StringP("ConfigData", "c", "", "Qbus-manager ConfigData file path")
	version = pflag.BoolP("version", "v", false, "show version info.")
	zk      = pflag.StringP("zookeeper", "z", "123.57.45.66:2181", "The connection string for the zookeeper connection in the form host:port. Multiple hosts can be given to allow fail-over.")
)

var (
	Cfg          string
	SysVersion   bool
	ZookeeperURL []string
)

func Parse() {
	pflag.Parse()
	Cfg = *cfg
	SysVersion = *version
	ZookeeperURL = strings.Split(*zk, ",")
}
