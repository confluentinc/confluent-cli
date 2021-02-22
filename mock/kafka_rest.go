package mock

import (
	"context"
	"fmt"

	nethttp "net/http"

	krsdk "github.com/confluentinc/kafka-rest-sdk-go/kafkarestv3"
)

// Compile-time check interface adherence
var _ krsdk.TopicApi = (*Topic)(nil)

type Topic struct {
}

func NewTopicMock() *Topic {
	return &Topic{}
}

func (m *Topic) ClustersClusterIdTopicsGet(_ctx context.Context, _clusterId string) (krsdk.TopicDataList, *nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: 200,
	}
	return krsdk.TopicDataList{
		Kind:     "",
		Metadata: krsdk.ResourceCollectionMetadata{},
		Data: []krsdk.TopicData{
			{
				Kind:                   "",
				Metadata:               krsdk.ResourceMetadata{},
				ClusterId:              _clusterId,
				TopicName:              "NAME",
				IsInternal:             false,
				ReplicationFactor:      0,
				Partitions:             krsdk.Relationship{},
				Configs:                krsdk.Relationship{},
				PartitionReassignments: krsdk.Relationship{},
			},
		},
	}, httpResp, nil
}

func (m *Topic) ClustersClusterIdTopicsPost(_ctx context.Context, _clusterId string, _localVarOptionals *krsdk.ClustersClusterIdTopicsPostOpts) (krsdk.TopicData, *nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: 201,
	}
	return krsdk.TopicData{}, httpResp, nil
}

func (m *Topic) ClustersClusterIdTopicsTopicNameDelete(_ctx context.Context, _clusterId string, _topicName string) (*nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: 204,
	}
	return httpResp, nil
}

func (m *Topic) ClustersClusterIdTopicsTopicNameGet(_ctx context.Context, _clusterId string, _topicName string) (krsdk.TopicData, *nethttp.Response, error) {
	return krsdk.TopicData{}, nil, nil
}

// Compile-time check interface adherence
var _ krsdk.ACLApi = (*ACL)(nil)

type ACL struct {
}

func NewACLMock() *ACL {
	return &ACL{}
}

func (m *ACL) ClustersClusterIdAclsDelete(_ctx context.Context, _clusterId string, _localVarOptionals *krsdk.ClustersClusterIdAclsDeleteOpts) (krsdk.InlineResponse200, *nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: 200,
	}
	return krsdk.InlineResponse200{
		Data: []krsdk.AclData{
			{},
		},
	}, httpResp, nil
}

func (m *ACL) ClustersClusterIdAclsGet(_ctx context.Context, _clusterId string, _localVarOptionals *krsdk.ClustersClusterIdAclsGetOpts) (krsdk.AclDataList, *nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: 200,
	}
	return krsdk.AclDataList{
		Kind:     "",
		Metadata: krsdk.ResourceCollectionMetadata{},
		Data: []krsdk.AclData{
			{
				Kind:         "KIND",
				Metadata:     krsdk.ResourceMetadata{},
				ClusterId:    _clusterId,
				ResourceType: "TYPE",
				ResourceName: "NAME",
				PatternType:  "PATTERN",
				Principal:    "PRINCIPAL",
				Host:         "HOST",
				Operation:    "OP",
				Permission:   "PERMISSION",
			},
		},
	}, httpResp, nil
}

func (m *ACL) ClustersClusterIdAclsPost(_ctx context.Context, _clusterId string, _localVarOptionals *krsdk.ClustersClusterIdAclsPostOpts) (*nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: 201,
	}
	return httpResp, nil
}

// Compile-time check interface adherence
var _ krsdk.PartitionApi = (*Partition)(nil)

type Partition struct {
}

func NewPartitionMock() *Partition {
	return &Partition{}
}

func (m *Partition) ClustersClusterIdTopicsPartitionsReassignmentGet(_ctx context.Context, _clusterId string) (krsdk.ReassignmentDataList, *nethttp.Response, error) {
	return krsdk.ReassignmentDataList{}, nil, nil
}

func (m *Partition) ClustersClusterIdTopicsTopicNamePartitionsGet(_ctx context.Context, clusterId string, topicName string) (krsdk.PartitionDataList, *nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: 200,
	}
	return krsdk.PartitionDataList{
		Kind:     "",
		Metadata: krsdk.ResourceCollectionMetadata{},
		Data: []krsdk.PartitionData{
			{
				Kind:         "",
				Metadata:     krsdk.ResourceMetadata{},
				ClusterId:    clusterId,
				TopicName:    topicName,
				PartitionId:  0,
				Leader:       krsdk.Relationship{},
				Replicas:     krsdk.Relationship{},
				Reassignment: krsdk.Relationship{},
			},
		},
	}, httpResp, nil
}

