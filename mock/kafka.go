package mock

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/confluentinc/ccloud-sdk-go"
	authv1 "github.com/confluentinc/ccloudapis/auth/v1"
	kafkav1 "github.com/confluentinc/ccloudapis/kafka/v1"
)

// Compile-time check interface adherence
var _ ccloud.Kafka = (*Kafka)(nil)

type Kafka struct {
	Expect chan interface{}
}

func NewKafkaMock(expect chan interface{}) *Kafka {
	return &Kafka{expect}
}

func (m *Kafka) CreateAPIKey(_ context.Context, apiKey *authv1.ApiKey) (*authv1.ApiKey, error) {
	return apiKey, nil
}

func (m *Kafka) List(_ context.Context, cluster *kafkav1.KafkaCluster) ([]*kafkav1.KafkaCluster, error) {
	return []*kafkav1.KafkaCluster{cluster}, nil
}

func (m *Kafka) Describe(_ context.Context, cluster *kafkav1.KafkaCluster) (*kafkav1.KafkaCluster, error) {
	return cluster, nil
}

func (m *Kafka) Create(_ context.Context, config *kafkav1.KafkaClusterConfig) (*kafkav1.KafkaCluster, error) {
	return &kafkav1.KafkaCluster{}, nil
}

func (m *Kafka) Delete(_ context.Context, cluster *kafkav1.KafkaCluster) error {
	return nil
}

func (m *Kafka) ListTopics(ctx context.Context, cluster *kafkav1.KafkaCluster) ([]*kafkav1.TopicDescription, error) {
	return []*kafkav1.TopicDescription{
		{Name: "test1"},
		{Name: "test2"},
		{Name: "test3"}}, nil
}

func (m *Kafka) DescribeTopic(ctx context.Context, cluster *kafkav1.KafkaCluster, topic *kafkav1.Topic) (*kafkav1.TopicDescription, error) {
	node := &kafkav1.KafkaNode{Id: 1}
	tp := &kafkav1.TopicPartitionInfo{Leader: node, Replicas: []*kafkav1.KafkaNode{node}}
	return &kafkav1.TopicDescription{Partitions: []*kafkav1.TopicPartitionInfo{tp}},
		assertEquals(topic, <-m.Expect)
}

func (m *Kafka) CreateTopic(ctx context.Context, cluster *kafkav1.KafkaCluster, topic *kafkav1.Topic) error {
	return assertEquals(topic, <-m.Expect)
}

func (m *Kafka) DeleteTopic(ctx context.Context, cluster *kafkav1.KafkaCluster, topic *kafkav1.Topic) error {
	return assertEquals(topic, <-m.Expect)
}

func (m *Kafka) ListTopicConfig(ctx context.Context, cluster *kafkav1.KafkaCluster, topic *kafkav1.Topic) (*kafkav1.TopicConfig, error) {
	return nil, assertEquals(topic, <-m.Expect)
}

func (m *Kafka) UpdateTopic(ctx context.Context, cluster *kafkav1.KafkaCluster, topic *kafkav1.Topic) error {
	return assertEquals(topic, <-m.Expect)
}

func (m *Kafka) ListACL(ctx context.Context, cluster *kafkav1.KafkaCluster, filter *kafkav1.ACLFilter) ([]*kafkav1.ACLBinding, error) {
	expect := <-m.Expect
	if filter.PatternFilter.PatternType == kafkav1.PatternTypes_ANY {
		expect.(*kafkav1.ACLFilter).PatternFilter.PatternType = kafkav1.PatternTypes_ANY
	}
	if filter.EntryFilter.Operation == kafkav1.ACLOperations_ANY {
		expect.(*kafkav1.ACLFilter).EntryFilter.Operation = kafkav1.ACLOperations_ANY
	}
	if filter.EntryFilter.PermissionType == kafkav1.ACLPermissionTypes_ANY {
		expect.(*kafkav1.ACLFilter).EntryFilter.PermissionType = kafkav1.ACLPermissionTypes_ANY
	}
	return nil, assertEquals(filter, expect)
}

func (m *Kafka) CreateACL(ctx context.Context, cluster *kafkav1.KafkaCluster, binding []*kafkav1.ACLBinding) error {
	return assertEquals(binding[0], <-m.Expect)
}

func (m *Kafka) DeleteACL(ctx context.Context, cluster *kafkav1.KafkaCluster, filter *kafkav1.ACLFilter) error {
	return assertEquals(filter, <-m.Expect)
}

func assertEquals(actual interface{}, expected interface{}) error {
	actualMessage := actual.(proto.Message)
	expectedMessage := expected.(proto.Message)

	if !proto.Equal(actualMessage, expectedMessage) {
		return fmt.Errorf("actual: %+v\nexpected: %+v", actual, expected)
	}
	return nil
}