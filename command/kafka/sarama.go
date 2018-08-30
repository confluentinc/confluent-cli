package kafka

import (
	"strings"

	"github.com/Shopify/sarama"

	"github.com/confluentinc/cli/shared"
)

// NewSaramaKafka returns a sarama.Client configured for the KafkaClusterConfig
func NewSaramaKafka(kafka shared.KafkaClusterConfig) (sarama.Client, error) {
	return sarama.NewClient(strings.Split(kafka.Bootstrap, ","), saramaConf(kafka))
}

// NewSaramaAdmin returns a sarama.ClusterAdmin configured for the KafkaClusterConfig
func NewSaramaAdmin(kafka shared.KafkaClusterConfig) (sarama.ClusterAdmin, error) {
	return sarama.NewClusterAdmin(strings.Split(kafka.Bootstrap, ","), saramaConf(kafka))
}

// NewSaramaKafkaForConfig returns a sarama.Client configured for the CLI config
func NewSaramaKafkaForConfig(config *shared.Config) (sarama.Client, error) {
	cluster, err := config.KafkaClusterConfig()
	if err != nil {
		return nil, err
	}
	return NewSaramaKafka(cluster)
}

// NewSaramaAdminForConfig returns a sarama.ClusterAdmin configured for the CLI config
func NewSaramaAdminForConfig(config *shared.Config) (sarama.ClusterAdmin, error) {
	cluster, err := config.KafkaClusterConfig()
	if err != nil {
		return nil, err
	}
	return NewSaramaAdmin(cluster)
}

func saramaConf(kafka shared.KafkaClusterConfig) *sarama.Config {
	saramaConf := sarama.NewConfig()
	saramaConf.Version = sarama.V1_1_0_0
	saramaConf.Net.TLS.Enable = true
	saramaConf.Net.SASL.Enable = true
	saramaConf.Net.SASL.User = kafka.APIKey
	saramaConf.Net.SASL.Password = kafka.APISecret
	return saramaConf
}