func (m *Partition) ClustersClusterIdTopicsTopicNamePartitionsPartitionIdGet(_ctx context.Context, _clusterId string, _topicName string, _partitionId int32) (krsdk.PartitionData, *nethttp.Response, error) {
	return krsdk.PartitionData{}, nil, nil
}

func (m *Partition) ClustersClusterIdTopicsTopicNamePartitionsPartitionIdReassignmentGet(_ctx context.Context, _clusterId string, _topicName string, _partitionId int32) (krsdk.ReassignmentData, *nethttp.Response, error) {
	return krsdk.ReassignmentData{}, nil, nil
}

func (m *Partition) ClustersClusterIdTopicsTopicNamePartitionsReassignmentGet(_ctx context.Context, _clusterId string, _topicName string) (krsdk.ReassignmentDataList, *nethttp.Response, error) {
	return krsdk.ReassignmentDataList{}, nil, nil
}

// Compile-time check interface adherence
var _ krsdk.ReplicaApi = (*Replica)(nil)

type Replica struct {
}

func NewReplicaMock() *Replica {
	return &Replica{}
}

func (m *Replica) ClustersClusterIdBrokersBrokerIdPartitionReplicasGet(_ctx context.Context, _clusterId string, _brokerId int32) (krsdk.ReplicaDataList, *nethttp.Response, error) {
	return krsdk.ReplicaDataList{}, nil, nil
}

func (m *Replica) ClustersClusterIdTopicsTopicNamePartitionsPartitionIdReplicasBrokerIdGet(_ctx context.Context, _clusterId string, _topicName string, partitionId int32, brokerId int32) (krsdk.ReplicaData, *nethttp.Response, error) {
	return krsdk.ReplicaData{}, nil, nil
}

func (m *Replica) ClustersClusterIdTopicsTopicNamePartitionsPartitionIdReplicasGet(_ctx context.Context, clusterId string, topicName string, partitionId int32) (krsdk.ReplicaDataList, *nethttp.Response, error) {
	return krsdk.ReplicaDataList{
		Kind:     "",
		Metadata: krsdk.ResourceCollectionMetadata{},
		Data: []krsdk.ReplicaData{
			{
				Kind:        "",
				Metadata:    krsdk.ResourceMetadata{},
				ClusterId:   clusterId,
				TopicName:   topicName,
				PartitionId: partitionId,
				BrokerId:    42,
				IsLeader:    true,
				IsInSync:    true,
				Broker:      krsdk.Relationship{},
			},
		},
	}, nil, nil
}

// Compile-time check interface adherence
var _ krsdk.ConfigsApi = (*Configs)(nil)

type Configs struct {
}

func NewConfigsMock() *Configs {
	return &Configs{}
}

func (m *Configs) ClustersClusterIdBrokerConfigsGet(_ctx context.Context, _clusterId string) (krsdk.ClusterConfigDataList, *nethttp.Response, error) {
	return krsdk.ClusterConfigDataList{}, nil, nil
}

func (m *Configs) ClustersClusterIdBrokerConfigsNameDelete(_ctx context.Context, _clusterId string, name string) (*nethttp.Response, error) {
	return nil, nil
}

func (m *Configs) ClustersClusterIdBrokerConfigsNameGet(_ctx context.Context, _clusterId string, name string) (krsdk.ClusterConfigData, *nethttp.Response, error) {
	return krsdk.ClusterConfigData{}, nil, nil
}

func (m *Configs) ClustersClusterIdBrokerConfigsNamePut(_ctx context.Context, _clusterId string, name string, localVarOptionals *krsdk.ClustersClusterIdBrokerConfigsNamePutOpts) (*nethttp.Response, error) {
	return nil, nil
}

func (m *Configs) ClustersClusterIdBrokerConfigsalterPost(_ctx context.Context, _clusterId string, localVarOptionals *krsdk.ClustersClusterIdBrokerConfigsalterPostOpts) (*nethttp.Response, error) {
	return nil, nil
}

func (m *Configs) ClustersClusterIdBrokersBrokerIdConfigsGet(_ctx context.Context, _clusterId string, brokerId int32) (krsdk.BrokerConfigDataList, *nethttp.Response, error) {
	return krsdk.BrokerConfigDataList{}, nil, nil
}

