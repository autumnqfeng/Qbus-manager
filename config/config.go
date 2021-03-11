package config

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type Config struct {
	RunMode      string `yaml:"run_mode"`
	Addr         string `yaml:"addr"`
	Name         string `yaml:"name"`
	Url          string `yaml:"url"`
	MaxPingCount int    `yaml:"max_ping_count"`
}

var DataYaml *Config

func parseConfig() *Config {
	cfg := Config{
		RunMode:      viper.GetString("run_mode"),
		Addr:         viper.GetString("addr"),
		Name:         viper.GetString("name"),
		Url:          viper.GetString("url"),
		MaxPingCount: viper.GetInt("max_ping_count"),
	}
	return &cfg
}

func Init() (*gin.Engine, error) {

	// init ConfigData properties
	if err := initConfig(); err != nil {
		return nil, err
	}

	// init log package
	initLog()
	// monitor configuration file changes and hot load the program
	watchConfig()

	return initGin(), nil
}

func initConfig() error {
	if Cfg != "" {
		viper.SetConfigFile(Cfg)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("QBUS_MANAGER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	DataYaml = parseConfig()

	return nil
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed：%s", e.Name)
	})
}
