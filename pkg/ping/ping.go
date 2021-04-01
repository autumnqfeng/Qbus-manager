package ping

import (
	"errors"
	"net/http"
	"time"

	"qbus-manager/configs"

	"go.uber.org/zap"
)

// Ping the server to make sure the router is working.
func Start() {
	if err := pingServer(); err != nil {
		zap.L().Fatal("The router has no response, or it might took too long to start up.", zap.Error(err))
	}
	zap.L().Info("The router has been deployed successfully.")
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < configs.Conf.MaxPingCount; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(configs.Conf.Url + "/check/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		zap.L().Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("can not connect to the router")
}
