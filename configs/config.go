package configs

import "github.com/spf13/viper"

type Config struct {
	Mode         string `yaml:"mode"`
	Port         string `yaml:"port"`
	Name         string `yaml:"name"`
	Url          string `yaml:"url"`
	MaxPingCount int    `yaml:"max_ping_count"`
	*LogConfig   `yaml:"log"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxsize"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}

var Conf *Config

func ParseConfig() *Config {
	return &Config{
		Mode:         viper.GetString("mode"),
		Port:         viper.GetString("port"),
		Name:         viper.GetString("name"),
		Url:          viper.GetString("url"),
		MaxPingCount: viper.GetInt("max_ping_count"),
		LogConfig: &LogConfig{
			Level:      viper.GetString("log.level"),
			Filename:   viper.GetString("log.filename"),
			MaxSize:    viper.GetInt("log.maxsize"),
			MaxAge:     viper.GetInt("log.max_age"),
			MaxBackups: viper.GetInt("log.max_backups"),
		},
	}
}
