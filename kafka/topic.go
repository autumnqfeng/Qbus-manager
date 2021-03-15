package kafka

import (
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
