package kafka

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func adminClient() (*kafka.AdminClient, error) {
	adminConfig := &kafka.ConfigMap{"bootstrap.servers": strings.Join(Brokers.KafkaUrlList(), ",")}
	return kafka.NewAdminClient(adminConfig)
}
