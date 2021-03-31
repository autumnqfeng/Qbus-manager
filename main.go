package main

import (
	"net/http"
	"os"
	"os/signal"

	"qbus-manager/_init"
	"qbus-manager/pkg/ping"
	"qbus-manager/pkg/version"
	"qbus-manager/pkg/zookeeper"

	"github.com/lexkong/log"
)

func main() {
	_init.Parse()
	// print system version
	if version.PrintVersion() {
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// SystemInit
	g, err := _init.SystemInit()
	if err != nil {
		panic(err)
	}

	if err := zookeeper.Init(_init.ZookeeperURL); err != nil {
		log.Errorf(err, "zk_connect")
		os.Exit(1)
		return
	}

	// Ping the server to make sure the router is working.
	go ping.Start()

	log.Infof("Start to listening the incoming requests on http address: %s", _init.DataYaml.Addr)
	log.Info(http.ListenAndServe(_init.DataYaml.Addr, g).Error())

	<-signals
	zookeeper.StopAll()
}