func (m *Configs) ClustersClusterIdBrokersBrokerIdConfigsNameDelete(_ctx context.Context, _clusterId string, brokerId int32, name string) (*nethttp.Response, error) {
	return nil, nil
}

func (m *Configs) ClustersClusterIdBrokersBrokerIdConfigsNameGet(_ctx context.Context, _clusterId string, brokerId int32, name string) (krsdk.BrokerConfigData, *nethttp.Response, error) {
	return krsdk.BrokerConfigData{}, nil, nil
}

func (m *Configs) ClustersClusterIdBrokersBrokerIdConfigsNamePut(_ctx context.Context, _clusterId string, brokerId int32, name string, localVarOptionals *krsdk.ClustersClusterIdBrokersBrokerIdConfigsNamePutOpts) (*nethttp.Response, error) {
	return nil, nil
}

func (m *Configs) ClustersClusterIdBrokersBrokerIdConfigsalterPost(_ctx context.Context, _clusterId string, brokerId int32, localVarOptionals *krsdk.ClustersClusterIdBrokersBrokerIdConfigsalterPostOpts) (*nethttp.Response, error) {
	return nil, nil
}

func (m *Configs) ClustersClusterIdTopicsTopicNameConfigsGet(_ctx context.Context, _clusterId string, topicName string) (krsdk.TopicConfigDataList, *nethttp.Response, error) {
	v := "configValue1"
	return krsdk.TopicConfigDataList{
		Kind:     "",
		Metadata: krsdk.ResourceCollectionMetadata{},
		Data: []krsdk.TopicConfigData{
			{
				Kind:        "",
				Metadata:    krsdk.ResourceMetadata{},
				ClusterId:   _clusterId,
				Name:        "ConfigProperty1",
				Value:       &v,
				IsDefault:   false,
				IsReadOnly:  false,
				IsSensitive: false,
				Source:      "",
				Synonyms:    []krsdk.ConfigSynonymData{},
				TopicName:   topicName,
			},
		},
	}, nil, nil
}

func (m *Configs) ClustersClusterIdTopicsTopicNameConfigsNameDelete(_ctx context.Context, _clusterId string, topicName string, name string) (*nethttp.Response, error) {
	return nil, nil
}

func (m *Configs) ClustersClusterIdTopicsTopicNameConfigsNameGet(_ctx context.Context, _clusterId string, topicName string, name string) (krsdk.TopicConfigData, *nethttp.Response, error) {
	return krsdk.TopicConfigData{}, nil, nil
}

func (m *Configs) ClustersClusterIdTopicsTopicNameConfigsNamePut(_ctx context.Context, _clusterId string, topicName string, name string, localVarOptionals *krsdk.ClustersClusterIdTopicsTopicNameConfigsNamePutOpts) (*nethttp.Response, error) {
	return nil, nil
}

func (m *Configs) ClustersClusterIdTopicsTopicNameConfigsalterPost(_ctx context.Context, _clusterId string, topicName string, localVarOptionals *krsdk.ClustersClusterIdTopicsTopicNameConfigsalterPostOpts) (*nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: 204,
	}
	return httpResp, nil
}

// Compile-time check interface adherence
var _ krsdk.ClusterLinkingApi = (*ClusterLinking)(nil)

type ClusterLinking struct {
}

func NewClusterLinkingMock() *ClusterLinking {
	return &ClusterLinking{}
}

