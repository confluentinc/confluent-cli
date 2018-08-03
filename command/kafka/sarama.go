package kafka

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"

	"github.com/confluentinc/cli/shared"
)

// NewSaramaKafka returns a sarama.Client configured for the KafkaCluster
func NewSaramaKafka(kafka shared.KafkaCluster) (sarama.Client, error) {
	return sarama.NewClient(strings.Split(kafka.Bootstrap, ","), saramaConf(kafka))
}

// NewSaramaAdmin returns a sarama.ClusterAdmin configured for the KafkaCluster
func NewSaramaAdmin(kafka shared.KafkaCluster) (sarama.ClusterAdmin, error) {
	return sarama.NewClusterAdmin(strings.Split(kafka.Bootstrap, ","), saramaConf(kafka))
}

// NewSaramaKafkaForConfig returns a sarama.Client configured for the CLI config
func NewSaramaKafkaForConfig(config *shared.Config) (sarama.Client, error) {
	cluster, err := kafkaCluster(config)
	if err != nil {
		return nil, err
	}
	return NewSaramaKafka(cluster)
}

// NewSaramaAdminForConfig returns a sarama.ClusterAdmin configured for the CLI config
func NewSaramaAdminForConfig(config *shared.Config) (sarama.ClusterAdmin, error) {
	cluster, err := kafkaCluster(config)
	if err != nil {
		return nil, err
	}
	return NewSaramaAdmin(cluster)
}

func kafkaCluster(config *shared.Config) (shared.KafkaCluster, error) {
	cfg, err := config.Context()
	if err != nil {
		return shared.KafkaCluster{}, err
	}
	cluster, found := config.Platforms[cfg.Platform].KafkaClusters[cfg.Kafka]
	if !found {
		e := fmt.Errorf("no auth found for Kafka %s, please run `confluent kafka cluster auth` first", cfg.Kafka)
		return shared.KafkaCluster{}, shared.NotAuthenticatedError(e)
	}
	return cluster, nil
}

func saramaConf(kafka shared.KafkaCluster) *sarama.Config {
	saramaConf := sarama.NewConfig()
	saramaConf.Version = sarama.V1_1_0_0
	saramaConf.Net.TLS.Enable = true
	saramaConf.Net.SASL.Enable = true
	saramaConf.Net.SASL.User = kafka.APIKey
	saramaConf.Net.SASL.Password = kafka.APISecret
	return saramaConf
}
