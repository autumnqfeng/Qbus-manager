package _init

import (
	"fmt"
	"go.uber.org/zap"
	"qbus-manager/configs"
	"qbus-manager/pkg/logger"
	"qbus-manager/pkg/zookeeper"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Init() (*gin.Engine, error) {

	// _init ConfigData properties
	if err := viperInit(); err != nil {
		return nil, err
	}

	// _init log package
	//logInit()

	if err := logger.InitLogger(configs.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return nil, err
	}

	// monitor configuration file changes and hot load the program
	watchConfig()

	if err := zookeeper.Init(ZookeeperURL); err != nil {
		zap.L().Error("zk_connect filed", zap.Error(err))
		return nil, err
	}

	return ginInit(), nil
}

func viperInit() error {
	if Cfg != "" {
		viper.SetConfigFile(Cfg)
	} else {
		viper.AddConfigPath("configs")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("QBUS_MANAGER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	configs.Conf = configs.ParseConfig()

	return nil
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		zap.L().Info("Config file changed", zap.String("event_name", e.Name))
	})
}
