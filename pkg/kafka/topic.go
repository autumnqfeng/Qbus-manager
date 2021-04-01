package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const timeOut = 15 * time.Second

func CreateTopic(name string, partitions, replicationFactor int, topicConfig map[string]string) error {
	a, err := adminClient()
	if err != nil {
		zap.L().Error(fmt.Sprintf("create_topic"), zap.Error(err))
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
		zap.L().Error(fmt.Sprintf("create_topic"), zap.Error(err))
		return err
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			zap.L().Error(fmt.Sprintf("create_topic"), zap.Error(result.Error))
			return result.Error
		}
	}

	return nil
}

func DeleteTopic(name string) error {
	a, err := adminClient()
	if err != nil {
		zap.L().Error(fmt.Sprintf("delete_topic"), zap.Error(err))
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results, err := a.DeleteTopics(
		ctx,
		[]string{name},
		kafka.SetAdminOperationTimeout(timeOut))

	if err != nil {
		zap.L().Error(fmt.Sprintf("delete_topic"), zap.Error(err))
		return err
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			zap.L().Error(fmt.Sprintf("delete_topic"), zap.Error(result.Error))
			return result.Error
		}
	}

	return nil
}
