package topic

import (
	"go.uber.org/zap"
	"strconv"

	"qbus-manager/handler"
	"qbus-manager/pkg/errno"
	"qbus-manager/pkg/kafka"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	zap.L().Info("Topic Create function called.")
	var t CreateTopicVo
	if err := c.Bind(&t); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	config := make(map[string]string)
	config["retention.ms"] = strconv.Itoa(t.Retention)
	config["max.message.bytes"] = strconv.Itoa(t.MaxMessage)

	if err := kafka.CreateTopic(t.Topic, t.Partitions, t.Replications, config); err != nil {
		handler.SendResponse(c, errno.ErrCreateTopic, nil)
		return
	}
	handler.SendResponse(c, errno.OK, nil)
}

func Delete(c *gin.Context) {
	c.Param("clustername")
	c.Param("topic")

}
