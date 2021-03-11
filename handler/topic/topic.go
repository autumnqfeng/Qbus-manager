package topic

import (
	"Qbus-manager/handler"
	"Qbus-manager/pkg/errno"
	"Qbus-manager/store"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func Create(c *gin.Context) {
	log.Info("Topic Create function called.")
	var t CreateTopic
	if err :=c.Bind(&t); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
	}

	config := make(map[string]string)
	config["retention.ms"] = strconv.Itoa(t.Retention)
	config["max.message.bytes"] = strconv.Itoa(t.MaxMessage)

	if err := store.CreateTopic(t.Topic, t.Partitions, t.Replications, config); err != nil {
		handler.SendResponse(c, errno.ErrCreateTopic, nil)
	}
	handler.SendResponse(c, errno.OK, nil)
}

func Delete(c *gin.Context) {
	c.Param("clustername")
	c.Param("topic")



}
