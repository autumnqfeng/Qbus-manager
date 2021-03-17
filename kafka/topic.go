package kafka

import (
	"Qbus-manager/pkg/errno"
	"Qbus-manager/zookeeper"
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lexkong/log"
	"time"
)

const timeOut = 15 * time.Second

func CreateTopic(name string, partitions, replicationFactor int, topicConfig map[string]string) error {
	a, err := adminClient()
	if err != nil {
		log.Errorf(err, "create_topic")
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results, err := a.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:             name,
			NumPartitions:     partitions,
			ReplicationFactor: replicationFactor,
			Config:            topicConfig,
		}},
		kafka.SetAdminOperationTimeout(timeOut))

	if err != nil {
		log.Errorf(err, "create_topic")
		return err
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			log.Errorf(result.Error, "create_topic")
			return result.Error
		}
	}

	return nil
}

func DeleteTopic(name string) error {
	a, err := adminClient()
	if err != nil {
		log.Errorf(err, "delete_topic")
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results, err := a.DeleteTopics(
		ctx,
		[]string{name},
		kafka.SetAdminOperationTimeout(timeOut))

	if err != nil {
		log.Errorf(err, "delete_topic")
		return err
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			log.Errorf(result.Error, "delete_topic")
			return result.Error
		}
	}

	return nil
}

func GetTopicByCluster(cc *zookeeper.ClusterConfig) ([]string, error) {
	clusterName := cc.ClusterName

	conn, err := zookeeper.GetConn(clusterName)
	if err != nil {
		log.Errorf(errno.ErrClusterConnect, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrClusterConnect.Code, errno.ErrClusterConnect.Message, clusterName)
		return nil, errno.ErrClusterConnect
	}

	data, _, err := conn.Children("/brokers/topics")
	if err != nil {
		log.Errorf(errno.ErrGetBroker, "err_code: `%v`, err_msg: `%v`, clusterName: `%v`", errno.ErrGetBroker.Code, errno.ErrGetBroker.Message, clusterName)
		return nil, errno.ErrGetBroker
	}
	return data, nil
}