func (m *ClusterLinking) ClustersClusterIdLinksGet(ctx context.Context, clusterId string) (krsdk.ListLinksResponseDataList, *nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusOK,
	}
	return krsdk.ListLinksResponseDataList{
		Kind:     "",
		Metadata: krsdk.ResourceCollectionMetadata{},
		Data: []krsdk.ListLinksResponseData{
			{
				Kind:        "",
				Metadata:    krsdk.ResourceMetadata{},
				ClusterId:   clusterId,
				LinkName:    "link-1",
				LinkId:      "LinkId",
				TopicsNames: []string{"topic-1", "topic-2", "topic-3"},
			},
		},
	}, httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameConfigsConfigNameDelete(ctx context.Context, clusterId string, linkName string, configName string) (*nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusNoContent,
	}
	return httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameConfigsConfigNameGet(ctx context.Context, clusterId string, linkName string, configName string) (krsdk.ListLinkConfigsResponseData, *nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusOK,
	}
	return krsdk.ListLinkConfigsResponseData{
		Kind:      "",
		Metadata:  krsdk.ResourceMetadata{},
		ClusterId: clusterId,
		Name:      configName,
		Value:     "value",
		ReadOnly:  false,
		Sensitive: false,
		Source:    "DYNAMIC_CLUSTER_LINK_CONFIG",
		Synonyms:  nil,
		LinkName:  linkName,
	}, httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameConfigsConfigNamePut(ctx context.Context, clusterId string, linkName string, configName string, localVarOptionals *krsdk.ClustersClusterIdLinksLinkNameConfigsConfigNamePutOpts) (*nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusNoContent,
	}
	return httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameConfigsGet(ctx context.Context, clusterId string, linkName string) (krsdk.ListLinkConfigsResponseDataList, *nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusOK,
	}
	return krsdk.ListLinkConfigsResponseDataList{
		Kind:     "",
		Metadata: krsdk.ResourceCollectionMetadata{},
		Data: []krsdk.ListLinkConfigsResponseData{
			{
				Kind:      "",
				Metadata:  krsdk.ResourceMetadata{},
				ClusterId: clusterId,
				Name:      "consumer.offset.sync.enable",
				Value:     "value",
				ReadOnly:  false,
				Sensitive: false,
				Source:    "DYNAMIC_CLUSTER_LINK_CONFIG",
				Synonyms:  nil,
				LinkName:  linkName,
			},
		},
	}, httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameConfigsalterPut(ctx context.Context, clusterId string, linkName string, localVarOptionals *krsdk.ClustersClusterIdLinksLinkNameConfigsalterPutOpts) (*nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusNoContent,
	}
	return httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameDelete(ctx context.Context, clusterId string, linkName string) (*nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusNoContent,
	}
	return httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameGet(ctx context.Context, clusterId string, linkName string) (krsdk.ListLinksResponseData, *nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusOK,
	}
	return krsdk.ListLinksResponseData{
		Kind:        "",
		Metadata:    krsdk.ResourceMetadata{},
		ClusterId:   clusterId,
		LinkName:    linkName,
		LinkId:      "link-1",
		TopicsNames: []string{"topic-1", "topic-2", "topic-3"},
	}, httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameMirrorsDestinationTopicNameGet(ctx context.Context, clusterId string, linkName string, destinationTopicName string) (krsdk.ListMirrorTopicsResponseData, *nethttp.Response, error) {
	// TODO: Remove std output before merging
	fmt.Println(destinationTopicName)

	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusOK,
	}
	return krsdk.ListMirrorTopicsResponseData{
		Kind:                 "",
		Metadata:             krsdk.ResourceMetadata{},
		DestinationTopicName: "dest-topic-1",
		SourceTopicName:      "src-topic-1",
		MirrorTopicStatus:    "active",
		StateTimeMs:          44444444,
	}, httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameMirrorsFailoverPost(ctx context.Context, clusterId string, linkName string, localVarOptionals *krsdk.ClustersClusterIdLinksLinkNameMirrorsFailoverPostOpts) (krsdk.AlterMirrorStatusResponseDataList, *nethttp.Response, error) {
	// TODO: Remove std output before merging
	fmt.Println(localVarOptionals.ValidateOnly.Value())
	fmt.Println(localVarOptionals.AlterMirrorsRequestData.Value().(krsdk.AlterMirrorsRequestData).DestinationTopicNames)

	return m.AlterMirrorResultResponse()
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameMirrorsGet(ctx context.Context, clusterId string, linkName string, localVarOptionals *krsdk.ClustersClusterIdLinksLinkNameMirrorsGetOpts) (krsdk.ListMirrorTopicsResponseDataList, *nethttp.Response, error) {
	// TODO: Remove std output before merging
	fmt.Println(localVarOptionals.MirrorStatus.Value().(krsdk.MirrorTopicStatus))

	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusOK,
	}
	return krsdk.ListMirrorTopicsResponseDataList{
		Kind:     "",
		Metadata: krsdk.ResourceCollectionMetadata{},
		Data: []krsdk.ListMirrorTopicsResponseData{
			{
				Kind:                 "",
				Metadata:             krsdk.ResourceMetadata{},
				DestinationTopicName: "dest-topic-1",
				SourceTopicName:      "src-topic-1",
				MirrorTopicStatus:    "active",
				StateTimeMs:          44444444,
			},
			{
				Kind:                 "",
				Metadata:             krsdk.ResourceMetadata{},
				DestinationTopicName: "dest-topic-2",
				SourceTopicName:      "src-topic-2",
				MirrorTopicStatus:    "active",
				StateTimeMs:          55555555,
			},
		},
	}, httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameMirrorsPausePost(ctx context.Context, clusterId string, linkName string, localVarOptionals *krsdk.ClustersClusterIdLinksLinkNameMirrorsPausePostOpts) (krsdk.AlterMirrorStatusResponseDataList, *nethttp.Response, error) {
	// TODO: Remove std output before merging
	fmt.Println(localVarOptionals.ValidateOnly.Value())
	fmt.Println(localVarOptionals.AlterMirrorsRequestData.Value().(krsdk.AlterMirrorsRequestData).DestinationTopicNames)

	return m.AlterMirrorResultResponse()
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameMirrorsPost(ctx context.Context, clusterId string, linkName string, localVarOptionals *krsdk.ClustersClusterIdLinksLinkNameMirrorsPostOpts) (*nethttp.Response, error) {
	// TODO: Remove std output before merging
	fmt.Println(localVarOptionals.CreateMirrorTopicRequestData.Value().(krsdk.CreateMirrorTopicRequestData).SourceTopicName)
	fmt.Println(localVarOptionals.CreateMirrorTopicRequestData.Value().(krsdk.CreateMirrorTopicRequestData).ReplicationFactor)
	fmt.Println(localVarOptionals.CreateMirrorTopicRequestData.Value().(krsdk.CreateMirrorTopicRequestData).Configs[0].Name)
	fmt.Println(localVarOptionals.CreateMirrorTopicRequestData.Value().(krsdk.CreateMirrorTopicRequestData).Configs[0].Value)
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusNoContent,
	}
	return httpResp, nil
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameMirrorsPromotePost(ctx context.Context, clusterId string, linkName string, localVarOptionals *krsdk.ClustersClusterIdLinksLinkNameMirrorsPromotePostOpts) (krsdk.AlterMirrorStatusResponseDataList, *nethttp.Response, error) {
	// TODO: Remove std output before merging
	fmt.Println(localVarOptionals.ValidateOnly.Value())
	fmt.Println(localVarOptionals.AlterMirrorsRequestData.Value().(krsdk.AlterMirrorsRequestData).DestinationTopicNames)

	return m.AlterMirrorResultResponse()
}

func (m *ClusterLinking) ClustersClusterIdLinksLinkNameMirrorsResumePost(ctx context.Context, clusterId string, linkName string, localVarOptionals *krsdk.ClustersClusterIdLinksLinkNameMirrorsResumePostOpts) (krsdk.AlterMirrorStatusResponseDataList, *nethttp.Response, error) {
	// TODO: Remove std output before merging
	fmt.Println(localVarOptionals.ValidateOnly.Value())
	fmt.Println(localVarOptionals.AlterMirrorsRequestData.Value().(krsdk.AlterMirrorsRequestData).DestinationTopicNames)

	return m.AlterMirrorResultResponse()
}

func (m *ClusterLinking) ClustersClusterIdLinksPost(ctx context.Context, clusterId string, linkName string, localVarOptionals *krsdk.ClustersClusterIdLinksPostOpts) (*nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusCreated,
	}
	return httpResp, nil
}

func (m *ClusterLinking) AlterMirrorResultResponse() (krsdk.AlterMirrorStatusResponseDataList, *nethttp.Response, error) {
	httpResp := &nethttp.Response{
		StatusCode: nethttp.StatusOK,
	}

	errorMsg := "Not authorized"
	var errorCode int32 = 401

	return krsdk.AlterMirrorStatusResponseDataList{
		Kind:     "",
		Metadata: krsdk.ResourceCollectionMetadata{},
		Data: []krsdk.AlterMirrorStatusResponseData{
			{
				Kind:                 "",
				Metadata:             krsdk.ResourceMetadata{},
				DestinationTopicName: "dest-topic-1",
				ErrorMessage:         nil,
				ErrorCode:            nil,
			},
			{
				Kind:                 "",
				Metadata:             krsdk.ResourceMetadata{},
				DestinationTopicName: "dest-topic-2",
				ErrorMessage:         &errorMsg,
				ErrorCode:            &errorCode,
			},
		},
	}, httpResp, nil
}
