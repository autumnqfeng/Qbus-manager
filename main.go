package main

import (
	"net/http"
	"os"
	"os/signal"

	"Qbus-manager/config"
	"Qbus-manager/pkg/ping"
	"Qbus-manager/pkg/version"
	"Qbus-manager/zookeeper"

	"github.com/lexkong/log"
)

func main() {
	config.Parse()
	// print system version
	if version.PrintVersion() {
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Init
	g, err := config.Init()
	if err != nil {
		panic(err)
	}

	if err := zookeeper.Connect(config.ZookeeperURL); err != nil {
		log.Errorf(err, "zk_connect")
		os.Exit(1)
		return
	}

	// Ping the server to make sure the router is working.
	go ping.Start()

	log.Infof("Start to listening the incoming requests on http address: %s", config.DataYaml.Addr)
	log.Info(http.ListenAndServe(config.DataYaml.Addr, g).Error())

	<-signals
	zookeeper.Stop()
}
