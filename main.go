package main

import (
	"net/http"
	"os"
	"os/signal"

	"qbus-manager/_init"
	"qbus-manager/configs"
	"qbus-manager/pkg/ping"
	"qbus-manager/pkg/version"
	"qbus-manager/pkg/zookeeper"

	"go.uber.org/zap"
)

func main() {
	_init.Parse()
	// print system version
	if version.PrintVersion() {
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Init
	g, err := _init.Init()
	if err != nil {
		panic(err)
	}

	// Ping the server to make sure the router is working.
	go ping.Start()

	zap.L().Info("Start to listening the incoming requests on http address.", zap.String("port", configs.Conf.Port))
	zap.L().Info(http.ListenAndServe(configs.Conf.Port, g).Error())

	<-signals
	zookeeper.StopAll()
}
